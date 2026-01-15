package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud_balance_exporter/pkg/metrics"
	"cloud_balance_exporter/utils/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed config/default.yaml
var defaultConfig []byte

//go:embed version.txt
var version string

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

	// 收集所有启用账户的百度云指标
	for _, cfg := range config.Cfg.Cloud.Baidu {
		if cfg.Enabled {
			go collectBaiduMetrics(cfg)
		}
	}

	// 收集所有启用账户的京东云指标
	for _, cfg := range config.Cfg.Cloud.Jdcloud {
		if cfg.Enabled {
			go collectJdcloudMetrics(cfg)
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
