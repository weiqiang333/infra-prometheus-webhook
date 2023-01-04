# infra-prometheus-webhook

This is a prometheus + alertmanager + webhook alert notification receiver.

Gather: prometheus-webhook-dingtalk; prometheus-webhook-yunpian; prometheus-webhook-weixin ;prometheus-webhook-telegram.

What is prometheus? See here [`prometheus`](https://prometheus.io/docs/introduction/overview/#what-is-prometheus)

Here, the [`DingTalk`](https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq) group robot receiver 
and the [`yunpian`](https://github.com/yunpian/yunpian-go-sdk) Yunpian SMS + voice receiver
and the [Weixin robot](https://developer.work.weixin.qq.com/document/path/91770)
and the [Telegram Bot](https://core.telegram.org/bots/api).

[中文](README-cn.md)

## Goals and status
- The main goals of this webhook are:
```text
    1. Receive from prometheus -> alertmanager -> HTTP POST request.
    2. According to the parameters and data of the Post, it is judged to send different levels, different media, and aggregated alarm information.
```
- supported media
```text
1. DingTalk News: Use the DingTalk group bot to inform the people who deal with it
2. YunPian sms/voice ()
    Specific rules: Call the police only from 00:00 to 08:00.
        Telephone is a specific solution and other media have not awakened the on-duty personnel to solve the problem. (Because it is more expensive or intrusive)
3. Enterprise WeChat
4. Telegram Bot
```

# Usage

Here is the initialization of infra-prometheus-webhook, 
and using systemd to manage infra-prometheus-webhook

- install
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

Note to initialize your configuration file (configs/production.yaml)


# Use

- infra-prometheus-webhook -h
```
      --check                   check's cron: Used to check the infrastructure
      --config string           config file (default "configs/production.yaml")
      --listen_address string   server listen address. (default "0.0.0.0:8099")
```

# prometheus and alertmanager and infra-prometheus-webhook Cases used together
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