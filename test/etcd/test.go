package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go EtcdClient(ctx)
	go EtcdPut()
	time.Sleep(2 * time.Second)
	EtcdGet(ctx)
	// EtcdGet(ctx)

	// etcdClient := GetClient(ctx)
	// lease, _ := etcdClient.Grant(ctx, 3)
	// fmt.Println(etcdClient.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision("/foo"), "=", 0)).Then(clientv3.OpPut("/foo", "bar", clientv3.WithLease(lease.ID))).Commit())
	// t := time.NewTimer(2 * time.Second)
	// <-t.C
	// a, err := etcdClient.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision("/foo"), "=", 0)).Then(clientv3.OpPut("/foo", "bar", clientv3.WithLease(lease.ID))).Commit()
	// fmt.Println(a.Succeeded, err)
	<-ctx.Done()
	time.Sleep(2 * time.Second)
}

func GetClient(ctx context.Context) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // 你的etcd服务器地址
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "123456",
	})
	if err != nil {
		log.Println(err, 22222)
	}
	return cli
}

func EtcdClient(ctx context.Context) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // 你的etcd服务器地址
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "123456",
	})
	if err != nil {
		log.Println(err, 22222)
	}
	defer cli.Close()
	respond := cli.Watch(ctx, "/foo", clientv3.WithPrefix())
	for wresp := range respond {
		fmt.Println(11111111111111111)
		for _, ev := range wresp.Events {
			// var data TestMessage
			// json.Unmarshal(ev.Kv.Value, &data)
			fmt.Println(string(ev.Kv.Key), string(ev.Kv.Value), 222222)
		}
	}
	fmt.Println(3333333333333)
}

func EtcdGet(ctx context.Context) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // 你的etcd服务器地址
		DialTimeout: 10 * time.Second,
		Username:    "root",
		Password:    "123456",
	})
	// cli.Delete(ctx, "/foo", clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}
	defer cli.Close()
	a, _ := cli.Get(ctx, "/foo/0", clientv3.WithCountOnly())
	fmt.Println(a.Count, a.OpResponse(), len(a.Kvs))
	fmt.Println(11111111111111)
	// result, err := cli.Get(ctx, "/foo", clientv3.WithPrefix())
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result, 222222, len(result.Kvs))
	// for _, ev := range result.Kvs {
	// 	fmt.Println("kvs", ev)
	// 	fmt.Println(string(ev.Key), string(ev.Value))
	// }
}

func EtcdPut() {
	cCtx, cCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cCancel()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // 你的etcd服务器地址
		DialTimeout: 10 * time.Second,
		Username:    "root",
		Password:    "123456",
	})
	if err != nil {
		log.Println(err)
	}
	defer cli.Close()
	time.Sleep(1 * time.Second)
	listen := GetListener(cCtx, 3, cli)
	cli.Put(cCtx, "/foo/"+fmt.Sprintf("%d", 0), "aaa", clientv3.WithPrevKV(), clientv3.WithLease(listen.ID))
	keepRespChan := KeepAlive(cCtx, listen.ID, cli)
	go func() {
		for {
			select {
			case <-cCtx.Done():
				log.Println("ctx.Done")
				return
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					log.Println(66666666666666)
					return
				}
			}
		}
	}()

	time.Sleep(1 * time.Second)
	cli.Put(cCtx, "/foo/"+fmt.Sprintf("%d", 0), "bbb", clientv3.WithPrevKV(), clientv3.WithLease(listen.ID))
	cli.Delete(cCtx, "/foo/"+fmt.Sprintf("%d", 0))
	fmt.Println(1111)
	// count := 0
	// for {
	// 	count++
	// 	cli.Put(cCtx, "/foo/"+fmt.Sprintf("%d", count), fmt.Sprintf("%d", count), clientv3.WithPrevKV())
	// 	fmt.Println(count)
	// 	time.Sleep(1 * time.Second)
	// 	cli.Delete(cCtx, "/foo/"+fmt.Sprintf("%d", count))
	// }
	// t := time.NewTicker(5 * time.Second)
	// defer t.Stop()
	// <-t.C
	// count := 0
	// data := TestMessage{}
	// for {
	// 	if count != 0 {
	// 		<-t.C
	// 	}
	// 	if count == 20 {
	// 		cancel()
	// 	}
	// 	data.Id = uint32(count)
	// 	data.Data = fmt.Sprintf("%d", count)
	// 	d, _ := json.Marshal(data)
	// 	_, err = cli.Put(ctx, "/foo/"+fmt.Sprintf("%d", count), string(d), clientv3.WithPrevKV(), clientv3.WithLease(listen.ID))
	// 	if err != nil {
	// 		log.Println(err, 1111)
	// 	}
	// 	count++
	// }
}

// 创建租约
func GetListener(ctx context.Context, ttl int64, client *clientv3.Client) *clientv3.LeaseGrantResponse {
	leaseGrantResp, err := client.Grant(ctx, ttl)
	if err != nil {
		log.Println("创建租约失败", err)
	}
	return leaseGrantResp
}

// 续约
func KeepAlive(ctx context.Context, leaseID clientv3.LeaseID, client *clientv3.Client) <-chan *clientv3.LeaseKeepAliveResponse {
	keepRespChan, err := client.KeepAlive(ctx, leaseID)
	if err != nil {
		log.Println("续约失败", err)
	}
	return keepRespChan
}

type TestMessage struct {
	Id   uint32
	Data string
}
