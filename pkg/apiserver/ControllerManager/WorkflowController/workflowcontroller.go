package workflowcontroller

import (
	"encoding/json"
	"errors"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
	"strconv"
	"strings"
)

func GetWorkflow(workflowName string) (*entity.Workflow, error) {
	// 从etcd中拿出对应的workflow
	cli, err := etcdctl.NewClient()
	if err != nil {
		log.PrintE("etcd client connetc error")
	}
	defer cli.Close()

	out, err := etcdctl.Get(cli, "Workflow/"+workflowName)
	if err != nil {
		log.PrintE("No such Workflow!")
		return nil, err
	}

	// 解析Workflow并返回
	workflow := &entity.Workflow{}
    json.Unmarshal(out.Kvs[0].Value, workflow)
	return workflow, nil
}

func GetWorkflowNodeByName(Name string, workflow *entity.Workflow) (*entity.WorkflowNode, error) {
    for _, workflowNode := range workflow.WorkflowNodes {
		if (workflowNode.Name == Name) {
			return &workflowNode, nil
		}
	}
	return nil, errors.New("No such Node!")
} 

func SelectChoice(data string, choices []entity.Choice) string {
	for _, choice := range choices {
		Number := GetVariable(data, choice.Variable)
		match := false
		switch choice.Condition {
		case "NumericEquals":
			match = (Number==choice.Number)
			break
		case "NumericNotEquals":
			match = (Number!=choice.Number)
			break
		case "NumericLessThan":
			match = (Number<choice.Number)
			break
		case "NumericGreaterThan":
			match = (Number>choice.Number)
			break		             
		case "NumericGreaterThanOrEqual":
			match = (Number>=choice.Number)
			break	
		case "NumericLessThanOrEqual":
			match = (Number<=choice.Number)
			break							
		default :
			return "End"
		}
		if match {
			return choice.Next
		} 
	}
	return "End"
}

func GetVariable(data string, variable string) int64 {
	// 查找 "finalGrade" 的起始索引位置
	startIndex := strings.Index(data, variable)
	if startIndex == -1 {
		log.PrintE("未找到 \"finalGrade\"")
		return -1
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
		log.PrintE("转换为int64类型失败:", err)
		return -1
	}

	// 打印结果
	log.Print(finalGradeInt)

	return finalGradeInt
}