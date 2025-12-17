package alibaba

// AccountBalance represents common account balance information
// Each cloud provider may extend this with additional fields

type AccountBalance struct {
	AccountID   string
	AccountType int
	Amount      float64
	Currency    string
}

// AlibabaCloudBalance represents Alibaba Cloud specific account balance information
type AlibabaCloudBalance struct {
	AccountBalance
	AvailableAmount float64
	PendingAmount   float64
	CreditAmount    float64
	CashAmount      float64
	CouponAmount    float64
}
