package message

import (
	"fmt"
	"strconv"

	msginterface "github.com/kongshui/gnw/msg/msginterface"
)

func getLoadAvgMessage(_ string, msgConn msginterface.MsgConn, data []byte, extra string) {
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		ziLog.Error(fmt.Sprintf("获取负载失败 %v", err), debug)
	}
	msgConn.SetLoad(f)
}
