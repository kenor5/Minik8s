package etcdctl

import (
	"context"
	"errors"
	"fmt"
	"os"

	// "os"
	// "strings"
	"os/exec"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
	reference to:
	https://www.tizi365.com/archives/574.html
*/

func NewClient() (*clientv3.Client, error) {
	return clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		})
}

func Put(cli *clientv3.Client, k string, v string) error {
	if cli == nil {
		return errors.New("client is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := cli.Put(ctx, k, v)
	cancel()
	if err != nil {
		return errors.New("etcd put error")
	}
	return nil
}

func Get(client *clientv3.Client, k string) (*clientv3.GetResponse, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	ret, err := client.Get(ctx, k)
	cancel()
	if err != nil {
		return nil, errors.New("etcd get error")
	}
	return ret, nil
}

func GetWithPrefix(client *clientv3.Client, k string) (*clientv3.GetResponse, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ret, err := client.Get(ctx, k, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, errors.New("etcd get error")
	}
	return ret, nil
}

func Delete(client *clientv3.Client, k string) error {
	if client == nil {
		return errors.New("client is null")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err1 := client.Delete(ctx, k)
	cancel()

	if err1 != nil {
		return errors.New("etcd delete error")
	}
	return nil
}

func Watch(client *clientv3.Client, k string) (clientv3.WatchChan, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	return client.Watch(context.Background(), k), nil
}

func Start(dirPath string) (*clientv3.Client, error) {
	// 获取到的是调用这个函数的文件路径，进行以下处理
	// calledPath, err := os.Getwd()
	// rootPath := calledPath[:strings.LastIndex(calledPath, "minik8s")]
	

	cmd := exec.Command("etcdctl", "member", "list")
	err := cmd.Run()
	if err != nil {
        fmt.Println("etcd not start")
		fmt.Println("try to start etcd")

			// data_dir := "/home/os/minik8s/minik8s/tools/etcdstore"
		// listen_client_url := "http://127.0.0.1:2379"
		advertise_client_url := "http://localhost:2379"
		
		args := []string {
				// "--data-dir",
				// data_dir,
				// "--listen-client-urls",
				// listen_client_url, 
				"-advertise_client_urls",
				advertise_client_url,
			}
		cmd := exec.Command("etcd", args...)
		
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			fmt.Print("etcd start error : ")
			fmt.Println(err)
		} else {
			fmt.Println("started etcd")
		}
    }else {
		fmt.Println("etcd already started")
	}
	
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		})

	if err != nil {
		fmt.Println("etcd connect error")
		return nil, err
	}else {
		fmt.Println("connetced to etcd server")
	}

	return cli, nil

}
