package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"os/exec"
	"runtime"
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
const(
	clusterCheckTemplate="kubectl --kubeconfig=%s --namespace=kube-system cluster-info"
	nodeCheckTemplate="kubectl --kubeconfig=%s --namespace=kube-system get nodes %s --template  --template {{.metadata.name}}"
	nsCheckTemplate="kubectl --kubeconfig=%s --namespace=%s get ns %s  -o jsonpath=\"{.status.phase}\""
	nsCheckTemplate2="kubectl --kubeconfig=%s --namespace=%s get pods  -o jsonpath=\"{.items[*].metadata.name}\""
)
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
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.request.ClusterName)
	fullOutFolder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(clusterCheckTemplate, outPath)
	output, err := execute(strings.Split(cmd, " "))
	status := model.Unhealty
	message := "--"
	if err != nil {
		message=fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), "kubernetes master") {
			status = model.Deployed
		} else {
			message = output
		}
	}
	c.print(model.HealthCheckResponse{
		Type: "Cluster",
		Cluster: cl.Name,
		Node: "--",
		Instance: "--",
		Status: status,
		Message: message,
	})
	return nil
}

func (c Executor) verifyNode() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node without cluster name and node name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	if ! cl.Contains(c.request.NodeName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node, cluster name %s has not node: %s", c.request.ClusterName, c.request.NodeName))
	}
	var nd model.Node
	nd = cl.GetByName(c.request.NodeName)[0]

	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.request.ClusterName)
	fullOutFolder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(nodeCheckTemplate, outPath, nd.Host)
	output, err := execute(strings.Split(cmd, " "))
	status := model.Unhealty
	message := "--"
	if err != nil {
		message=fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), nd.Host) {
			status = model.Deployed
		} else {
			message = output
		}
	}
	c.print(model.HealthCheckResponse{
		Type: "Node",
		Cluster: cl.Name,
		Node: nd.Name,
		Instance: "--",
		Status: status,
		Message: message,
	})
	return nil
}

func (c Executor) verifyInstance() error {
	if c.request.ClusterName == "" || c.request.NodeName == "" || c.request.Instance == "" {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance without cluster name, node name and instance name information"))
	}
	if ! c.internal.Contains(c.request.ClusterName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s doesn't exists", c.request.ClusterName))
	}
	var cl model.Cluster
	cl, _ = c.internal.Get(c.request.ClusterName)
	if ! cl.Contains(c.request.NodeName) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s has not node: %s", c.request.ClusterName, c.request.NodeName))
	}
	var nd model.Node
	nd = cl.GetByName(c.request.NodeName)[0]
	if ! nd.Contains(c.request.Instance) {
		return errors.New(fmt.Sprintf("Error, could not verify a cluster node instance, cluster name %s on node: %s has not instance named: %s ", c.request.ClusterName, c.request.NodeName, c.request.Instance))
	}
	inst, _ := nd.GetByName(c.request.Instance)
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.request.ClusterName)
	fullOutFolder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(nsCheckTemplate, outPath, inst.Namespace, inst.Namespace)
	output, err := execute(strings.Split(cmd, " "))
	status := model.Unhealty
	message := "--"
	if err != nil {
		message=fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), "active"){
			cmd = fmt.Sprintf(nsCheckTemplate2, outPath, inst.Namespace)
			output, err = execute(strings.Split(cmd, " "))
			if err != nil {
				message=fmt.Sprintf("%v", err)
			} else {
				if "" != output {
					status = model.Deployed
				} else {
					message = output
				}
			}
		} else {
			message = output
		}

	}
	c.print(model.HealthCheckResponse{
		Type: "Instance",
		Cluster: cl.Name,
		Node: nd.Name,
		Instance: inst.Name,
		Status: status,
		Message: message,
	})
	return nil
}
