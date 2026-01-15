package jdcloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// JdcloudBalanceResponse represents the response from JD Cloud API
type JdcloudBalanceResponse struct {
	Result struct {
		AvailableBalance float64 `json:"availableBalance"`
	} `json:"result"`
	RequestID string `json:"requestId"`
}

// generateSignature generates the signature for JD Cloud API
func generateSignature(sk, method, path, date, contentType, body string) string {
	// Create the string to sign
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", method, path, date, contentType)
	if body != "" {
		signStr += body
	}

	// Generate the signature
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(signStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// JdcloudYunBalance returns the account balance information from JD Cloud
func JdcloudYunBalance(ak, sk string) ([]JdcloudCloudBalance, error) {
	// API endpoint
	url := "https://billing.jdcloud-api.com/v1/account/balance"
	path := "/v1/account/balance"

	// Generate current date in RFC 1123 format
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	contentType := "application/json"

	// Generate signature
	signature := generateSignature(sk, "GET", path, date, contentType, "")

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Host", "billing.jdcloud-api.com")
	req.Header.Set("Date", date)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Jdcloud-Request-Id", fmt.Sprintf("%d", time.Now().UnixNano()))
	req.Header.Set("Authorization", fmt.Sprintf("JDCloud-HMAC-SHA256 credential=%s, signedHeaders=content-type;date;host, signature=%s", ak, signature))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JD Cloud API request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var response JdcloudBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Create balance struct
	balance := JdcloudCloudBalance{
		AccountBalance: AccountBalance{
			AccountID:   "", // JD Cloud API doesn't return account ID in this endpoint
			AccountType: 1,  // Default to main account
			Amount:      response.Result.AvailableBalance,
			Currency:    "CNY", // JD Cloud uses CNY
		},
		CashBalance: response.Result.AvailableBalance,
	}

	return []JdcloudCloudBalance{balance}, nil
}
