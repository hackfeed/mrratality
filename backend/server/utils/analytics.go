package utils

import (
	storagedb "backend/db/storage"
	"fmt"
	"time"
)

var layout = "2006-01-02"
var layoutDB = time.RFC3339Nano

// Refactor this shit
func GetAnalytics(userID, fileID, periodStart, periodEnd string) error {
	periodStartDate, _ := time.Parse(layout, periodStart)
	periodEndDate, _ := time.Parse(layout, periodEnd)

	months := getMonthsBetween(periodEndDate, periodStartDate)
	monthsCount := len(months)

	ufID := fmt.Sprintf("%v.%v", userID, fileID)

	ins := true

	mpp := storagedb.MPP{ufID, periodStart, periodEnd, months}
	err := storagedb.CreateMPPTable(storagedb.DB, mpp)
	if err != nil {
		if err.Error() == "table already exists" {
			ins = false
		}
	}

	npe := periodEndDate.AddDate(0, 1, -1)
	npestr := npe.Format(layout)

	invoices, err := storagedb.ReadByPeriod("mrr.storage", storagedb.AllFields, userID, fileID, periodStart, npestr)
	if err != nil {
		return err
	}

	mppEntries := formMPPEntries(invoices, monthsCount, periodStartDate)
	if ins {
		for _, entry := range mppEntries {
			err := storagedb.CreateDynamic(fmt.Sprintf("mrr.`%v-%v-%v`", ufID, periodStart, periodEnd), entry)
			if err != nil {
				return err
			}
		}
	}

	data, err := storagedb.ReadDynamic(fmt.Sprintf(`mrr."%v-%v-%v"`, ufID, periodStart, periodEnd))
	if err != nil {
		return err
	}
	fmt.Println(monthsCount)
	totalMRR := make([]storagedb.MRR, monthsCount)

	for _, client := range data {
		clientFloat := []float32{}
		for _, x := range client {
			clientFloat = append(clientFloat, x.(float32))
		}
		clientMRR := calculateMRR(clientFloat)

		for i := 0; i < monthsCount; i++ {
			totalMRR[i].New += clientMRR[i].New
			totalMRR[i].Old += clientMRR[i].Old
			totalMRR[i].Reactivation += clientMRR[i].Reactivation
			totalMRR[i].Expansion += clientMRR[i].Expansion
			totalMRR[i].Contraction += clientMRR[i].Contraction
			totalMRR[i].Churn += clientMRR[i].Churn
		}
	}

	return nil
}

func calculateMRR(moneyPerMonth []float32) []storagedb.MRR {
	monthsCount := len(moneyPerMonth)
	clientMRR := make([]storagedb.MRR, monthsCount)

	isNew := true

	for i := range moneyPerMonth {
		var monthMRR storagedb.MRR

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
