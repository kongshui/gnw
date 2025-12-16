package nodeinit

import (
	"fmt"

	"github.com/google/uuid"
	dao_etcd "github.com/kongshui/danmu/dao/etcd"

	conf "github.com/kongshui/danmu/conf/node"

	"github.com/kongshui/danmu/zilog"
)

var (
	Ectd_client = dao_etcd.NewEtcd() //etcd
	Config      = conf.GetConf()     //配置
	Zilog       zilog.LogStruct
	Project     = Config.Project
	NodeUuid    = uuid.New() //uuid
)

func init() {
	Ectd_client.InitEtcd(Config.Etcd.Addr, Config.Etcd.Username, Config.Etcd.Password)
	Zilog.Init(Config.Logging.LogPath, Config.Logging.Level, Config.Logging.MaxSize, Config.Logging.MaxBackups, Config.Logging.MaxAge, 3600)
	Zilog.Info(fmt.Sprintf("node init success, uuid: %v", NodeUuid.String()), true)
}
