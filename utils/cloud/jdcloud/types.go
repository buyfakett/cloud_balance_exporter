package jdcloud

// AccountBalance represents common account balance information
type AccountBalance struct {
	AccountID   string
	AccountType int
	Amount      float64
	Currency    string
}

// JdcloudCloudBalance represents JD Cloud specific account balance information
type JdcloudCloudBalance struct {
	AccountBalance
	CashBalance float64
}
