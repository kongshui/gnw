package message

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/kongshui/gnw/node/nodeinit"

	"github.com/kongshui/danmu/model/pmsg"

	"github.com/shirou/gopsutil/v3/cpu"
)

func loadMessageHandler() {
	t := time.NewTicker(3 * time.Second)
	for {
		var sendLoad []byte
		<-t.C
		if runtime.GOOS == "windows" {
			cpuloads, err := cpu.Percent(0, false)
			if err != nil {
				ziLog.Error(fmt.Sprintf("获取负载失败 %v", err), debug)
			}
			var load float64
			for _, f := range cpuloads {
				load += f
			}
			sendLoad = []byte(strconv.FormatFloat(load/float64(len(cpuloads)), 'f', 2, 64))
		} else {
			load, err := nodeinit.GetLoadAvg()
			if err != nil {
				ziLog.Error(fmt.Sprintf("获取负载失败 %v", err), debug)
			}
			sendLoad = []byte(strconv.FormatFloat(load.Load1*100, 'f', 2, 64))
		}
		for _, v := range MessageMap.GetConnAllMap() {
			if err := sendMessage("", v, pmsg.MessageId_NodeLoad, sendLoad, ""); err != nil {
				ziLog.Error(fmt.Sprintf("发送负载失败 %v", err), debug)
			}
		}
	}

}
