## 1. dingtalk
- api: /alerts/dingtalk/:priority
- POST URL: http://127.0.0.1:8099/alerts/dingtalk/p0
- POST Body JSON
```json
{
    "receiver": "web\\.hook-P0",
    "status": "firing",
    "alerts": [
        {
            "status": "firing",
            "labels": {
                "alertname": "Instance Is Down",
                "host_name": "infra0",
                "instance": "localhost:9100",
                "job": "infra0",
                "priority": "P0"
            },
            "annotations": {
                "description": "localhost:9100/infra0 of job infra0 has been down for more than 3 minutes.",
                "summary": "Instance is down. It could happen:\n\\t\\t 1. node_exporter service run is failed\n\\t\\t 2. A critical error has occurred on the Instance, cause instance is down"
            },
            "startsAt": "2022-12-02T10:00:20.381Z",
            "endsAt": "0001-01-01T00:00:00Z",
            "generatorURL": "http://infra0:9090/graph?g0.expr=up+%3D%3D+0&g0.tab=1",
            "fingerprint": "4e4c4c3d872bf6d8"
        }
    ],
    "groupLabels": {
        "alertname": "Instance Is Down",
        "priority": "P0"
    },
    "commonLabels": {
        "alertname": "Instance Is Down",
        "host_name": "infra0",
        "instance": "localhost:9100",
        "job": "infra0",
        "priority": "P0"
    },
    "commonAnnotations": {
        "description": "localhost:9100/infra0 of job infra0 has been down for more than 3 minutes.",
        "summary": "Instance is down. It could happen:\n\\t\\t 1. node\\_exporter service run is failed;\n\\t\\t 2. A critical error has occurred on the Instance, cause instance is down"
    },
    "externalURL": "http://infra0:9093",
    "version": "4",
    "groupKey": "{}/{priority=\"P0\"}:{alertname=\"Instance Is Down\", priority=\"P0\"}",
    "truncatedAlerts": 0
}
```

## 2. yunpian
- api: /alerts/yunpian/:sendtype/:priority
- POST URL
```text
    http://127.0.0.1:8099/alerts/yunpian/voice/p0
    http://127.0.0.1:8099/alerts/yunpian/sms/p0
```

## 3. Weixin robot
- api: /alerts/weixin/:priority
- POST URL: http://127.0.0.1:8099/alerts/weixin/p0

## 4. Telegram Bot
- api: /alerts/telegram/:priority
- POST URL: http://127.0.0.1:8099/alerts/telegram/p0