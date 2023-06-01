// package main

// /**
// * Example of creating and destroying a single container
// * 1. create a new container according to name and image
// * 2. start the container
// * 3. stop the container
// * 4. remove the container
//  */
// func main() {

// 	// containerName := "container1"
// 	// image := "nginx:latest"
// 	// id := docker.CreateContainer(containerName, image)
// 	// fmt.Printf("container %s created\n", id)
// 	// fmt.Printf("It will be started in 1s\n")
// 	// time.Sleep(time.Second * 1)
// 	// docker.StartContainer(id)
// 	// fmt.Printf("after 40s, container %s will be stopped\n", id)
// 	// time.Sleep(time.Second * 40)
// 	// docker.StopContainer(id)
// 	// fmt.Printf("after 3s, container %s will be removed\n", id)
// 	// time.Sleep(time.Second * 3)
// 	// id, err := docker.RemoveContainer(id)
// 	// if err == nil {
// 	// 	fmt.Printf("container %s is removed\n", id)
// 	// }
// }

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main_6() {
	data := `{"finalGrade": 93.0}`
	variable := "finalGrade"
	// 查找 "finalGrade" 的起始索引位置
	startIndex := strings.Index(data, variable)
	if startIndex == -1 {
		fmt.Println("未找到 \"finalGrade\"")
		return
	}

	// 截取 "finalGrade" 后的部分
	substring := data[startIndex+len(`"finalGrade"`):]

	// 移除不需要的字符
	substring = strings.Trim(substring, ` :`)

	// 提取数字字符串
	endIndex := strings.IndexFunc(substring, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if endIndex == -1 {
		endIndex = len(substring)
	}
	numberString := substring[:endIndex]

	// 转换为int64类型
	finalGradeInt, err := strconv.ParseInt(numberString, 10, 64)
	if err != nil {
		fmt.Println("转换为int64类型失败:", err)
		return
	}

	// 打印结果
	fmt.Println(finalGradeInt)
}
