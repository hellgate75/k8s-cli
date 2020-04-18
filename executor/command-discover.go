package executor

import (
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"runtime"
	"strings"
)

const(
	clusterDiscoverTemplate="kubectl --kubeconfig=%s --namespace=kube-system get nodes"
)

func (c Executor) discover() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "nodes":
		return c.discoverNodes()
	default:
		return errors.New(fmt.Sprintf("Command prepare, sub-command: %s is unknown", c.request.SubCommand))
	}
}
func (c Executor) discoverNodes() error {
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
		outFolder := fmt.Sprintf("clusters%c%s", os.PathSeparator, c.request.ClusterName)
		fullOutFolder := fmt.Sprintf("%s%c%s", c.baseFolder, os.PathSeparator, outFolder)
		outPath := fmt.Sprintf("%s%c%s", fullOutFolder, os.PathSeparator, cl.ClusterFile)
		if strings.ToLower(runtime.GOOS) == "windows" {
			outPath = strings.ReplaceAll(outPath, fmt.Sprintf("%c", os.PathSeparator), "/")
		}
		cmd := fmt.Sprintf(clusterDiscoverTemplate, outPath)
		output, err := execute(strings.Split(cmd, " "))
		lines := strings.Split(output, "\n")
		var hosts = make([]string, 0)
		var found = false
		for _, l := range lines {
			if ! strings.Contains(strings.ToLower(l), "name ") && l != "" && strings.Contains(strings.ToLower(l), "ready") {
				host := strings.Split(l, " ")[0]
				if "" != host {
					if ! cl.ContainsHost(host) {
						prefix:=fmt.Sprintf("%s-Node", cl.Name)
						idx := 1
						for cl.Contains(fmt.Sprintf("%s-%v", prefix, idx)) {
							idx += 1
						}
						n := model.NewNode(fmt.Sprintf("%s-%v", prefix, idx), host, c.request.NodeSlots)
						cl.Nodes = append(cl.Nodes, n)
						found = true
						hosts = append(hosts, host)
					}
				}
			}
		}
		if found {
			state := c.internal.UpdateCluster(cl)
			if !state {
				return errors.New(fmt.Sprintf("Couldn't update cluster %s in cluster list", c.request.ClusterName))
			}
			err = c.commit()
			if err != nil {
				return err
			}
		}
		status := "Created"
		if len(hosts) == 0 {
			status = "Not Created"
		}
		c.print(model.SuccessResponse{
			Command: c.request.Command,
			Subject: c.request.SubCommand,
			Status:  status,
			Message: fmt.Sprintf("Cluster %s has created following nodes: %v", cl.Name, hosts),
		})
	}
	return nil
}
