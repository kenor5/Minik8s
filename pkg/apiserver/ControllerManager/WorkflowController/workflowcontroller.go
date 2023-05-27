package workflowcontroller

import (
	"encoding/json"
	"errors"
	"minik8s/entity"
	"minik8s/tools/etcdctl"
	"minik8s/tools/log"
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

func SelectChoice(Number int64, choices []entity.Choice) string {
	for _, choice := range choices {
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