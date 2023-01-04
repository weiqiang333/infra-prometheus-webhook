# infra-prometheus-webhook

infra-prometheus-webhook 是一个 prometheus + alertmanager + webhook 警报通知接收器钩子

集合: prometheus-webhook-dingtalk; prometheus-webhook-yunpian; prometheus-webhook-weixin ;prometheus-webhook-telegram。

什么是 prometheus? 请看这里 [`Prometheus`](https://prometheus.io/docs/introduction/overview/#what-is-prometheus)

它实现了 [`DingTalk`](https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq) 群机器人接收器
和 [`yunpian`](https://github.com/yunpian/yunpian-go-sdk) 云片短信+语音接收器
以及 [Weixin robot](https://developer.work.weixin.qq.com/document/path/91770)
和 [Telegram Bot](https://core.telegram.org/bots/api).

[English](README.md)

## 目标和状态
- 该 webhook 的主要目标是：
```text
    1. 接收来自 prometheus -> alertmanager -> HTTP POST 请求.
    2. 根据 Post 的参数及数据来判断发送不同等级、不同媒介、聚合的报警信息.
```
- 支持的媒介
```text
1. 钉钉消息: 使用钉钉群机器人来向处理故障的人知晓它
2. 云片短信+语音,yunpian sms/voice
    特定规则: 只有00:00至08:00点才会进行电话报警. 
        电话是特定解决其它媒介都未唤醒值班人员来解决问题.(因为它的费用或者干扰度更高)
3. 企业微信
4. Telegram 机器人
```

# 用法

初始化并使用 infra-prometheus-webhook, 
使用 systemd 管理你的 infra-prometheus-webhook 服务.

- 安装
```bash
version=v2.1
wget https://github.com/weiqiang333/infra-prometheus-webhook/releases/download/${version}/infra-prometheus-webhook-linux-amd64-${version}.tar.gz
mkdir -p /usr/local/infra-prometheus-webhook/log
tar -zxf infra-prometheus-webhook-linux-amd64-${version}.tar.gz -C /usr/local/infra-prometheus-webhook
chmod +x /usr/local/infra-prometheus-webhook/infra-prometheus-webhook
# systemd manager serivce
cp /usr/local/infra-prometheus-webhook/configs/systemd/infra-prometheus-webhook.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now infra-prometheus-webhook
systemctl status infra-prometheus-webhook
```
- API
[api test](./docs/api_test/api_test.md)
```text
/
    # health check
/-/reload
    # reload config file
/alerts/dingtalk/:priority
/alerts/phonecall/:role
/alerts/yunpian/:sendtype/:priority
    # sendtype: sms/voice
/alerts/weixin/:priority
/alerts/telegram/:priority
```

## configs
注意初始化你的配置文件 (configs/production.yaml)


## Use

- infra-prometheus-webhook -h
```
      --check                   check's cron: Used to check the infrastructure
      --config string           config file (default "configs/production.yaml")
      --listen_address string   server listen address. (default "0.0.0.0:8099")
```

# 使用案例
prometheus and alertmanager and infra-prometheus-webhook Cases used together
- prometheus's alerting_rules.yml file
```yaml
groups:
- name: Instances
  rules:
  - alert: Instance Is Down
    expr: up == 0
    for: 3m
    labels:
      priority: P0
    annotations:
      description: '{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 3 minutes.'
      summary: 'Instance is down. It could happen: 1. node_exporter service run is failed; 2. A critical error has occurred on the Instance, cause instance is down;'
- name: Disk
  rules:
  - alert: Disk Is Pressure
    expr: (1-  (node_filesystem_free_bytes{fstype=~"ext3|ext4|xfs"} / node_filesystem_size_bytes{fstype=~"ext3|ext4|xfs"}) ) * 100 > 80
    for: 3m
    labels:
      priority: P1
    annotations:
      description: '{{ $labels.instance }} of job {{ $labels.job }} in mountpoint {{ $labels.mountpoint }} Disk utilization up to 80%.'
      summary: 'Instance Disk Insufficient available resources utilization up to 80%'
```

- alertmanager's alertmanager.yml file
```yaml
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 3h
  receiver: 'default-receiver_web.hook'
  routes:
  - receiver: 'web.hook-P0'
    group_by: [priority, alertname]
    group_wait: 10s
    repeat_interval: 1h
    matchers:
    - priority="P0"
  - receiver: 'web.hook-P1'
    group_by: [priority, alertname]
    group_wait: 20s
    repeat_interval: 3h
    matchers:
    - priority="P1"
receivers:
  - name: 'default-receiver_web.hook'
    webhook_configs:
      - url: 'http://127.0.0.1:8099/alerts/weixin/p0'
  - name: 'web.hook-P0'
    webhook_configs:
      - url: 'http://127.0.0.1:8099/alerts/weixin/p0'
  - name: 'web.hook-P1'
    webhook_configs:
      - url: 'http://127.0.0.1:8099/alerts/weixin/p1'
```
- infra-prometheus-webhook's Weixin messages
```text
状态: PROBLEM

等级: P1

告警: Disk Is Pressure

Item values: 

1. 172.16.82.100:9101 of job gitlab in mountpoint /var/opt Disk utilization up to 80%.
2. 172.16.9.46:9100 of job jenkins_master in mountpoint / Disk utilization up to 80%.

故障: Instance Disk Insufficient available resources utilization up to 80%
```