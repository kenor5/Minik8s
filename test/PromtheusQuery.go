package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"minik8s/tools/log"
	"time"
)

func PromtheusQuery(query string) {
	// 设置Prometheus主机地址
	cfg := api.Config{
		Address: "http://localhost:9090",
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 创建v1.API对象
	v1api := v1.NewAPI(client)

	// 查询容器CPU使用率
	//query:= "container_cpu_usage_seconds_total{name=\"autoscale\"}"
	//SumQuery := "sum(" + query + ")"
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()
	// 执行查询并获取结果
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		panic(err)
	}
	if len(warnings) > 0 {
		fmt.Println("Warnings:", warnings)
	}

	// 处理查询结果
	if err != nil {
		log.PrintE("fail to get cpu usage from prometheus:", err)
	}
	if len(warnings) > 0 {
		log.PrintE("warnings from prometheus", warnings)
	}
	if result.(model.Vector).Len() == 0 {
		log.PrintE("query is null")
	}

	fmt.Printf("cpu usage: %v\n", float64(result.(model.Vector)[0].Value))
	//matrix := result.(model.Vector)
	//for _, stream := range matrix {
	//	// 打印指标名称
	//	fmt.Println("Metric:", stream.Metric["name"])
	//
	//	// 打印时间序列数据
	//	for _, point := range stream.Values {
	//		fmt.Printf("Timestamp:%v, Value:%v\n", point.Timestamp, point.Value)
	//	}
	//}
}
