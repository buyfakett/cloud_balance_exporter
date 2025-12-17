package alibaba

import (
	"log"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/bssopenapi-20171214/v2/client"
	opopenapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AlibabaYunBalance 返回账户余额信息
func parseAmount(amountStr string) (float64, error) {
	// parseAmount 从金额字符串中解析浮点数（例如："1,157.74" -> 1157.74）
	cleanStr := strings.ReplaceAll(amountStr, ",", "")
	return strconv.ParseFloat(cleanStr, 64)
}

// AlibabaCloudBalance returns the account balances information
func AlibabaYunBalance(ak string, sk string) ([]AlibabaCloudBalance, error) {
	// 创建配置
	config := &opopenapi.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
		RegionId:        tea.String("cn-shanghai"),
	}

	// 创建客户端
	c, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 发送请求（无需参数）
	response, err := c.QueryAccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []AlibabaCloudBalance

	// 提取余额信息
	if response.Body != nil && response.Body.Data != nil {
		// 转换字符串金额为float64，带错误处理和逗号移除
		availableAmountStr := tea.StringValue(response.Body.Data.AvailableAmount)
		availableAmount, err := parseAmount(availableAmountStr)
		if err != nil {
			log.Printf("解析可用余额失败: '%s', 错误: %v", availableAmountStr, err)
			availableAmount = 0.0
		}

		creditAmountStr := tea.StringValue(response.Body.Data.CreditAmount)
		creditAmount, err := parseAmount(creditAmountStr)
		if err != nil {
			log.Printf("解析信用余额失败: '%s', 错误: %v", creditAmountStr, err)
			creditAmount = 0.0
		}

		availableCashAmountStr := tea.StringValue(response.Body.Data.AvailableCashAmount)
		availableCashAmount, err := parseAmount(availableCashAmountStr)
		if err != nil {
			log.Printf("解析可用现金余额失败: '%s', 错误: %v", availableCashAmountStr, err)
			availableCashAmount = 0.0
		}

		// 从响应中获取货币
		currency := tea.StringValue(response.Body.Data.Currency)
		if currency == "" {
			currency = "CNY" // 默认为人民币如果未提供
		}

		balance := AlibabaCloudBalance{
			AccountBalance: AccountBalance{
				AccountID:   "", // 阿里云此端点不返回账户ID
				AccountType: 1,  // 默认为主账户
				Amount:      availableAmount,
				Currency:    currency,
			},
			AvailableAmount: availableAmount,
			CreditAmount:    creditAmount,
			PendingAmount:   0.0, // 响应中不可用
			CashAmount:      availableCashAmount,
			CouponAmount:    0.0, // 响应中不可用
		}

		balances = append(balances, balance)
	} else {
		log.Printf("阿里云API响应体或数据为空")
	}

	return balances, nil
}
