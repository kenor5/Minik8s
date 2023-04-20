package entity

/*
	reference to
	https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/workload-resources/pod-v1/#PodSpec
	https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/common-definitions/object-meta/#ObjectMeta

	``的用法： https://blog.csdn.net/zsf211/article/details/106534361
*/

type ObjectMeta struct {
	// Name 在命名空间内必须是唯一的。创建资源时需要
	Name string `json:"name" yaml:"name"`

	// Namespace 定义了一个值空间，其中每个名称必须唯一。“omitempty”表示可为空，默认为default命名空间
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// 可用于组织和分类（确定范围和选择）对象的字符串键和值的映射。
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Uid 是该对象在时间和空间上的唯一值。它通常由服务器在成功创建资源时生成，并且不允许使用 PUT 操作更改。
	Uid string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

// PodTemplate /*------------------------------------------------------*/
type PodTemplate struct {
	Metadata ObjectMeta      `json:"metadata" yaml:"metadata"`
	Template PodTemplateSpec `json:"template" yaml:"template"`
}

type PodTemplateSpec struct {
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     PodSpec    `json:"spec" yaml:"spec"`
}


type ReplicaSet struct {
	Metadata ObjectMeta       `json:"metadata" yaml:"metadata"`
	Spec     ReplicaSetSpec   `json:"spec" yaml:"spec"`
	Status   ReplicaSetStatus `json:"status" yaml:"status"`
}