package etcdctl

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os/exec"
	"time"
)

/*
	reference to:
	https://www.tizi365.com/archives/574.html
*/

func Put(cli *clientv3.Client, k string, v string) error {
	if cli == nil {
		return errors.New("client is null")
	}

	ctx, err1 := context.WithTimeout(context.Background(), 5*time.Second)
	if err1 != nil {
		return errors.New("context open err")
	}
	_, err := cli.Put(ctx, k, v)

	if err != nil {
		return errors.New("etcd put error")
	}
	return nil
}

func Get(client *clientv3.Client, k string) (*clientv3.GetResponse, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	ctx, err1 := context.WithTimeout(context.Background(), 5*time.Second)
	if err1 != nil {
		return nil, errors.New("context open err")
	}
	ret, err := client.Get(ctx, k)

	if err != nil {
		return nil, errors.New("etcd get error")
	}
	return ret, nil
}

func GetWithPrefix(client *clientv3.Client, k string) (*clientv3.GetResponse, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	ctx, err1 := context.WithTimeout(context.Background(), 5*time.Second)
	if err1 != nil {
		return nil, errors.New("context open err")
	}
	ret, err := client.Get(ctx, k, clientv3.WithPrefix())

	if err != nil {
		return nil, errors.New("etcd get error")
	}
	return ret, nil
}

func Delete(client *clientv3.Client, k string) error {
	if client == nil {
		return errors.New("client is null")
	}

	ctx, err := context.WithTimeout(context.Background(), 5*time.Second)
	if err != nil {
		return errors.New("context open err")
	}

	_, err1 := client.Delete(ctx, k)

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
	cmd := exec.Command("./etcd")
	cmd.Dir = dirPath
	err := cmd.Start()
	if err != nil {
		fmt.Println("etcd already running")
	} else {
		fmt.Println("started etcd")
	}

	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})

	if err != nil {
		fmt.Println("etcd connect error")
		return nil, err
	}

	return cli, nil

}
