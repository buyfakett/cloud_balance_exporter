package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ServerConfig 表示服务器配置
type ServerConfig struct {
	Port           int    `mapstructure:"port"`
	Name           string `mapstructure:"name"`
	Author         string `mapstructure:"author"`
	Zone           string `mapstructure:"zone"`
	ScrapeInterval int    `mapstructure:"scrape_interval"` // 分钟为单位
}

// HuaweiCloudConfig 表示华为云特定配置
type HuaweiCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
}

// AlibabaCloudConfig 表示阿里云特定配置
type AlibabaCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
	// 阿里云特定配置字段将在这里添加
}

// TencentCloudConfig 表示腾讯云特定配置
type TencentCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
	// 腾讯云特定配置字段将在这里添加
}

// VolcengineCloudConfig 表示火山引擎特定配置
type VolcengineCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
	// 火山引擎特定配置字段将在这里添加
}

// BaiduCloudConfig 表示百度云特定配置
type BaiduCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
	// 百度云特定配置字段将在这里添加
}

// JdcloudCloudConfig 表示京东云特定配置
type JdcloudCloudConfig struct {
	Name             string  `mapstructure:"name"`
	AccessKey        string  `mapstructure:"access_key"`
	SecretKey        string  `mapstructure:"secret_key"`
	Enabled          bool    `mapstructure:"enabled"`
	BalanceThreshold float64 `mapstructure:"balance_threshold"`
	// 京东云特定配置字段将在这里添加
}

// CloudConfig 表示云提供商配置
type CloudConfig struct {
	Huawei     []HuaweiCloudConfig     `mapstructure:"huawei"`
	Alibaba    []AlibabaCloudConfig    `mapstructure:"alibaba"`
	Tencent    []TencentCloudConfig    `mapstructure:"tencent"`
	Volcengine []VolcengineCloudConfig `mapstructure:"volcengine"`
	Baidu      []BaiduCloudConfig      `mapstructure:"baidu"`
	Jdcloud    []JdcloudCloudConfig    `mapstructure:"jdcloud"`
}

// AppConfig 表示整个应用程序配置
type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	Cloud  CloudConfig  `mapstructure:"cloud"`
}

var Cfg AppConfig

func InitConfig(defaultConfigContent []byte) {
	// 1. 处理命令行参数
	ParseCLI()

	// 如果显示版本信息，直接退出
	if CliCfg.ShowVersion {
		return
	}

	v := viper.New()

	// 2. 加载外部配置文件（如果存在）
	if CliCfg.ConfigFile != "" {
		if _, err := os.Stat(CliCfg.ConfigFile); err == nil {
			v.SetConfigFile(CliCfg.ConfigFile)
			if err := v.ReadInConfig(); err != nil {
				fmt.Printf("加载外部配置失败: %v (路径: %s)\n", err, CliCfg.ConfigFile)
				os.Exit(1)
			}
		} else {
			fmt.Printf("警告: 外部配置文件不存在，使用默认配置 (路径: %s)\n", CliCfg.ConfigFile)
		}
	}

	// 3. 环境变量覆盖
	v.SetEnvPrefix("CLOUD_BALANCE_EXPORTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 4. 合并命令行参数
	_ = v.BindPFlags(pflag.CommandLine)

	// 5. 映射到结构体
	if err := v.Unmarshal(&Cfg); err != nil {
		fmt.Println("解析配置失败:", err)
		os.Exit(1)
	}

	// 6. 设置默认值
	if Cfg.Server.Port <= 0 {
		Cfg.Server.Port = CliCfg.Port
	}
	if Cfg.Server.Name == "" {
		Cfg.Server.Name = ServerName
	}
	if Cfg.Server.Author == "" {
		Cfg.Server.Author = Author
	}
	if Cfg.Server.ScrapeInterval <= 0 {
		Cfg.Server.ScrapeInterval = 5 // 默认5分钟
	}
}
