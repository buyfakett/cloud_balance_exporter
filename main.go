package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud_balance_exporter/pkg/metrics"
	"cloud_balance_exporter/utils/cloud/alibaba"
	cloud "cloud_balance_exporter/utils/cloud/huawei"
	"cloud_balance_exporter/utils/cloud/tencent"
	"cloud_balance_exporter/utils/cloud/volcengine"
	"cloud_balance_exporter/utils/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed config/default.yaml
var defaultConfig []byte

//go:embed version.txt
var version string

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
	}

	log.Printf("成功收集华为云余额指标，账户 %s", cfg.Name)
}

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
	}

	log.Printf("成功收集阿里云余额指标，账户 %s", cfg.Name)
}

// collectTencentMetrics 从腾讯云收集余额指标
func collectTencentMetrics(cfg config.TencentCloudConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("收集腾讯云指标时发生恐慌，账户 %s: %v", cfg.Name, r)
		}
	}()

	balances, err := tencent.TencentYunBalance(cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		log.Printf("收集腾讯云余额指标失败，账户 %s: %v", cfg.Name, err)
		return
	}

	// Update account balance metrics
	for _, balance := range balances {
		// 设置实际余额 - 腾讯云使用空account_id
		var threshold string = fmt.Sprint(cfg.BalanceThreshold)
		metrics.SetAccountBalance(
			"tencent", // 账户类型现在是云提供商类型
			cfg.Name,  // 使用配置中的名称
			threshold,
			balance.Amount,
		)

		// 设置余额状态 (1=正常, 0=异常) - 腾讯云使用空account_id
		var status float64 = 1.0 // Default to normal
		if cfg.BalanceThreshold > 0 && balance.Amount < cfg.BalanceThreshold {
			status = 0.0 // Abnormal if balance < threshold
		}
		metrics.SetAccountBalanceStatus(
			"tencent",
			cfg.Name,
			threshold,
			status,
		)
	}

	log.Printf("成功收集腾讯云余额指标，账户 %s", cfg.Name)
}

// collectVolcengineMetrics 从火山引擎收集余额指标
func collectVolcengineMetrics(cfg config.VolcengineCloudConfig) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("收集火山引擎指标时发生恐慌，账户 %s: %v", cfg.Name, r)
		}
	}()

	balances, err := volcengine.VolcengineYunBalance(cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		log.Printf("收集火山引擎余额指标失败，账户 %s: %v", cfg.Name, err)
		return
	}

	// Update account balance metrics
	for _, balance := range balances {
		// 设置实际余额 - 火山引擎使用空account_id
		var threshold string = fmt.Sprint(cfg.BalanceThreshold)
		metrics.SetAccountBalance(
			"volcengine", // 账户类型现在是云提供商类型
			cfg.Name,     // 使用配置中的名称
			threshold,
			balance.Amount,
		)

		// 设置余额状态 (1=正常, 0=异常) - 火山引擎使用空account_id
		var status float64 = 1.0 // Default to normal
		if cfg.BalanceThreshold > 0 && balance.Amount < cfg.BalanceThreshold {
			status = 0.0 // Abnormal if balance < threshold
		}
		metrics.SetAccountBalanceStatus(
			"volcengine",
			cfg.Name,
			threshold,
			status,
		)
	}

	log.Printf("成功收集火山引擎余额指标，账户 %s", cfg.Name)
}

// collectMetricsFromConfig 从配置中所有启用的云提供商收集指标
func collectMetricsFromConfig() {
	// 收集所有启用账户的华为云指标
	for _, cfg := range config.Cfg.Cloud.Huawei {
		if cfg.Enabled {
			go collectHuaweiMetrics(cfg)
		}
	}

	// 收集所有启用账户的阿里云指标
	for _, cfg := range config.Cfg.Cloud.Alibaba {
		if cfg.Enabled {
			go collectAlibabaMetrics(cfg)
		}
	}

	// 收集所有启用账户的腾讯云指标
	for _, cfg := range config.Cfg.Cloud.Tencent {
		if cfg.Enabled {
			go collectTencentMetrics(cfg)
		}
	}

	// 收集所有启用账户的火山引擎指标
	for _, cfg := range config.Cfg.Cloud.Volcengine {
		if cfg.Enabled {
			go collectVolcengineMetrics(cfg)
		}
	}
}

func main() {
	// 初始化配置（包括命令行参数处理）
	config.InitConfig(defaultConfig)

	// 如果显示版本信息，直接退出
	if config.CliCfg.ShowVersion {
		config.ShowVersionAndExit(version)
	}

	// 注册指标
	metrics.Register()

	// 在goroutine中启动指标收集
	go func() {
		// 初始收集
		collectMetricsFromConfig()

		// 根据配置设置定期收集
		skipInterval := time.Duration(config.Cfg.Server.ScrapeInterval) * time.Minute
		ticker := time.NewTicker(skipInterval)
		defer ticker.Stop()

		for range ticker.C {
			collectMetricsFromConfig()
		}
	}()

	// 为指标设置HTTP服务器
	serverAddress := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("在 %s 启动服务器", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
