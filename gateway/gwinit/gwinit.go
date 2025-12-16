package gwinit

import (
	dao_etcd "github.com/kongshui/danmu/dao/etcd"

	conf "github.com/kongshui/danmu/conf/gateway"

	"github.com/kongshui/danmu/zilog"
)

var (
	Ectd_client = dao_etcd.NewEtcd() //etcd
	Config      = conf.GetConf()     //配置
	Zilog       zilog.LogStruct
	Debug       = false
)

func init() {
	Ectd_client.InitEtcd(Config.Etcd.Addr, Config.Etcd.Username, Config.Etcd.Password)
	Zilog.Init(Config.Logging.LogPath, Config.Logging.Level, Config.Logging.MaxSize, Config.Logging.MaxBackups, Config.Logging.MaxAge, 3600)
	if Config.Logging.Level == "debug" {
		Debug = true
	}
}
