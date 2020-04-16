package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"strings"
)

func (c Executor) show() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "clusters":
		return c.showClusters()
	case "nodes":
		return c.showNodes()
	case "instances":
		return c.showInstances()
	default:
		return errors.New(fmt.Sprintf("Command show, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) showClusters() error {
	if c.request.VerifySlots {
		var names = make([]model.ClusterType, 0)
		for _, cl := range c.internal.Clusters {
			names = append(names, model.ClusterType{
				Name:  cl.Name,
				Slots: cl.TotalSlots(),
				Used:  cl.UsedSlots(),
				Free:  cl.FreeSlots(),
			})
		}
		c.print(model.SuccessTypeResponse{
			Type:    "List",
			SubType: "Cluster",
			Content: names,
		})

	} else {
		var names = make([]string, 0)
		for _, cl := range c.internal.Clusters {
			names = append(names, cl.Name)
		}
		c.print(model.SuccessTypeResponse{
			Type:    "List",
			SubType: "Cluster",
			Content: names,
		})
	}
	return nil
}

func (c Executor) showNodes() error {
	if c.request.ClusterName == "" {
		return errors.New(fmt.Sprintf("Error, could not show a cluster nodes without cluster name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not show a cluster nodes, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); !test {
		return errors.New(fmt.Sprintf("Error, could not show a cluster nodes, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var nodes = make([]model.NodeType, 0)
		for _, nd := range cl.Nodes {
			nodes = append(nodes, model.NodeType{
				Name:     nd.Name,
				Hostname: nd.Host,
				Slots:    nd.Slots,
				Free:     nd.FreeSlots(),
				Used:     nd.UsedSlots(),
			})
		}
		c.print(model.SuccessTypeResponse{
			Type:    "List",
			SubType: "Node",
			Content: nodes,
		})
	}
	return nil
}

func (c Executor) showInstances() error {
	if c.request.ClusterName == "" && c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not show a cluster node instances without cluster name and node name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not show a cluster node instances, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); !test {
		return errors.New(fmt.Sprintf("Error, could not show a cluster node instances, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var nds = []model.Node{}
		if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
			return errors.New(fmt.Sprintf("Error, could not show a cluster node instances, cluster name %s has no node named %s", c.request.ClusterName, c.request.NodeName))
		} else {
			nd := nds[0]
			var nodes = make([]model.InstanceType, 0)
			for _, in := range nd.Instances {
				nodes = append(nodes, model.InstanceType{
					Name:      in.Name,
					NameSpace: in.Namespace,
					Status:    in.Status,
				})
			}
			c.print(model.SuccessTypeResponse{
				Type:    "List",
				SubType: "Instance",
				Content: nodes,
			})

		}
	}
	return nil
}
