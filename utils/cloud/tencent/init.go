package tencent

import (
	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// TencentYunBalance returns the account balances information
func TencentYunBalance(ak string, sk string) ([]TencentCloudBalance, error) {
	// Configure authentication with provided credentials
	credential := common.NewCredential(ak, sk)

	// Set region - default to ap-shanghai if not specified
	region := "ap-shanghai"

	// Create client profile
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "billing.tencentcloudapi.com"

	// Create billing client
	client, err := billing.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}

	// Create request
	request := billing.NewDescribeAccountBalanceRequest()

	// Send request
	response, err := client.DescribeAccountBalance(request)
	if err != nil {
		return nil, err
	}

	var balances []TencentCloudBalance

	// Extract balance information from response
	if response.Response != nil {
		// Convert int64 to float64 and divide by 100 (Tencent uses fen, 1元 = 100分)
		balanceFloat := float64(*response.Response.Balance) / 100.0

		balance := TencentCloudBalance{
			AccountBalance: AccountBalance{
				AccountID:   "", // Tencent API doesn't return account ID in this endpoint
				AccountType: 1,  // Default to 1 for main account
				Amount:      balanceFloat,
				Currency:    "CNY", // Default to CNY for Tencent Cloud
			},
			AvailableAmount: balanceFloat,
			CreditAmount:    0.0, // Not available in this endpoint
			PendingAmount:   0.0, // Not available in this endpoint
			CashAmount:      balanceFloat,
			VoucherAmount:   0.0, // Not available in this endpoint
		}

		balances = append(balances, balance)
	}

	return balances, nil
}
