package etcdctl

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := client.Delete(ctx, k)

	if err != nil {
		return errors.New("etcd get error")
	}
	return nil
}

func Watch(client *clientv3.Client, k string) (clientv3.WatchChan, error) {
	if client == nil {
		return nil, errors.New("client is null")
	}

	return client.Watch(context.Background(), k), nil
}
