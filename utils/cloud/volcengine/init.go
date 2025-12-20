package volcengine

import (
	"strconv"

	"github.com/volcengine/volcengine-go-sdk/service/billing"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

// VolcengineYunBalance returns the account balances information
func VolcengineYunBalance(ak string, sk string) ([]VolcengineCloudBalance, error) {
	// Create configuration
	config := volcengine.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(ak, sk, "")).
		WithRegion("cn-beijing") // Volcengine uses cn-beijing as default region

	// Create session
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	// Create billing service client
	svc := billing.New(sess)

	// Create request for account balance
	input := &billing.QueryBalanceAcctInput{}

	// Send request
	result, err := svc.QueryBalanceAcct(input)
	if err != nil {
		return nil, err
	}

	var balances []VolcengineCloudBalance

	// Extract balance information from response
	if result != nil {
		// Parse available balance
		availableBalance := 0.0
		if result.AvailableBalance != nil {
			availableBalance, _ = strconv.ParseFloat(*result.AvailableBalance, 64)
		}

		// Parse cash balance
		cashBalance := 0.0
		if result.CashBalance != nil {
			cashBalance, _ = strconv.ParseFloat(*result.CashBalance, 64)
		}

		// Parse credit limit
		creditLimit := 0.0
		if result.CreditLimit != nil {
			creditLimit, _ = strconv.ParseFloat(*result.CreditLimit, 64)
		}

		// Parse arrears balance
		arrearsBalance := 0.0
		if result.ArrearsBalance != nil {
			arrearsBalance, _ = strconv.ParseFloat(*result.ArrearsBalance, 64)
		}

		// Default currency to CNY for Volcengine
		currency := "CNY"

		balance := VolcengineCloudBalance{
			AccountBalance: AccountBalance{
				AccountID:   "", // Volcengine API doesn't return account ID in a usable format
				AccountType: 1,  // Default to 1 for main account
				Amount:      availableBalance,
				Currency:    currency,
			},
			AvailableAmount: availableBalance,
			CreditAmount:    creditLimit,
			PendingAmount:   arrearsBalance,
			CashAmount:      cashBalance,
			CouponAmount:    0.0, // Not available in this endpoint
		}

		balances = append(balances, balance)
	}

	return balances, nil
}
