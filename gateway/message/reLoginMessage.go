package message

import (
	"fmt"
	"path"
	"time"

	msginterface "github.com/kongshui/gnw/msg/msginterface"

	"github.com/kongshui/danmu/common"

	pmsg "github.com/kongshui/danmu/model/pmsg"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

// 重新登录消息
func reLoginMessageHandler(uidStr string, msgConn msginterface.MsgConn, data []byte, extra string) {
	d := &pmsg.ReloginMessage{}
	if err := proto.Unmarshal(data, d); err != nil {
		ziLog.Error(fmt.Sprintf("reLoginMessageHandler UnmarshalErr, 地址：%v, err：%v", msgConn.RemoteAddr().String(), err), debug)
		return
	}
	defer recoverRelogin(msgConn, d, extra)
	// 判断uidStr是否为空
	if uidStr == "" {
		msgConn.SetOnline(false)
		msgConn.Close()
		msgConn.Cancel()
		panic("ReLoginMessageHandler uidStr为空")
	}
	ziLog.Info(fmt.Sprintf("ReLoginMessageHandler 重新登录消息, 地址：%v, 消息：%v", msgConn.RemoteAddr().String(), d.String()), debug)
	// 注册uuid
	// if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join(common.Uuid_Register_key, msgConn.GetUuid()), commconet.Uuid.String(), clientv3.WithLease(msgConn.GetLease())); err != nil {
	// 	log.Println("ReLoginMessageHandler Uuid_Register_key 发送消息至etcd失败:", msgConn.RemoteAddr().String(), ", err： ", err, d.String())
	// 	msgConn.SetOnline(false)
	// 	msgConn.Close()
	// 	msgConn.Cancel()
	// 	return
	// }
	// 注册openId
	if d.OpenId != "" {
		// resp, err := ectdClient.Client.Get(msgConn.GetCtx(), path.Join("/", config.Project, common.OpenId_Register_Uid_key, d.GetOpenId()), clientv3.WithCountOnly())
		// if resp.Count > 0 || err != nil {
		// 	msgConn.SetOnline(false)
		// 	msgConn.Close()
		// 	msgConn.Cancel()
		// 	panic("OpenId_Register_Uid_key 已绑定")
		// }
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.OpenId_Register_Uid_key, d.GetOpenId()), msgConn.GetUuid(), clientv3.WithLease(msgConn.GetLease())); err != nil {

			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("OpenId_Register_Uid_key 发送消息至etcd失败")
		}
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.Uid_Register_OpenId_key, msgConn.GetUuid()), d.GetOpenId(), clientv3.WithLease(msgConn.GetLease())); err != nil {

			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("Uid_Register_OpenId_key 发送消息至etcd失败")
		}
	}
	// 注册userId
	if d.UserId != "" {
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.UserId_Register_Uid_key, d.GetOpenId()), msgConn.GetUuid(), clientv3.WithLease(msgConn.GetLease())); err != nil {

			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("UserId_Register_Uid_key 发送消息至etcd失败")
		}
		// if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.Uid_Register_OpenId_key, msgConn.GetUuid()), d.GetOpenId(), clientv3.WithLease(msgConn.GetLease())); err != nil {
		// 	log.Println("ReLoginMessageHandler UserId_Register_key 发送消息至etcd失败:", msgConn.RemoteAddr().String(), ", err： ", err, d.String())
		// 	msgConn.SetOnline(false)
		// 	msgConn.Close()
		// 	msgConn.Cancel()
		// 	return
		// }
		msgConn.SetUserId(d.UserId)
	}
	// roomid 和 openid 绑定
	if d.GetRoomId() != "" && d.GetOpenId() != "" {
		// resp, err := ectdClient.Client.Get(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_OpenId_Register_key, d.GetOpenId()), clientv3.WithCountOnly())
		// if resp.Count > 0 || err != nil {
		// 	msgConn.SetOnline(false)
		// 	msgConn.Close()
		// 	msgConn.Cancel()
		// 	panic("RoomId_OpenId_Register_key 已绑定")
		// }
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_OpenId_Register_key, d.GetOpenId()), d.GetRoomId(), clientv3.WithLease(msgConn.GetLease())); err != nil {
			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("RoomId_OpenId_Register_key 发送消息至etcd失败")
		}
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_OpenId_Register_key, d.GetRoomId()), d.GetOpenId(), clientv3.WithLease(msgConn.GetLease())); err != nil {

			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("RoomId_OpenId_Register_key 发送消息至etcd失败")
		}
	}
	// 注册roomId
	if d.RoomId != "" {
		// resp, err := ectdClient.Client.Get(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_Register_Uid_key, d.GetRoomId()), clientv3.WithCountOnly())
		// if resp.Count > 0 || err != nil {
		// 	msgConn.SetOnline(false)
		// 	msgConn.Close()
		// 	msgConn.Cancel()
		// 	panic("RoomId_Register_Uid_key 已绑定")
		// }
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_Register_Uid_key, d.RoomId), msgConn.GetUuid(), clientv3.WithLease(msgConn.GetLease())); err != nil {
			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("RoomId_Register_Uid_key 发送消息至etcd失败")
		}
		if _, err := ectdClient.Client.Put(msgConn.GetCtx(), path.Join("/", config.Project, common.Uid_Register_RoomId_key, msgConn.GetUuid()), d.RoomId, clientv3.WithLease(msgConn.GetLease())); err != nil {

			msgConn.SetOnline(false)
			msgConn.Close()
			msgConn.Cancel()
			panic("Uid_Register_RoomId_key 发送消息至etcd失败")
		}
		msgConn.SetRoomId(d.RoomId)
		sData := &pmsg.MessageBody{MsgId: pmsg.MessageId_ReLoginAck, MessageData: data, Timestamp: time.Now().UnixMilli(), Uuid: msgConn.GetUuid(), Extra: extra}
		sDataByte, _ := proto.Marshal(sData)
		if err := sendMessageToClient(msgConn.GetUuid(), pmsg.MessageId_ReLoginAck, sDataByte, extra); err != nil {
			ziLog.Error(fmt.Sprintf("ReLoginMessageHandler 发送消息至client失败, id:%v err:%v data:%v", msgConn.GetUuid(), err, sData), debug)
			panic("发送ReLoginAck消息失败")
		}
		reData := &pmsg.Reconnect{}
		reData.RoomId = d.GetRoomId()
		reData.UserId = d.GetOpenId()
		reDataByte, err := proto.Marshal(reData)
		if err != nil {
			ziLog.Error(fmt.Sprintf("SendDisConnectMsg marshal err:%v", err.Error()), debug)
			panic("发送Reconnect消息失败, marshal err:" + err.Error())
		}
		dData := &pmsg.MessageBody{}
		dData.MessageData = reDataByte
		dData.MsgId = pmsg.MessageId_ReLogin
		dData.Uuid = msgConn.GetUuid()
		sDataByte, err = proto.Marshal(dData)
		if err != nil {
			ziLog.Error(fmt.Sprintf("SendDisConnectMsg MessageBody marshal err:%v", err.Error()), debug)
			panic("发送Reconnect消息失败 , MessageBody marshal err:" + err.Error())
		}
		if err := sendMessageToNode(msgConn, pmsg.MessageId_Forward, sDataByte, extra); err != nil {
			ziLog.Error(fmt.Sprintf("转发消息toNode失败, id:%v err:%v data:%v", msgConn.GetUuid(), err, sData), debug)
			panic("转发Reconnect消息失败, sendMessageToNode err:" + err.Error())
		}

	}
	// 注册groupId
	// if d.GroupId != "" {
	// 	if _, err := ectdClient.Client.Put(msgConn.GetCtx(), config.Etcd.GroupIdKey+d.GroupId, msgConn.GetUuid(), clientv3.WithLease(msgConn.GetLease())); err != nil {
	// 		log.Println("ReLoginMessageHandler 发送消息至etcd失败, :", msgConn.RemoteAddr().String(), ", err： ", err, d.String())
	// 	}
	// 	if _, err := ectdClient.Client.Put(msgConn.GetCtx(), config.Etcd.GroupIdKey+msgConn.GetUuid(), d.GroupId, clientv3.WithLease(msgConn.GetLease())); err != nil {
	// 		log.Println("ReLoginMessageHandler 发送消息至etcd失败, :", msgConn.RemoteAddr().String(), ", err： ", err, d.String())
	// 	}
	// 	msgConn.SetGroupId(d.GroupId)
	// }
}

func recoverRelogin(msgConn msginterface.MsgConn, data *pmsg.ReloginMessage, extra string) {
	if err := recover(); err != nil {
		ziLog.Error(fmt.Sprintf("reLoginMessageHandler recoverRelogin err, 地址：%v, err：%v, data:%v", msgConn.RemoteAddr().String(), err, data.String()), debug)
		// 回滚删除etcd注册
		if fmt.Sprintf("%v", err) != "UnmarshalErr" {
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.OpenId_Register_Uid_key, data.GetOpenId()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.Uid_Register_OpenId_key, msgConn.GetUuid()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.UserId_Register_Uid_key, data.GetOpenId()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_OpenId_Register_key, data.GetOpenId()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_OpenId_Register_key, data.GetRoomId()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.RoomId_Register_Uid_key, data.GetRoomId()))
			ectdClient.Client.Delete(msgConn.GetCtx(), path.Join("/", config.Project, common.Uid_Register_RoomId_key, msgConn.GetUuid()))
		}
		// 后续补充处理
		dataByte, _ := proto.Marshal(data)
		sData := &pmsg.MessageBody{MsgId: pmsg.MessageId_ReLoginAck, MessageData: dataByte, Timestamp: time.Now().UnixMilli(), Uuid: msgConn.GetUuid(), Extra: extra}
		sDataByte, _ := proto.Marshal(sData)
		if err := sendMessageToClient(msgConn.GetUuid(), pmsg.MessageId_ReLoginAck, sDataByte, extra); err != nil {
			ziLog.Error(fmt.Sprintf("ReLoginMessageHandler 发送消息至client失败, id:%v, err:%v, data:%v", msgConn.GetUuid(), err, sData), debug)
		}
	}
}
