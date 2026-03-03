package main

import (
	"fmt"
	"log"

	"cloud_balance_exporter/pkg/metrics"
	"cloud_balance_exporter/utils/cloud/alibaba"
	"cloud_balance_exporter/utils/config"
)

// collectAlibabaMetrics 从阿里云收集余额指标
func collectAlibabaMetrics(cfg config.AlibabaCloudConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("收集阿里云指标时发生恐慌，账户 %s: %v", cfg.Name, r)
		}
	}()

	balances, err := alibaba.AlibabaYunBalance(cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		log.Printf("收集阿里云余额指标失败，账户 %s: %v", cfg.Name, err)
		return
	}

	// Update account balance metrics
	for _, balance := range balances {
		// 设置实际余额 - 阿里云使用空account_id
		var threshold string = fmt.Sprint(cfg.BalanceThreshold)
		metrics.SetAccountBalance(
			"alibaba", // 账户类型现在是云提供商类型
			cfg.Name,  // 使用配置中的名称
			threshold,
			balance.Amount,
		)

		// 设置余额状态 (1=正常, 0=异常) - 阿里云使用空account_id
		var status float64 = 1.0 // Default to normal
		if cfg.BalanceThreshold > 0 && balance.Amount < cfg.BalanceThreshold {
			status = 0.0 // Abnormal if balance < threshold
		}
		metrics.SetAccountBalanceStatus(
			"alibaba",
			cfg.Name,
			threshold,
			status,
		)
		log.Printf("成功收集阿里云余额指标，账户 %s ,余额为 %s", cfg.Name, fmt.Sprint(balance.Amount))
	}
}
