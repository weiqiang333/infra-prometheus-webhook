# infra-prometheus-webhook

This is a prometheus alert notification receiver

What is prometheus? See here [`prometheus`](https://prometheus.io/docs/introduction/overview/#what-is-prometheus)

Here, the [`DingTalk`](https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq) group robot receiver and the [`yunpian`](https://github.com/yunpian/yunpian-go-sdk) voice receiver are implemented and weixin.


# init

Here is the initialization of infra-prometheus-webhook, 
and using systemd to manage infra-prometheus-webhook

- install
```bash
    bash init/build.sh
    mkdir -p /apps/infra-prometheus-webhook/{log,configs}
    cp -p infra-prometheus-webhook /usr/sbin/infra-prometheus-webhook
    cp configs/production.yaml /apps/infra-prometheus-webhook/configs/production.yaml
    cp init/systemd/infra-prometheus-webhook.service /etc/systemd/system/
    systemctl daemon-reload
    systemctl start infra-prometheus-webhook
```

# configs

Note to initialize your configuration file (configs/production.yaml)


# Use

- infra-prometheus-webhook -h
```
Usage:
  infra-prometheus-webhook [command]

Available Commands:
  help        Help about any command
  version     View prometheus-webhook version
  webhook     prometheus alert receiver webhook

Flags:
      --config string   config file (default is $HOME/.infra-prometheus-webhook.yaml)
  -h, --help            help for infra-prometheus-webhook
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