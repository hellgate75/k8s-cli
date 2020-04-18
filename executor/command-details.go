package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"strings"
)

func (c Executor) details() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "cluster":
		return c.detailsOfCluster()
	case "node":
		return c.detailsOfNode()
	case "instance":
		return c.detailsOfInstance()
	default:
		return errors.New(fmt.Sprintf("Command ontain details, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) detailsOfCluster() error {
	if c.request.ClusterName == "" {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node without cluster name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var index int
	if index = c.internal.IndexOf(c.request.ClusterName); index < 0 {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node, cluster name %s isn't available", c.request.ClusterName))
	} else {
		cl, _ = c.internal.Get(c.request.ClusterName)
		c.print(model.ClusterType{
			Name: cl.Name,
			Slots: cl.TotalSlots(),
			Used: cl.UsedSlots(),
			Free: cl.FreeSlots(),
		})
	}
	return nil
}

func (c Executor) detailsOfNode() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node without name and node name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); ! test {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var nds =[]model.Node{}
		if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
			return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances, cluster name %s has no node named %s", c.request.ClusterName, c.request.NodeName))
		} else {
			nd := nds[0]
			c.print(model.NodeType{
				Name: nd.Name,
				Hostname: nd.Host,
				Slots: nd.Slots,
				Used: nd.UsedSlots(),
				Free: nd.FreeSlots(),
			})
		}
	}
	return nil
}

func (c Executor) detailsOfInstance() error {
	if c.request.ClusterName == "" && c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances without cluster name and node name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); ! test {
		return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var nds =[]model.Node{}
		if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
			return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances, cluster name %s has no node named %s", c.request.ClusterName, c.request.NodeName))
		} else {
			nd := nds[0]
			var inst model.Instance
			if inst, test = nd.GetByName(c.request.Instance); !test {
				return errors.New(fmt.Sprintf("Error, could not ontain details a cluster node instances, cluster name %s in node %s has no instance named %s", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			c.print(model.InstanceType{
				Name: inst.Name,
				NameSpace: inst.Namespace,
				Status: inst.Status,
			})
		}
	}
	return nil
}


