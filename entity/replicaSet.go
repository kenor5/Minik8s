package entity


type ReplicaSetSpec struct {
	// Selector 是针对 Pod 的标签查询，应与副本计数匹配。
	Selector LabelSelector `json:"selector" yaml:"selector"`
	// Template 是描述 Pod 的一个对象，将在检测到副本不足时创建此对象。
	Template PodTemplateSpec `json:"template" yaml:"template"`
	// Replicas 是预期副本的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 1。
	Replicas int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}

type LabelSelector struct {
	// MatchLabels 映射中的单个 {key,value} 键值对相当于 matchExpressions 的一个元素，其键字段为 key，运算符为 In，values 数组仅包含 value。
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}

type ReplicaSetStatus struct {
	// Replicas 是最近观测到的副本数量。
	Replicas int32 `json:"replicas" yaml:"replicas"`

	// 此副本集可用副本（至少 minReadySeconds 才能就绪）的数量
	AvailableReplicas int32 `json:"availableReplicas,omitempty" yaml:"availableReplicas,omitempty"`
}