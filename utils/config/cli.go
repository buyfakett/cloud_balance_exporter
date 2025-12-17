package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// CLIConfig 表示命令行配置
type CLIConfig struct {
	ShowVersion bool
	ConfigFile  string
	Port        int
}

var CliCfg CLIConfig

// ParseCLI 解析命令行参数
func ParseCLI() {
	// 定义命令行参数
	flag.BoolVar(&CliCfg.ShowVersion, "version", false, "显示版本信息")
	flag.BoolVar(&CliCfg.ShowVersion, "v", false, "显示版本信息 (简写)")
	flag.StringVar(&CliCfg.ConfigFile, "config", "", "配置文件路径")
	flag.StringVar(&CliCfg.ConfigFile, "c", "", "配置文件路径 (简写)")
	flag.IntVar(&CliCfg.Port, "port", 8888, "服务端口")
	flag.IntVar(&CliCfg.Port, "p", 8888, "服务端口 (简写)")

	// 解析命令行参数
	flag.Parse()
}

// ShowVersionAndExit 显示版本信息并退出
func ShowVersionAndExit(version string) {
	version = strings.TrimSpace(version)
	fmt.Printf("cloud_balance_exporter %s\n", version)
	os.Exit(0)
}
