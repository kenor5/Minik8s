package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func getDocekrIdByname(podname string) error {

	//ctx := context.Background()
	//cli, err := dockerclient.NewClientWithOpts()
	//if err != nil {
	//	panic(err)
	//}
	//defer cli.Close()
	//
	//// 获取运行中的容器列表
	//containers, err := cli.ContainerList(ctx, dockertypes.ContainerListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	////ctx := context.Background()
	////dockerClient, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv)
	////if err != nil {
	////	panic(err)
	////}
	////defer dockerClient.Close()
	//
	////// 设置name参数为过滤器，筛选名称包含"podname"的容器
	////containerListOptions := types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{
	////	Key:   "name",
	////	Value: "*" + podname + "*",
	////})}
	////containers, err := dockerClient.ContainerList(ctx, containerListOptions)
	////
	////if err != nil {
	////	panic(err)
	////}
	////for _, container := range containers {
	////	fmt.Println(container.ID)
	////}
	//
	//// 获取带有关键字的容器名称
	//keyword := podname // 关键字为Podname
	//for _, container := range containers {
	//	for _, containername := range container.Names {
	//		if strings.Contains(string(containername), keyword) {
	//			fmt.Println(containername)
	//		}
	//	}
	//}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// 创建一个 filter 查询条件
	filter := filters.NewArgs()
	filter.Add("name", podname)

	// 获取容器列表
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     false,
		Filters: filter,
	})
	if err != nil {
		panic(err)
	}

	// 遍历容器列表输出 ID
	for _, container := range containers {
		fmt.Println(container.ID)
	}
	return nil
}
