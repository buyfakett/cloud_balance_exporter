package baidu

// AccountBalance represents common account balance information
type AccountBalance struct {
	AccountID   string
	AccountType int
	Amount      float64
	Currency    string
}

// BaiduCloudBalance represents Baidu Cloud specific account balance information
type BaiduCloudBalance struct {
	AccountBalance
	CashBalance float64
}
