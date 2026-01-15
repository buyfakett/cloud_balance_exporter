package baidu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// BaiduCloudBalanceResponse represents the response from Baidu Cloud API
type BaiduCloudBalanceResponse struct {
	CashBalance float64 `json:"cashBalance"`
}

// generateAuthorization generates the authorization string for Baidu Cloud API
func generateAuthorization(ak, sk, method, path, date string) string {
	// Create the canonical request string
	canonicalRequest := fmt.Sprintf("%s\n%s\n\n", strings.ToUpper(method), path)

	// Create the string to sign
	stringToSign := fmt.Sprintf("bce-auth-v1/%s/%s/1800\n%s", ak, date, canonicalRequest)

	// Generate the signature
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Create the authorization string
	authorization := fmt.Sprintf("bce-auth-v1/%s/%s/1800/%s", ak, date, signature)

	return authorization
}

// BaiduYunBalance returns the account balance information from Baidu Cloud
func BaiduYunBalance(ak, sk string) ([]BaiduCloudBalance, error) {
	// API endpoint
	url := "https://billing.baidubce.com/v1/finance/cash/balance"
	path := "/v1/finance/cash/balance"

	// Generate current date in RFC 2616 format
	date := time.Now().UTC().Format("2006-01-02 15:04:05")

	// Generate authorization string
	authorization := generateAuthorization(ak, sk, "POST", path, date)

	// Create request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Host", "billing.baidubce.com")
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var response BaiduCloudBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Create balance struct
	balance := BaiduCloudBalance{
		AccountBalance: AccountBalance{
			AccountID:   "", // Baidu Cloud API doesn't return account ID
			AccountType: 1,  // Default to main account
			Amount:      response.CashBalance,
			Currency:    "CNY", // Baidu Cloud uses CNY
		},
		CashBalance: response.CashBalance,
	}

	return []BaiduCloudBalance{balance}, nil
}
