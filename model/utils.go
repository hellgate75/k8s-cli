package model

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	clusterCheckTemplate = "kubectl --kubeconfig=%s --namespace=kube-system cluster-info"
	nodeCheckTemplate    = "kubectl --kubeconfig=%s --namespace=kube-system get nodes %s --template  --template {{.metadata.name}}"
	nsCheckTemplate      = "kubectl --kubeconfig=%s --namespace=%s get ns %s  -o jsonpath=\"{.status.phase}\""
	nsCheckTemplate2     = "kubectl --kubeconfig=%s --namespace=%s get pods  -o jsonpath=\"{.items[*].metadata.name}\""
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

func (c *Cluster) HealthCheck(baseFolder string) (Status, string, error) {
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.Name)
	fullOutFolder := fmt.Sprintf("%s%c%s", baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, c.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(clusterCheckTemplate, outPath)
	output, err := execute(strings.Split(cmd, " "))
	status := Unhealty
	message := "--"
	if err != nil {
		message = fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), "kubernetes master") {
			status = Deployed
		} else {
			message = output
		}
	}
	return status, message, err
}

func (nd *Node) HealthCheck(cl Cluster, baseFolder string) (Status, string, error) {
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, cl.Name)
	fullOutFolder := fmt.Sprintf("%s%c%s", baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(nodeCheckTemplate, outPath, nd.Host)
	output, err := execute(strings.Split(cmd, " "))
	status := Unhealty
	message := "--"
	if err != nil {
		message = fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), nd.Host) {
			status = Deployed
		} else {
			message = output
		}
	}
	return status, message, err
}

func (inst *Instance) HealthCheck(cl Cluster, baseFolder string) (Status, string, error) {
	outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, cl.Name)
	fullOutFolder := fmt.Sprintf("%s%c%s", baseFolder, os.PathSeparator, outFolder)
	outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
	if strings.ToLower(runtime.GOOS) == "windows" {
		outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
	}
	cmd := fmt.Sprintf(nsCheckTemplate, outPath, inst.Namespace, inst.Namespace)
	output, err := execute(strings.Split(cmd, " "))
	status := Unhealty
	message := "--"
	if err != nil {
		message = fmt.Sprintf("%v", err)
	} else {
		if strings.Contains(strings.ToLower(output), "active") {
			cmd = fmt.Sprintf(nsCheckTemplate2, outPath, inst.Namespace)
			output, err = execute(strings.Split(cmd, " "))
			if err != nil {
				message = fmt.Sprintf("%v", err)
			} else {
				if "" != output {
					status = Deployed
				} else {
					message = output
				}
			}
		} else {
			message = output
		}

	}
	return status, message, err
}
