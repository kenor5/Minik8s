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
			Endpoints:   []string{"localhost:2379"},
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
