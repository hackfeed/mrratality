package utils

import (
	cachedb "backend/db/cache"
	storagedb "backend/db/storage"
	"backend/server/models"
	"fmt"
	"time"
)

var layout = "2006-01-02"
var layoutDB = time.RFC3339Nano

// Refactor this shit
func GetAnalytics(userID, fileID, periodStart, periodEnd string) ([]string, []models.MRR, error) {
	periodStartDate, _ := time.Parse(layout, periodStart)
	periodEndDate, _ := time.Parse(layout, periodEnd)

	months := getMonthsBetween(periodEndDate, periodStartDate)
	monthsCount := len(months)

	ufID := fmt.Sprintf("%v.%v", userID, fileID)
	ufp := fmt.Sprintf("%v-%v-%v", ufID, periodStart, periodEnd)

	mrr, err := cachedb.Read("cache", ufp)
	if err != nil {
		fmt.Println(err)
	}

	var totalMRR []models.MRR

	if len(mrr) != 0 {
		totalMRR = convertFetchedMRR(mrr, monthsCount)
	} else {
		ins := true

		mpp := models.MPP{ufID, periodStart, periodEnd, months}
		err = storagedb.CreateMPPTable(storagedb.DB, mpp)
		if err != nil {
			if err.Error() == "table already exists" {
				ins = false
			}
		}

		npe := periodEndDate.AddDate(0, 1, -1)
		npestr := npe.Format(layout)

		invoices, err := storagedb.ReadByPeriod("mrr.storage", storagedb.AllFields, userID, fileID, periodStart, npestr)
		if err != nil {
			return nil, nil, err
		}

		mppEntries := formMPPEntries(invoices, monthsCount, periodStartDate)
		if ins {
			for _, entry := range mppEntries {
				err := storagedb.CreateDynamic(fmt.Sprintf("mrr.`%v`", ufp), entry)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		totalMoneyPerMonth, err := storagedb.ReadDynamic(fmt.Sprintf("mrr.`%v`", ufp))
		if err != nil {
			return nil, nil, err
		}

		totalMRR = calculateTotalMRR(totalMoneyPerMonth)

		totalMRRMap := map[string]interface{}{"mrr": totalMRR}

		err = cachedb.Create("cache", []interface{}{ufp, totalMRRMap, 1})
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(totalMRR)

	return months, totalMRR, nil
}

func convertFetchedMRR(fetched []interface{}, monthsCount int) []models.MRR {
	totalMRRFetched := make([]models.MRR, monthsCount)

	mrrMap := fetched[0].([]interface{})[1].(map[interface{}]interface{})["mrr"].([]interface{})
	for i, v := range mrrMap {
		curMRR := models.MRR{}
		for k, vv := range v.(map[interface{}]interface{}) {
			if k == "New" {
				curMRR.New = vv.(float32)
			}
			if k == "Old" {
				curMRR.Old = vv.(float32)
			}
			if k == "Reactivation" {
				curMRR.Reactivation = vv.(float32)
			}
			if k == "Expansion" {
				curMRR.Expansion = vv.(float32)
			}
			if k == "Contraction" {
				curMRR.Contraction = vv.(float32)
			}
			if k == "Churn" {
				curMRR.Churn = vv.(float32)
			}
		}
		totalMRRFetched[i] = curMRR
	}

	return totalMRRFetched
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

	fy, fm, _ := fdate.Date()
	sy, sm, _ := sdate.Date()

	cnt := int(sm - fm)
	cnt += 12*(sy-fy) + 1

	months := []string{}

	for i := 0; i < cnt; i++ {
		dy := (int(fm) + i) / 12
		m := (int(fm) + i) % 12
		if m == 0 {
			dy -= 1
			m = 12
		}
		months = append(months, fmt.Sprintf("%v.%v", m, fy+dy))
	}

	return months
}

func getMonthsDiff(fdate, sdate time.Time) int {
	if fdate.Location() != sdate.Location() {
		sdate = sdate.In(fdate.Location())
	}

	fy, fm, _ := fdate.Date()
	sy, sm, _ := sdate.Date()

	if fy > sy {
		fy, sy = sy, fy
	}

	cnt := int(fm - sm)
	cnt += 12 * (sy - fy)

	return cnt
}
