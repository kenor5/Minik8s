package etcdctl

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

// 测试时候，先把本机上的etcd运行起来
func TestGet(t *testing.T) {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()
	res, _ := Get(cli, "key")
	fmt.Println(res)
}

func TestPut(t *testing.T) {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()
	err = Put(cli, "key", "val")
	fmt.Println(err)
}

func TestStart(t *testing.T) {
	start, err := Start("/home/os/minik8s/minik8s/tools/etcdctl")
	if err != nil {
		fmt.Println("start error")
		return
	}

	Put(start, "k1", "v1")
	get, err := Get(start, "k1")
	if err != nil {

		fmt.Println("get error")
	}

	fmt.Println(get)
	defer func(start *clientv3.Client) {
		err := start.Close()
		if err != nil {

		}
	}(start)
}
