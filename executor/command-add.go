package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/io"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"strings"
)

func (c Executor) add() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "cluster":
		return c.addCluster()
	case "node":
		return c.addNode()
	case "instance":
		return c.addInstance()
	default:
		return errors.New(fmt.Sprintf("Command add, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) addCluster() error {
	if c.request.ClusterName == "" || c.request.KubeCtlFile == "" {
		return errors.New(fmt.Sprintf("Error, could not create a cluster without name and kubectl yaml file information"))
	}
	if c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not create a cluster, cluster name %s already in use", c.request.ClusterName))
	}
	if _, err := os.Stat(c.request.KubeCtlFile); err != nil {
		return errors.New(fmt.Sprintf("Error, could not create a cluster, kubectl file %s connet be accessed", c.request.KubeCtlFile))
	}
	outFile := fmt.Sprintf("%s.yaml", io.GetUniqueId())
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.request.ClusterName)
	fullOutFolder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, outFolder)
	if _, err := os.Stat(fullOutFolder); err != nil {
		err = os.MkdirAll(fullOutFolder, 0660)
		if err != nil {
			return err
		}
	}
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, outFile)
	err := io.CopyFile(c.request.KubeCtlFile, outPath)
	if err != nil {
		return err
	}
	cl := model.NewCluster(c.request.ClusterName, outFile, outFolder)
	c.internal.Clusters = append(c.internal.Clusters, cl)
	err = c.commit()
	if err == nil {
		c.print(model.SuccessResponse{
			Command: c.request.Command,
			Subject: c.request.SubCommand,
			Status:  "Created",
			Message: fmt.Sprintf("Cluster %s has been created successfully", c.request.ClusterName),
		})
	}
	return err
}

func (c Executor) addNode() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node without name and node name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); !test {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node, cluster name %s isn't available", c.request.ClusterName))
	} else {
		if cl.Contains(c.request.NodeName) || cl.ContainsHost(c.request.HostName) {
			return errors.New(fmt.Sprintf("Error, could not create a cluster node, cluster name %s already has a node with name %s or hostname %s", c.request.ClusterName, c.request.NodeName, c.request.HostName))
		}
		if c.request.NodeSlots <= 0 {
			return errors.New(fmt.Sprintf("Error, could not create a cluster node, too low value for number of slots %v, minimum is one", c.request.ClusterName, c.request.NodeSlots))
		}
		n := model.NewNode(c.request.NodeName, c.request.HostName, c.request.NodeSlots)
		cl.Nodes = append(cl.Nodes, n)
		state := c.internal.UpdateCluster(cl)
		if !state {
			return errors.New(fmt.Sprintf("Couldn't update cluster %s in cluster list", c.request.ClusterName))
		}
		err := c.commit()
		if err != nil {
			return err
		}
		c.print(model.SuccessResponse{
			Command: c.request.Command,
			Subject: c.request.SubCommand,
			Status:  "Created",
			Message: fmt.Sprintf("Node %s of Cluster %s has been created successfully", c.request.NodeName, c.request.ClusterName),
		})
	}
	return nil
}

func (c Executor) addInstance() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" || c.request.Instance == "" {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node instance without name, node name and instance name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); !test {
		return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s isn't available", c.request.ClusterName))
	} else {
		if !cl.Contains(c.request.NodeName) {
			return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s doesn't exists", c.request.ClusterName))
		} else {
			if !cl.Contains(c.request.NodeName) {
				return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s, with node name %s doesn't exist", c.request.ClusterName, c.request.NodeName))
			} else {
				var nds = make([]model.Node, 0)
				if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
					return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s, with node name %s has following occcurances %v instead of one", c.request.ClusterName, c.request.NodeName, len(nds)))
				} else {
					var nd = nds[0]
					if nd.FreeSlots() <= 0 {
						return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s, with node name %s has no free slot(s)", c.request.ClusterName, c.request.NodeName))
					}
					var found = false
					for _, nd := range cl.Nodes {
						if !found {
							for _, in := range nd.Instances {
								if !found && in.Namespace == c.request.Namespace {
									found = true
								}
							}
						}
					}
					if found {
						return errors.New(fmt.Sprintf("Error, could not create a cluster node instance, cluster name %s, has already a namespace named: %s", c.request.ClusterName, c.request.Namespace))
					}
					i := model.NewInstance(c.request.Instance, c.request.Namespace, c.request.ClusterName, c.request.NodeName)
					nd.Instances = append(nd.Instances, i)
					success := cl.UpdateNode(nd)
					if !success {
						return errors.New(fmt.Sprintf("Couldn't update node named %s in cluster %s nodes list", c.request.NodeName, c.request.ClusterName))
					}
					state := c.internal.UpdateCluster(cl)
					if !state {
						return errors.New(fmt.Sprintf("Couldn't update cluster %s in cluster list", c.request.ClusterName))
					}
					err := c.commit()
					if err != nil {
						return err
					}
					c.print(model.SuccessResponse{
						Command: c.request.Command,
						Subject: c.request.SubCommand,
						Status:  "Created",
						Message: fmt.Sprintf("Instance %s in Node %s of Cluster %s has been created successfully", c.request.Instance, c.request.NodeName, c.request.ClusterName),
					})
				}
			}
		}
	}
	return nil
}
