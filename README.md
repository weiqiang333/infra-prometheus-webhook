# infra-prometheus-webhook

This is a prometheus alert notification receiver

What is prometheus? See here [`prometheus`](https://prometheus.io/docs/introduction/overview/#what-is-prometheus)

Here, the [`DingTalk`](https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq) group robot receiver and the [`yunpian`](https://github.com/yunpian/yunpian-go-sdk) voice receiver are implemented.


# init

Here is the initialization of infra-prometheus-webhook, 
and using systemd to manage infra-prometheus-webhook


# configs

Note to initialize your configuration


# Use

- infra-prometheus-webhook -h
```
Usage:
  infra-prometheus-webhook [command]

Available Commands:
  help        Help about any command
  webhook     prometheus alert receiver webhook

Flags:
      --config string   config file (default is $HOME/.infra-prometheus-webhook.yaml)
  -h, --help            help for infra-prometheus-webhook
```