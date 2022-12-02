#!/bin/bash

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

version=$1
if [ -z $version ]; then
    version=v0.1
fi

go build -o infra-prometheus-webhook infra-prometheus-webhook.go
chmod u+x infra-prometheus-webhook
tar -zcvf infra-prometheus-webhook-linux-amd64-${version}.tar.gz \
  infra-prometheus-webhook configs/production.yaml configs/systemd/ README.md README-cn.md
