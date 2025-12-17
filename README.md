# 云余额监控导出器 (Cloud Balance Exporter)

一个支持多云平台账户余额监控的 Prometheus 导出器，支持阿里云、华为云和腾讯云。

## 功能特性

- 阿里云余额监控
- 华为云余额监控  
- 腾讯云余额监控
- 多账户支持
- 余额阈值告警
- Prometheus 指标导出
- 自动单位转换

## 配置说明

### 基本配置

```yaml
cloud:
  # 阿里云配置
  alibaba:
    - name: "alibaba-account-1"
      access_key: "你的阿里云AccessKey"
      secret_key: "你的阿里云SecretKey"
      enabled: true
      balance_threshold: 1000.0  # 余额告警阈值（元）

  # 华为云配置
  huawei:
    - name: "huawei-account-1"
      access_key: "你的华为云AccessKey"
      secret_key: "你的华为云SecretKey"
      enabled: true
      balance_threshold: 1000.0  # 余额告警阈值（元）

  # 腾讯云配置
  tencent:
    - name: "tencent-account-1"
      access_key: "你的腾讯云SecretId"
      secret_key: "你的腾讯云SecretKey"
      region: "ap-shanghai"      # 地域，可选，默认ap-shanghai
      enabled: true
      balance_threshold: 10.0    # 余额告警阈值（元）
```

### 多账户配置

每个云平台都支持配置多个账户：

```yaml
cloud:
  tencent:
    - name: "tencent-main-account"
      access_key: "主账户SecretId"
      secret_key: "主账户SecretKey"
      region: "ap-shanghai"
      enabled: true
      balance_threshold: 50.0
      
    - name: "tencent-sub-account"
      access_key: "子账户SecretId"
      secret_key: "子账户SecretKey"
      region: "ap-beijing"
      enabled: true
      balance_threshold: 20.0
```

## 运行

```bash
# 使用默认配置
./cloud_balance_exporter

# 使用自定义配置文件
./cloud_balance_exporter --config /path/to/your/config.yaml

# 指定端口
./cloud_balance_exporter --port 9090
```

## 监控指标

访问 `http://localhost:8081/metrics` 查看导出的指标：

### 余额指标
```
cloud_account_balance{account_id="",account_type="tencent",account_name="tencent-account-1",currency="CNY"} 1.0
```

### 余额状态指标
```
cloud_account_balance_status{account_id="",account_type="tencent",account_name="tencent-account-1"} 1
```

状态说明：
- `1`：正常（余额 ≥ 阈值）
- `0`：异常（余额 < 阈值）
