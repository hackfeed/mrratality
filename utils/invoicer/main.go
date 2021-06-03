package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var header = []string{"customer_id", "period_start", "paid_plan", "paid_amount", "period_end"}
var paidPlans = []string{"annually", "monthly"}
var layout = "02.01.2006"

func main() {
	file, err := os.OpenFile(os.Args[1], os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Write(header)
	for i := 0; i < 1000; i++ {
		customerID := randRange(1, 10000)
		periodStart := fmt.Sprintf("01.%02d.%d", randRange(1, 12), randRange(2018, 2021))
		periodStartDate, _ := time.Parse(layout, periodStart)
		paidPlan := paidPlans[randRange(0, 1)]
		var periodEndDate time.Time
		if paidPlan == "annually" {
			periodEndDate = periodStartDate.AddDate(1, 0, -1)
		} else {
			periodEndDate = periodStartDate.AddDate(0, 1, -1)
		}
		periodEnd := periodEndDate.Format(layout)
		paidAmount := randRange(1, 1000)
		entry := []string{fmt.Sprint(customerID), periodStart, paidPlan, fmt.Sprint(paidAmount), periodEnd}
		csvWriter.Write(entry)
	}
	csvWriter.Flush()
}

func randRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
