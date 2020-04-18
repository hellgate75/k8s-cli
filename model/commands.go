package model

import (
	"fmt"
	"strings"
)

type Cluster struct {
	Name        string `yaml:"name" json:"name"`
	ClusterFile string `yaml:"kubeConfig" json:"kubeConfig"`
	Folder      string `yaml:"folder" json:"folder"`
	Nodes       []Node `yaml:"nodes" json:"nodes"`
}

func NewCluster(name string, file string, folder string) Cluster {
	return Cluster{
		Name:        name,
		ClusterFile: file,
		Folder:      folder,
		Nodes:       make([]Node, 0),
	}
}

func (c *Cluster) Contains(name string) bool {
	nmLow := strings.ToLower(name)
	for _, n := range c.Nodes {
		if strings.ToLower(n.Name) == nmLow {
			return true
		}
	}
	return false
}

func (c *Cluster) ContainsInstance(name string) bool {
	nmLow := strings.ToLower(name)
	for _, n := range c.Nodes {
		for _, i := range n.Instances {
			if strings.ToLower(i.Name) == nmLow {
				return true
			}
		}
	}
	return false
}

func (c *Cluster) ContainsNameSpace(name string) bool {
	nmLow := strings.ToLower(name)
	for _, n := range c.Nodes {
		for _, i := range n.Instances {
			if strings.ToLower(i.Namespace) == nmLow {
				return true
			}
		}
	}
	return false
}

func (c *Cluster) IndexOf(name string) int {
	nmLow := strings.ToLower(name)
	for index, n := range c.Nodes {
		if strings.ToLower(n.Name) == nmLow {
			return index
		}
	}
	return -1
}

func (c *Cluster) IndexByHost(hostname string) int {
	nmLow := strings.ToLower(hostname)
	for index, n := range c.Nodes {
		if strings.ToLower(n.Host) == nmLow {
			return index
		}
	}
	return -1
}

func (c *Cluster) ContainsHost(host string) bool {
	hsLow := strings.ToLower(host)
	for _, n := range c.Nodes {
		if strings.ToLower(n.Host) == hsLow {
			return true
		}
	}
	return false
}

func (c *Cluster) ContainsFree() bool {
	for _, n := range c.Nodes {
		if n.FreeSlots() > 0 {
			return true
		}
	}
	return false
}

func (c *Cluster) GetByName(name string) []Node {
	var nodes = make([]Node, 0)
	nmLow := strings.ToLower(name)
	for _, n := range c.Nodes {
		if strings.ToLower(n.Name) == nmLow {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (c *Cluster) GetByHost(host string) []Node {
	var nodes = make([]Node, 0)
	hsLow := strings.ToLower(host)
	for _, n := range c.Nodes {
		if strings.ToLower(n.Host) == hsLow {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (c *Cluster) GetSlotsFree() []Node {
	var nodes = make([]Node, 0)
	for _, n := range c.Nodes {
		if n.FreeSlots() > 0 {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (c *Cluster) Remove(index int) bool {
	if index < len(c.Nodes) {
		if index > 0 && index < (len(c.Nodes)-1) {
			var newArr = make([]Node, 0)
			newArr = append(newArr, c.Nodes[:index]...)
			newArr = append(newArr, c.Nodes[index+1:]...)
			c.Nodes = newArr
		} else if index == 0 {
			c.Nodes = c.Nodes[1:]
		} else {
			c.Nodes = c.Nodes[:(len(c.Nodes) - 1)]
		}
		return true
	}
	return false
}

func (c *Cluster) UpdateNode(n Node) bool {

	for index, nd := range c.Nodes {
		if nd.Name == n.Name {
			c.Nodes[index] = n
			return true
		}
	}
	return false
}

func (c *Cluster) FreeSlots() int {
	var value = 0
	for _, nd := range c.Nodes {
		value += nd.FreeSlots()
	}
	return value
}

func (c *Cluster) UsedSlots() int {
	var value = 0
	for _, nd := range c.Nodes {
		value += nd.UsedSlots()
	}
	return value
}

func (c *Cluster) TotalSlots() int {
	var value = 0
	for _, nd := range c.Nodes {
		value += nd.Slots
	}
	return value
}

type ClusterData struct {
	Clusters []Cluster `yaml:"clusters" json:"clusters"`
}

func (cd *ClusterData) UpdateCluster(cl Cluster) bool {

	for index, c := range cd.Clusters {
		if c.Name == cl.Name {
			cd.Clusters[index] = cl
			return true
		}
	}
	return false
}

func (cd *ClusterData) Contains(name string) bool {
	nmLow := strings.ToLower(name)
	for _, cl := range cd.Clusters {
		if strings.ToLower(cl.Name) == nmLow {
			return true
		}
	}
	return false
}

func (cd *ClusterData) Remove(index int) bool {
	if index < len(cd.Clusters) {
		if index > 0 && index < (len(cd.Clusters)-1) {
			var newArr = make([]Cluster, 0)
			newArr = append(newArr, cd.Clusters[:index]...)
			newArr = append(newArr, cd.Clusters[index+1:]...)
			cd.Clusters = newArr
		} else if index == 0 {
			cd.Clusters = cd.Clusters[1:]
		} else {
			cd.Clusters = cd.Clusters[:(len(cd.Clusters) - 1)]
		}
		return true
	}
	return false
}

func (cd *ClusterData) IndexOf(name string) int {
	nmLow := strings.ToLower(name)
	for index, cl := range cd.Clusters {
		if strings.ToLower(cl.Name) == nmLow {
			return index
		}
	}
	return -1
}

func (cd *ClusterData) Get(name string) (Cluster, bool) {
	nmLow := strings.ToLower(name)
	for _, cl := range cd.Clusters {
		if strings.ToLower(cl.Name) == nmLow {
			return cl, true
		}
	}
	return Cluster{}, false
}

type Node struct {
	Name      string     `yaml:"name" json:"name"`
	Host      string     `yaml:"hostname" json:"hostname"`
	Slots     int        `yaml:"slots" json:"slots"`
	Instances []Instance `yaml:"instances" json:"instances"`
}

func NewNode(name string, host string, slots int) Node {
	return Node{
		Name:      name,
		Host:      host,
		Slots:     slots,
		Instances: make([]Instance, 0),
	}
}

func (nd *Node) Contains(name string) bool {
	nmLow := strings.ToLower(name)
	for _, i := range nd.Instances {
		if strings.ToLower(i.Name) == nmLow {
			return true
		}
	}
	return false
}

func (nd *Node) IndexOf(name string) int {
	nmLow := strings.ToLower(name)
	for index, i := range nd.Instances {
		if strings.ToLower(i.Name) == nmLow {
			return index
		}
	}
	return -1
}

func (nd *Node) GetByName(name string) (Instance, bool) {
	nmLow := strings.ToLower(name)
	for _, i := range nd.Instances {
		if strings.ToLower(i.Name) == nmLow {
			return i, true
		}
	}
	return Instance{}, false
}

func (nd *Node) FreeSlots() int {
	return nd.Slots - len(nd.Instances)
}

func (nd *Node) UsedSlots() int {
	return len(nd.Instances)
}

func (nd *Node) GetByNamespace(namespace string) (Instance, bool) {
	nsLow := strings.ToLower(namespace)
	for _, i := range nd.Instances {
		if strings.ToLower(i.Namespace) == nsLow {
			return i, true
		}
	}
	return Instance{}, false
}

func (nd *Node) GetByStatus(status Status) []Instance {
	var out = make([]Instance, 0)
	for _, i := range nd.Instances {
		if i.Status == status {
			out = append(out, i)
		}
	}
	return out
}

func (nd *Node) Remove(index int) bool {
	if index < len(nd.Instances) {
		if index > 0 && index < (len(nd.Instances)-1) {
			var newArr = make([]Instance, 0)
			newArr = append(newArr, nd.Instances[:index]...)
			newArr = append(newArr, nd.Instances[index+1:]...)
			nd.Instances = newArr
		} else if index == 0 {
			nd.Instances = nd.Instances[1:]
		} else {
			nd.Instances = nd.Instances[:(len(nd.Instances) - 1)]
		}
		return true
	}
	return false
}

func (nd *Node) UpdateInstance(i Instance) bool {

	for index, in := range nd.Instances {
		if in.Name == i.Name {
			nd.Instances[index] = i
			return true
		}
	}
	return false
}

type Status string

const (
	Created  Status = "CREATED"
	Ready    Status = "READY"
	Deployed Status = "DEPLOYED"
	Unhealty Status = "UNHEALTHY"
)

type Instance struct {
	Name         string `yaml:"name" json:"name"`
	Namespace    string `yaml:"namespace" json:"namespace"`
	File         string `yaml:"file" json:"file"`
	Status       Status `yaml:"status" json:"status"`
	Index        int    `yaml:"index" json:"index"`
	ClusterIndex int    `yaml:"clusterIndex" json:"clusterIndex"`
}

func NewInstance(name string, namespace string, clusterName string, nodeName string) Instance {
	return Instance{
		Name:      name,
		Namespace: namespace,
		File:      fmt.Sprintf("%s-%s-%s-%s.env", clusterName, nodeName, name, namespace),
		Status:    Created,
	}
}
