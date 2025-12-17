package cloud

// AccountBalance represents common account balance information
// Each cloud provider may extend this with additional fields

type AccountBalance struct {
	AccountID   string
	AccountType int
	Amount      float64
	Currency    string
}

// HuaweiCloudBalance represents Huawei Cloud specific account balance information
type HuaweiCloudBalance struct {
	AccountBalance
	DesignatedAmount float64
	CreditAmount     float64
	MeasureID        int
}

// AlibabaCloudBalance represents Alibaba Cloud specific account balance information
type AlibabaCloudBalance struct {
	AccountBalance
	// Alibaba specific fields will be added here
}

// TencentCloudBalance represents Tencent Cloud specific account balance information
type TencentCloudBalance struct {
	AccountBalance
	// Tencent specific fields will be added here
}
