package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"strings"
)

func (c Executor) ensure() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "instance":
		return c.ensureInstance()
	default:
		return errors.New(fmt.Sprintf("Command prepare, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) ensureInstance() error {
	if c.request.ClusterName == "" && c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not prepare a cluster node instances without cluster name and node name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not prepare a cluster node instances, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); !test {
		return errors.New(fmt.Sprintf("Error, could not prepare a cluster node instances, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var currNode *model.Node
		for _, nd := range cl.Nodes {
			if nd.FreeSlots() > 0 {
				currNode = &nd
				break
			}
		}
		if currNode != nil {
			c.print(model.FreeNodeResponse{
				Cluster:   cl.Name,
				Node:      currNode.Name,
				Available: true,
			})
		} else {
			c.print(model.FreeNodeResponse{
				Cluster:   cl.Name,
				Node:      "<null>",
				Available: false,
			})
		}
	}
	return nil
}
