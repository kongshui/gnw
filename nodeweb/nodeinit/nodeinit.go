package nodeinit

import (
	"context"

	"github.com/google/uuid"
	dao_etcd "github.com/kongshui/danmu/dao/etcd"

	conf "github.com/kongshui/danmu/conf/nodeweb"

	"github.com/kongshui/danmu/zilog"
)

var (
	Ectd_client *dao_etcd.Etcd //etcd
	Config      *conf.Config   //配置
	Zilog       *zilog.LogStruct
	Project     string
	NodeUuid    = uuid.New() //uuid
)

// init 初始化

func Init(ctx context.Context, etcd *dao_etcd.Etcd, zilog *zilog.LogStruct, config *conf.Config) {
	Ectd_client = etcd
	Zilog = zilog
	Config = config
	Project = config.Project
	// go nodeservice.Listen(ctx)
	// go nodeservice.RegisterToEtcd(ctx)
}
