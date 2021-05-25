package storagedb

var AllFields = []string{
	"user_id",
	"file_id",
	"customer_id",
	"period_start",
	"paid_plan",
	"paid_amount",
	"period_end",
}

type Invoice struct {
	UserID      string  `db:"user_id"`
	FileID      string  `db:"file_id"`
	CustomerId  uint32  `db:"customer_id"`
	PeriodStart string  `db:"period_start"`
	PaidPlan    string  `db:"paid_plan"`
	PaidAmount  float32 `db:"paid_amount"`
	PeriodEnd   string  `db:"period_end"`
}

type MRR struct {
	New          float32
	Old          float32
	Reactivation float32
	Expansion    float32
	Contraction  float32
	Churn        float32
}

type MPP struct {
	UserFileID  string
	PeriodStart string
	PeriodEnd   string
	Dates       []string
}
