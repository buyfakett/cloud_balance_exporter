package tencent

// AccountBalance represents common account balance information
// Each cloud provider may extend this with additional fields

type AccountBalance struct {
	AccountID   string
	AccountType int
	Amount      float64
	Currency    string
}

// TencentCloudBalance represents Tencent Cloud specific account balance information
type TencentCloudBalance struct {
	AccountBalance
	AvailableAmount float64
	CreditAmount    float64
	PendingAmount   float64
	CashAmount      float64
	VoucherAmount   float64
}
