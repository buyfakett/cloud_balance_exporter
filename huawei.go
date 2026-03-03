package main

import (
	"fmt"
	"log"

	"cloud_balance_exporter/pkg/metrics"
	cloud "cloud_balance_exporter/utils/cloud/huawei"
	"cloud_balance_exporter/utils/config"
)

// collectHuaweiMetrics 从华为云收集余额指标
func collectHuaweiMetrics(cfg config.HuaweiCloudConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("收集华为云指标时发生恐慌，账户 %s: %v", cfg.Name, r)
		}
	}()

	balances, err := cloud.HuaweiYunBalance(cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		log.Printf("收集华为云余额指标失败，账户 %s: %v", cfg.Name, err)
		return
	}

	// 更新账户余额指标
	for _, balance := range balances {
		// 设置实际余额 - 华为云使用空account_id
		var threshold string = fmt.Sprint(cfg.BalanceThreshold)
		metrics.SetAccountBalance(
			"huawei", // 账户类型现在是云提供商类型
			cfg.Name, // 添加配置中的名称到标签
			threshold,
			balance.Amount,
		)

		// 设置余额状态 (1=正常, 0=异常)
		var status float64 = 1.0 // 默认为正常
		if cfg.BalanceThreshold > 0 && balance.Amount < cfg.BalanceThreshold {
			status = 0.0 // 如果余额小于阈值则为异常
		}
		metrics.SetAccountBalanceStatus(
			"huawei",
			cfg.Name,
			threshold,
			status,
		)
		log.Printf("成功收集华为云余额指标，账户 %s ,余额为 %s", cfg.Name, fmt.Sprint(balance.Amount))
	}
}
