package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"os/exec"
	"strings"
)

func (c Executor) verify() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "cluster":
		return c.verifyCluster()
	case "node":
		return c.verifyNode()
	case "instance":
		return c.verifyInstance()
	default:
		return errors.New(fmt.Sprintf("Command verify, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func execute(command []string) (string, error) {
	cmdVal := command[0]
	cmdArgs := command[1:]
	cmd := exec.Command(cmdVal, cmdArgs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", stdoutStderr), nil
}

func (c Executor) verifyCluster() error {
	if c.request.ClusterName == "" {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster without cluster name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	status, message, err := cl.HealthCheck(c.baseFolder)
	if err != nil {
		message = fmt.Sprintf("%v", err)
	}
	c.print(model.HealthCheckResponse{
		Type:     "Cluster",
		Cluster:  cl.Name,
		Node:     "--",
		Instance: "--",
		Status:   status,
		Message:  message,
	})
	return nil
}

func (c Executor) verifyNode() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node without cluster name and node name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	if !cl.Contains(c.request.NodeName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node, cluster name %s has not node: %s", c.request.ClusterName, c.request.NodeName))
	}
	var nd model.Node
	nd = cl.GetByName(c.request.NodeName)[0]
	status, message, err := nd.HealthCheck(cl, c.baseFolder)
	if err != nil {
		message = fmt.Sprintf("%v", err)
	}
	c.print(model.HealthCheckResponse{
		Type:     "Node",
		Cluster:  cl.Name,
		Node:     nd.Name,
		Instance: "--",
		Status:   status,
		Message:  message,
	})
	return nil
}

func (c Executor) verifyInstance() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" || c.request.Instance == "" {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance without cluster name, node name and instance name information"))
	}
	if !c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	if !cl.Contains(c.request.NodeName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s has not node: %s", c.request.ClusterName, c.request.NodeName))
	}
	var nd model.Node
	nd = cl.GetByName(c.request.NodeName)[0]
	if !nd.Contains(c.request.Instance) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s on node: %s has not instance named: %s ", c.request.ClusterName, c.request.NodeName, c.request.Instance))
	}
	inst, _ := nd.GetByName(c.request.Instance)
	status, message, err := inst.HealthCheck(cl, c.baseFolder)
	if err != nil {
		message = fmt.Sprintf("%v", err)
	}
	c.print(model.HealthCheckResponse{
		Type:     "Instance",
		Cluster:  cl.Name,
		Node:     nd.Name,
		Instance: inst.Name,
		Status:   status,
		Message:  message,
	})
	return nil
}
