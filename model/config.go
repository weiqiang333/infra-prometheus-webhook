package model

type config struct {
	ListenPort string
	Dingtalk   map[string]string
	Weixin     map[string]string
	PhoneCall  map[string]string
	DataBase   map[string]string
}

var Config config
