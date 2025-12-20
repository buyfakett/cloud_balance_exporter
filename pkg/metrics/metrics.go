package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 全局Prometheus指标
var (
	// 账户余额指标
	accountBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloud_account_balance",
			Help: "Cloud provider account balance",
		},
		[]string{"account_type", "name", "threshold"},
	)
	// 账户余额状态指标 (1=正常, 0=异常)
	accountBalanceStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloud_account_balance_status",
			Help: "Cloud provider account balance status (1=normal, 0=abnormal)",
		},
		[]string{"account_type", "name", "threshold"},
	)
)

// Register 注册所有Prometheus指标
func Register() {
	prometheus.MustRegister(accountBalance, accountBalanceStatus)
}

// SetAccountBalance 设置账户余额指标
func SetAccountBalance(accountType, name, threshold string, amount float64) {
	accountBalance.WithLabelValues(accountType, name, threshold).Set(amount)
}

// SetAccountBalanceStatus 设置账户余额状态指标 (1=正常, 0=异常)
func SetAccountBalanceStatus(accountType, name string, threshold string, status float64) {
	accountBalanceStatus.WithLabelValues(accountType, name, threshold).Set(status)
}
