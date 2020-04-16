package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"strings"
)

func (c Executor) remove() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "cluster":
		return c.removeCluster()
	case "node":
		return c.removeNode()
	case "instance":
		return c.removeInstance()
	default:
		return errors.New(fmt.Sprintf("Command remove, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) removeCluster() error {
	if c.request.ClusterName == "" {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node without cluster name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var index int
	if index = c.internal.IndexOf(c.request.ClusterName); index < 0 {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s isn't available", c.request.ClusterName))
	} else {
		cl, _ = c.internal.Get(c.request.ClusterName)
		folder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, cl.Folder)
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
		done := c.internal.Remove(index)
		if ! done {
			return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s cannot be removed", c.request.ClusterName))
		}
		err = c.commit()
		if err != nil {
			return err
		}
		c.print(model.SuccessResponse{
			Command: c.request.Command,
			Subject: c.request.SubCommand,
			Status: "Removed",
			Message: fmt.Sprintf("Cluster %s has been removed", cl.Name),
		})
	}
	return nil
}

func (c Executor) removeNode() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node without name and node name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); ! test {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s isn't available", c.request.ClusterName))
	} else {
		if ! cl.Contains(c.request.NodeName) {
			return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s has no node with name %s", c.request.ClusterName, c.request.NodeName))
		}
		index := cl.IndexOf(c.request.NodeName)
		if index < 0 {
			return errors.New(fmt.Sprintf("Error, could not remove a cluster node, cluster name %s cannot stat any node with name %s", c.request.ClusterName, c.request.NodeName))
		}
		state := cl.Remove(index)
		if ! state {
			return errors.New(fmt.Sprintf("Cannot remove node %s of cluster %s", c.request.NodeName, c.request.ClusterName))
		}
		state = c.internal.UpdateCluster(cl)
		if ! state {
			return errors.New(fmt.Sprintf("Cannot save cluster node remove changes for node %s of cluster %s", c.request.NodeName, c.request.ClusterName))
		}
		err := c.commit()
		if err != nil {
			return err
		}
		c.print(model.SuccessResponse{
			Command: c.request.Command,
			Subject: c.request.SubCommand,
			Status: "Removed",
			Message: fmt.Sprintf("Successfully removed node %s in cluster %s", c.request.NodeName, c.request.ClusterName),
		})
	}
	return nil
}

func (c Executor) removeInstance() error {
	if c.request.ClusterName == "" && c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances without cluster name and node name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	var test bool
	if cl, test = c.internal.Get(c.request.ClusterName); ! test {
		return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s isn't available", c.request.ClusterName))
	} else {
		var nds =[]model.Node{}
		if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
			return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s has no node named %s", c.request.ClusterName, c.request.NodeName))
		} else {
			nd := nds[0]
			var inst model.Instance
			if inst, test = nd.GetByName(c.request.Instance); !test {
				return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s in node %s has no instance named %s", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			path := fmt.Sprintf("%s%c%s%c%s", c.baseFolder, os.PathSeparator, cl.Folder, os.PathSeparator, inst.File)
			if _, err := os.Stat(path); err == nil{
				err = os.Remove(path)
				if err != nil {
					return err
				}
			}
			folder := fmt.Sprintf("%s%c%s%c%s-%s", c.baseFolder, os.PathSeparator, cl.Folder, os.PathSeparator, inst.Namespace, inst.Name)
			if _, err := os.Stat(folder); err == nil{
				err = os.RemoveAll(folder)
				if err != nil {
					return err
				}
			}
			index := nd.IndexOf(c.request.Instance)
			if index < 0 {
				return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s in node %s cannot stat instance named %s", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			test = nd.Remove(index)
			if ! test {
				return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s in node %s cannot update instance %s node", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			test = cl.UpdateNode(nd)
			if ! test {
				return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s in node %s cannot update instance %s cluster", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			test = c.internal.UpdateCluster(cl)
			if ! test {
				return errors.New(fmt.Sprintf("Error, could not remove a cluster node instances, cluster name %s in node %s cannot update instance %s clusters collection", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			err := c.commit()
			if err != nil {
				return err
			}
			c.print(model.SuccessResponse{
				Command: c.request.Command,
				Subject: c.request.SubCommand,
				Status: "Removed",
				Message: fmt.Sprintf("Instance %s in node %s of cluster %s removed successfully!!", c.request.Instance, c.request.NodeName, c.request.ClusterName),
			})
		}
	}
	return nil
}


