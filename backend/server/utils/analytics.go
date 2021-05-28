package utils

import (
	cachedb "backend/db/cache"
	storagedb "backend/db/storage"
	"backend/server/models"
	"errors"
	"fmt"
	"time"
)

var layout = "2006-01-02"
var layoutDB = time.RFC3339Nano

func GetAnalytics(userID, fileID, periodStart, periodEnd string) ([]string, models.TotalMRR, error) {
	var (
		totalMRR models.TotalMRR
		months   []string
	)

	periodStartDate, _ := time.Parse(layout, periodStart)
	periodEndDate, _ := time.Parse(layout, periodEnd)

	if periodStartDate.After(periodEndDate) {
		return months, totalMRR, errors.New("Period start should be less than period end")
	}

	ufID := fmt.Sprintf("%v.%v", userID, fileID)
	ufp := fmt.Sprintf("%v-%v-%v", ufID, periodStart, periodEnd)

	mrr, err := cachedb.Read("cache", ufp)
	if err != nil {
		return months, totalMRR, err
	}

	months = getMonthsBetween(periodEndDate, periodStartDate)

	if len(mrr) != 0 {
		totalMRR = convertFetchedMRR(mrr)
	} else {
		mpp := models.MPP{ufID, periodStart, periodEnd, months}
		err = formMPPTable(mpp, userID, fileID, ufp, periodStartDate, periodEndDate)
		if err != nil {
			return months, totalMRR, err
		}

		totalMoneyPerMonth, err := storagedb.ReadDynamic(fmt.Sprintf("mrr.`%v`", ufp))
		if err != nil {
			return months, totalMRR, err
		}

		rawMRR := calculateTotalMRR(totalMoneyPerMonth)
		totalMRR = convertRawMRR(rawMRR)
		totalMRRMap := map[string]interface{}{"mrr": totalMRR}

		err = cachedb.Create("cache", []interface{}{ufp, totalMRRMap})
		if err != nil {
			return months, totalMRR, err
		}
	}

	return months, totalMRR, err
}

func convertRawMRR(rawMRR []models.MRR) models.TotalMRR {
	var totalMRR models.TotalMRR

	for _, mrr := range rawMRR {
		totalMRR.New = append(totalMRR.New, mrr.New)
		totalMRR.Old = append(totalMRR.Old, mrr.Old)
		totalMRR.Reactivation = append(totalMRR.Reactivation, mrr.Reactivation)
		totalMRR.Expansion = append(totalMRR.Expansion, mrr.Expansion)
		totalMRR.Contraction = append(totalMRR.Contraction, mrr.Contraction)
		totalMRR.Churn = append(totalMRR.Churn, mrr.Churn)
		totalMRR.Total = append(totalMRR.Total, mrr.New+mrr.Old+mrr.Reactivation+mrr.Expansion+mrr.Contraction+mrr.Churn)
	}

	return totalMRR
}

func convertFetchedMRR(fetchedMRR []interface{}) models.TotalMRR {
	var totalMRR models.TotalMRR

	mrrMap := fetchedMRR[0].([]interface{})[1].(map[interface{}]interface{})["mrr"].(map[interface{}]interface{})
	for k, v := range mrrMap {
		vfloat := []float32{}
		for _, vv := range v.([]interface{}) {
			vfloat = append(vfloat, vv.(float32))
		}
		if k == "New" {
			totalMRR.New = vfloat
		}
		if k == "Old" {
			totalMRR.Old = vfloat
		}
		if k == "Reactivation" {
			totalMRR.Reactivation = vfloat
		}
		if k == "Expansion" {
			totalMRR.Expansion = vfloat
		}
		if k == "Contraction" {
			totalMRR.Contraction = vfloat
		}
		if k == "Churn" {
			totalMRR.Churn = vfloat
		}
		if k == "Total" {
			totalMRR.Total = vfloat
		}
	}

	return totalMRR
}

func calculateTotalMRR(totalMoneyPerMonth [][]interface{}) []models.MRR {
	monthsCount := len(totalMoneyPerMonth[0])
	totalMRR := make([]models.MRR, monthsCount)

	for _, moneyPerMonth := range totalMoneyPerMonth {
		clientFloat := []float32{}
		for _, m := range moneyPerMonth {
			clientFloat = append(clientFloat, m.(float32))
		}
		clientMRR := calculateClientMRR(clientFloat)

		for i := 0; i < monthsCount; i++ {
			totalMRR[i].New += clientMRR[i].New
			totalMRR[i].Old += clientMRR[i].Old
			totalMRR[i].Reactivation += clientMRR[i].Reactivation
			totalMRR[i].Expansion += clientMRR[i].Expansion
			totalMRR[i].Contraction += clientMRR[i].Contraction
			totalMRR[i].Churn += clientMRR[i].Churn
		}
	}

	return totalMRR
}

func calculateClientMRR(moneyPerMonth []float32) []models.MRR {
	monthsCount := len(moneyPerMonth)
	clientMRR := make([]models.MRR, monthsCount)

	isNew := true

	for i := range moneyPerMonth {
		var monthMRR models.MRR

		if moneyPerMonth[i] > 0 && isNew {
			monthMRR.New = moneyPerMonth[i]
			isNew = false
		} else if i > 0 && moneyPerMonth[i] == moneyPerMonth[i-1] {
			monthMRR.Old = moneyPerMonth[i]
		} else if i > 0 && moneyPerMonth[i-1] == 0 && moneyPerMonth[i] > 0 && !isNew {
			monthMRR.Reactivation = moneyPerMonth[i]
		} else if i > 0 && moneyPerMonth[i] > moneyPerMonth[i-1] && moneyPerMonth[i-1] != 0 {
			monthMRR.Expansion = moneyPerMonth[i] - moneyPerMonth[i-1]
		} else if i > 0 && moneyPerMonth[i] < moneyPerMonth[i-1] && moneyPerMonth[i] != 0 {
			monthMRR.Contraction = moneyPerMonth[i] - moneyPerMonth[i-1]
		} else if i > 0 && moneyPerMonth[i] == 0 && moneyPerMonth[i-1] > 0 {
			monthMRR.Churn = -moneyPerMonth[i-1]
		}

		clientMRR[i] = monthMRR
	}

	return clientMRR
}

func formMPPTable(mpp models.MPP, userID, fileID, ufp string, periodStartDate, periodEndDate time.Time) error {
	ins := true

	err := storagedb.CreateMPPTable(storagedb.DB, mpp)
	if err != nil {
		if err.Error() == "table already exists" {
			ins = false
		}
	}

	normPeriodEndDate := periodEndDate.AddDate(0, 1, -1)
	normPeriodEnd := normPeriodEndDate.Format(layout)
	periodStart := periodStartDate.Format(layout)

	invoices, err := storagedb.ReadByPeriod("mrr.storage", storagedb.AllFields, userID, fileID, periodStart, normPeriodEnd)
	if err != nil {
		return err
	}

	if len(invoices) == 0 {
		return errors.New("No data found for given period")
	}

	if ins {
		mppEntries := formMPPEntries(invoices, len(mpp.Dates), periodStartDate)
		for _, entry := range mppEntries {
			err := storagedb.CreateDynamic(fmt.Sprintf("mrr.`%v`", ufp), entry)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func formMPPEntries(invoices []storagedb.Invoice, monthsCount int, periodStartDate time.Time) [][]interface{} {
	invoicesCount := len(invoices)
	mppEntries := make([][]interface{}, invoicesCount)

	for i, invoice := range invoices {
		moneyPerMonth := make([]float32, monthsCount)
		paidAmount := invoice.PaidAmount
		periodLen := 1
		invoicePeriodStartDate, _ := time.Parse(layoutDB, invoice.PeriodStart)
		startMonth := getMonthsDiff(invoicePeriodStartDate, periodStartDate)
		if startMonth < 0 {
			periodLen += startMonth
			startMonth = 0
		}
		if invoice.PaidPlan == "annually" {
			paidAmount /= 12
			periodLen = 12
		}
		for j := startMonth; j < monthsCount; j++ {
			if periodLen <= 0 {
				paidAmount = 0
			}
			moneyPerMonth[j] = paidAmount
			periodLen--
		}
		moneyFlow := []interface{}{invoice.CustomerId}
		for _, m := range moneyPerMonth {
			moneyFlow = append(moneyFlow, m)
		}
		mppEntries[i] = moneyFlow
	}

	return mppEntries
}

func getMonthsBetween(fdate, sdate time.Time) []string {
	if fdate.Location() != sdate.Location() {
		sdate = sdate.In(fdate.Location())
	}
	if fdate.After(sdate) {
		fdate, sdate = sdate, fdate
	}

	fyear, fmonth, _ := fdate.Date()
	syear, smonth, _ := sdate.Date()

	count := int(smonth - fmonth)
	count += 12*(syear-fyear) + 1

	months := []string{}

	for i := 0; i < count; i++ {
		yearsDif := (int(fmonth) + i) / 12
		month := (int(fmonth) + i) % 12
		if month == 0 {
			yearsDif -= 1
			month = 12
		}
		months = append(months, fmt.Sprintf("%v.%v", month, fyear+yearsDif))
	}

	return months
}

func getMonthsDiff(fdate, sdate time.Time) int {
	if fdate.Location() != sdate.Location() {
		sdate = sdate.In(fdate.Location())
	}

	fyear, fmonth, _ := fdate.Date()
	syear, smonth, _ := sdate.Date()

	if fyear > syear {
		fyear, syear = syear, fyear
	}

	count := int(fmonth - smonth)
	count += 12 * (syear - fyear)

	return count
}
