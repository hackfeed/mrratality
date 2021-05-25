package utils

import (
	storagedb "backend/db/storage"
	"fmt"
	"time"
)

var layout = "2006-01-02"

func GetAnalytics(userID, fileID, periodStart, periodEnd string) (string, error) {
	periodStartDate, err := time.Parse(layout, periodStart)
	if err != nil {
		return "", err
	}
	periodEndDate, err := time.Parse(layout, periodEnd)
	if err != nil {
		return "", err
	}

	months := getMonthsBetween(periodEndDate, periodStartDate)

	userfileID := fmt.Sprintf("%v.%v", userID, fileID)

	mpp := storagedb.MPP{userfileID, periodStart, periodEnd, months}
	str, err := storagedb.CreateMPPTable(storagedb.DB, mpp)
	if err != nil {
		fmt.Println(err)
		return str, err
	}

	return str, nil
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
	if cnt < 0 {
		cnt += 12
	}
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
