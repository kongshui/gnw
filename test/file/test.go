package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	dao_etcd "github.com/kongshui/danmu/dao/etcd"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	etcdClient = dao_etcd.NewEtcd()
)

func init() {
	etcdClient.InitEtcd([]string{"127.0.0.1:2379"}, "root", "123456")
}

// ConfigFileSend 发送配置文件
func ConfigFileSend(uidList []string) error {
	filepath.Walk("./test.log", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// 如果大于512k则分片读取发送
		if info.Size() > 512*1024 {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("ConfigFileSend 打开文件失败：%v", err)
			}
			defer file.Close()
			buf := make([]byte, 512*1024)
			labelCount := 0
			for {
				n, err := file.Read(buf)
				if err != nil && err.Error() != io.EOF.Error() {
					return fmt.Errorf("ConfigFileSend 读取文件失败：%v", err)
				}
				if n == 0 {
					fmt.Println("文件名：", info.Name(), "文件大小：", info.Size(), "文件内容：", string(buf[:n]))
					break
				}
				// 发送文件内容
				labelCount++
				fmt.Println("文件名111：", info.Name(), "文件大小：", info.Size(), "文件内容：", string(buf[:n]))
				sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				var isOk bool
				res := etcdClient.Client.Watch(sCtx, fmt.Sprintf("/config/%s/%s/%s/%s", "test", uidList[0], info.Name(),
					strconv.FormatInt(int64(labelCount), 10)))
				fmt.Println(33333333, fmt.Sprintf("/config/%s/%s/%s/%s", "test", uidList[0], info.Name(),
					strconv.FormatInt(int64(labelCount), 10)))
				go func(ctx context.Context) {
					t := time.NewTicker(1 * time.Second)
					defer t.Stop()
					count := 0

					for {
						select {
						case <-ctx.Done():
							fmt.Println("上下文取消，等待确认超时")
							return
						case <-t.C:
							count++
							fmt.Println("等待确认", count)
							if count >= 9 {
								return
							}
						}
					}
				}(sCtx)
				for event := range res {
					for _, ev := range event.Events {
						if ev.Type == mvccpb.PUT {
							isOk = true
						}
					}
					if isOk {
						break
					}
				}
				cancel()
				if !isOk {
					return fmt.Errorf("ConfigFileSend 等待确认超时")
				}
				fmt.Println(222222222222)
			}
		} else {
			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil && err.Error() != io.EOF.Error() {
				return fmt.Errorf("ConfigFileSend 读取文件失败：%v", err)
			}
			fmt.Println("文件名：", info.Name(), "文件大小：", info.Size(), "文件内容：", string(content))
			fmt.Println(111111111111)
		}
		return nil
	})
	return nil
}

func main() {
	etcdClient.Client.Delete(context.Background(), "/config/test/666/"+"test.log", clientv3.WithPrevKV())
	go func() {
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/1", "1")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/2", "2")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/3", "3")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/4", "4")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/5", "5")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/6", "6")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/7", "7")
		time.Sleep(3 * time.Second)
		etcdClient.Client.Put(context.Background(), "/config/test/666/"+"test.log/8", "8")
	}()
	if err := ConfigFileSend([]string{"666"}); err != nil {
		fmt.Println("ConfigFileSend 发送配置文件失败：", err)
	}

}
