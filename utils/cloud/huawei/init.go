package cloud

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	bss "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
	bssRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/region"
)

// HuaweiYunBalance returns the account balances information
func HuaweiYunBalance(ak string, sk string) ([]HuaweiCloudBalance, error) {
	// Configure authentication with provided credentials
	auth, err := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	// Get available region
	region, err := bssRegion.SafeValueOf("cn-north-1")
	if err != nil {
		return nil, err
	}

	// Create a service client
	hcClient, err := bss.BssClientBuilder().
		WithRegion(region).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	client := bss.NewBssClient(hcClient)

	request := &model.ShowCustomerAccountBalancesRequest{}
	response, err := client.ShowCustomerAccountBalances(request)
	if err != nil {
		return nil, err
	}

	var balances []HuaweiCloudBalance
	if response.AccountBalances != nil && len(*response.AccountBalances) > 0 {
		// Only take the first account balance
		ab := (*response.AccountBalances)[0]
		amount := ab.Amount.InexactFloat64()
		designatedAmount := ab.DesignatedAmount.InexactFloat64()
		creditAmount := ab.CreditAmount.InexactFloat64()

		balances = append(balances, HuaweiCloudBalance{
			AccountBalance: AccountBalance{
				AccountID:   ab.AccountId,
				AccountType: int(ab.AccountType),
				Amount:      amount,
				Currency:    ab.Currency,
			},
			DesignatedAmount: designatedAmount,
			CreditAmount:     creditAmount,
			MeasureID:        int(ab.MeasureId),
		})
	}

	return balances, nil
}
