package model

type ErrorResponse struct {
	Command string `yaml:"command" json:"command"`
	Subject string `yaml:"subject" json:"subject"`
	Status  string `yaml:"status" json:"status"`
	Code    int    `yaml:"code" json:"code"`
	Message string `yaml:"message" json:"message"`
}

type SuccessResponse struct {
	Command string `yaml:"command" json:"command"`
	Subject string `yaml:"subject" json:"subject"`
	Status  string `yaml:"status" json:"status"`
	Message string `yaml:"message" json:"message"`
}

type SuccessTypeResponse struct {
	Type    string      `yaml:"type" json:"type"`
	SubType string      `yaml:"subType" json:"subType"`
	Content interface{} `yaml:"content" json:"content"`
}

type NodeType struct {
	Name     string `yaml:"name" json:"name"`
	Hostname string `yaml:"hostname" json:"hostname"`
	Slots    int    `yaml:"totalSlots" json:"totalSlots"`
	Used     int    `yaml:"usedSlots" json:"usedSlots"`
	Free     int    `yaml:"freeSlots" json:"freeSlots"`
}

type ClusterType struct {
	Name  string `yaml:"name" json:"name"`
	Slots int    `yaml:"totalSlots" json:"totalSlots"`
	Used  int    `yaml:"usedSlots" json:"usedSlots"`
	Free  int    `yaml:"freeSlots" json:"freeSlots"`
}

type InstanceType struct {
	Name      string `yaml:"name" json:"name"`
	NameSpace string `yaml:"namespace" json:"namespace"`
	Status    Status `yaml:"status" json:"status"`
}

type HealthCheckResponse struct {
	Type     string `yaml:"type" json:"type"`
	Cluster  string `yaml:"clusterName" json:"clusterName"`
	Node     string `yaml:"nodeName" json:"nodeName"`
	Instance string `yaml:"instanceName" json:"instanceName"`
	Message  string `yaml:"message" json:"message"`
	Status   Status `yaml:"status" json:"status"`
}

type FreeNodeResponse struct {
	Cluster   string `yaml:"clusterName" json:"clusterName"`
	Node      string `yaml:"nodeName" json:"nodeName"`
	Available bool   `yaml:"available" json:"available"`
}
