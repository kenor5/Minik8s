package etcdctl

import (
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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
	start, err := Start(".")
	if err != nil {
		fmt.Println("start error")
		return
	}

	Put(start, "k1", "v1")
	get, err := Get(start, "k1")
	if err != nil {

		fmt.Println("get error")
	}

	fmt.Println(get.Kvs)
	defer func(start *clientv3.Client) {
		err := start.Close()
		if err != nil {

		}
	}(start)
}

func TestEtcdGetWithPrefix(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name string
		args args
		//want    *clientv3.GetResponse
		wantErr bool
	}{
		{
			name: "test 前缀查询",
			args: args{k: "Pod/nginx-deployment"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EtcdGetWithPrefix(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("EtcdGetWithPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("EtcdGetWithPrefix() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
