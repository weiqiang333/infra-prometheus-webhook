package model

import "time"

/*
{
  "receiver": "P0",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "KubePodContainerRestart",
        "namespace": "ingress-nginx",
        "pod": "nginx-deployment-76bf4969df-64fmc",
        "priority": "P0"
      },
      "annotations": {
        "description": "Pod ingress-nginx/nginx-deployment-76bf4969df-64fmc: 1.2492166370671725 该 pod 最近 10 分钟内重启次数超过 3 次",
        "summary": "Pod restart more than 3 times"
      },
      "startsAt": "2019-05-11T03:54:05.7007992Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus:9090/graph?g0.expr=avg+by%28pod%2C+namespace%29+%28increase%28kube_pod_container_status_restarts_total%5B10m%5D%29%29+%3E+1\u0026g0.tab=1"
    }
  ],
  "groupLabels": {
    "alertname": "KubePodContainerRestart",
    "priority": "P0"
  },
  "commonLabels": {
    "alertname": "KubePodContainerRestart",
    "namespace": "ingress-nginx",
    "pod": "nginx-deployment-76bf4969df-64fmc",
    "priority": "P0"
  },
  "commonAnnotations": {
    "description": "Pod ingress-nginx/nginx-deployment-76bf4969df-64fmc: 1.2492166370671725 该 pod 最近 10 分钟内重启次数超过 3 次",
    "summary": "Pod restart more than 3 times"
  },
  "externalURL": "http://prometheus:9093",
  "version": "4",
  "groupKey": "{}/{priority=\"P0\"}:{alertname=\"KubePodContainerRestart\", priority=\"P0\"}"
}
*/

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}
