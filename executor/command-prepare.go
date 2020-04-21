package executor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hellgate75/k8s-cli/model"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func (c Executor) prepare() error {

	switch strings.ToLower(c.request.SubCommand) {
	case "", "instance":
		return c.prepareInstance()
	default:
		return errors.New(fmt.Sprintf("Command prepare, sub-command: %s is unknown", c.request.SubCommand))
	}
}

func (c Executor) prepareInstance() error {
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
		var nds = []model.Node{}
		if nds = cl.GetByName(c.request.NodeName); len(nds) < 1 {
			return errors.New(fmt.Sprintf("Error, could not prepare a cluster node instances, cluster name %s has no node named %s", c.request.ClusterName, c.request.NodeName))
		} else {
			nd := nds[0]
			var inst model.Instance
			if inst, test = nd.GetByName(c.request.Instance); !test {
				return errors.New(fmt.Sprintf("Error, could not prepare a cluster node instances, cluster name %s in node %s has no instance named %s", c.request.ClusterName, c.request.NodeName, c.request.Instance))
			}
			kubefile := fmt.Sprintf("%s%c%s%c%s", c.baseFolder, os.PathSeparator, cl.Folder, os.PathSeparator, cl.ClusterFile)
			path := fmt.Sprintf("%s%c%s%c%s", c.baseFolder, os.PathSeparator, cl.Folder, os.PathSeparator, inst.File)
			if _, err := os.Stat(path); err == nil {
				err = os.Remove(path)
				if err != nil {
					return err
				}
			}
			folder := fmt.Sprintf("%s%c%s%c%s-%s", c.baseFolder, os.PathSeparator, cl.Folder, os.PathSeparator, inst.Namespace, inst.Name)
			if _, err := os.Stat(folder); err != nil {
				err = os.MkdirAll(folder, 0660)
				if err != nil {
					return err
				}
			}
			sep := fmt.Sprintf("%c", os.PathSeparator)
			if strings.ToLower(runtime.GOOS) == "windows" {
				sep = "/"
				folder = strings.ReplaceAll(folder, fmt.Sprintf("%c", os.PathSeparator), sep)
				kubefile = strings.ReplaceAll(kubefile, fmt.Sprintf("%c", os.PathSeparator), sep)
			}
			readyClusterIdx := 1
			if inst.ClusterIndex > 0 {
				readyClusterIdx = inst.ClusterIndex
			} else {
				mp := make(map[int]bool)
				for _, nd0 := range cl.Nodes {
					for _, ints := range nd0.Instances {
						if ints.ClusterIndex > 0 && ints.Name != inst.Name {
							mp[ints.ClusterIndex] = true
						}
					}
				}
				for mp[readyClusterIdx] {
					readyClusterIdx += 1
				}
				inst.ClusterIndex = readyClusterIdx
				_ = nd.UpdateInstance(inst)
				_ = cl.UpdateNode(nd)
				c.internal.UpdateCluster(cl)
				if err := c.commit(); err != nil {
					return err
				}
			}

			readyIdx := 1
			if inst.Index > 0 {
				readyIdx = inst.Index
			} else {
				mp := make(map[int]bool)
				for _, ints := range nd.Instances {
					if ints.Index > 0 && ints.Name != inst.Name {
						mp[ints.Index] = true
					}
				}
				for mp[readyIdx] {
					readyIdx += 1
				}
				inst.Index = readyIdx
				_ = nd.UpdateInstance(inst)
				_ = cl.UpdateNode(nd)
				c.internal.UpdateCluster(cl)
				if err := c.commit(); err != nil {
					return err
				}
			}
			buff := bytes.NewBuffer([]byte{})
			buff.Write([]byte(fmt.Sprintf("FOLDER=\"%s\"\n", folder)))
			buff.Write([]byte(fmt.Sprintf("NAMESPACE=\"%s\"\n", inst.Namespace)))
			buff.Write([]byte(fmt.Sprintf("KUBECONFIG=\"%s\"\n", kubefile)))
			buff.Write([]byte("KUBECTL_BASE=\"--kubeconfig=$KUBECONFIG --namespace=$NAMESPACE\"\n"))
			buff.Write([]byte("HELM_BASE=\"$KUBECTL_BASE --registry-config $HELM_DIR/registry.json --repository-cache $HELM_DIR/repository --repository-config $HELM_DIR/repositories.yaml\"\n"))
			buff.Write([]byte(fmt.Sprintf("HELM_HOME=\"%s%s.helm\"\n", folder, sep)))
			buff.Write([]byte("KUBECTL_CONFIG_FILE=$KUBECONFIG\n"))
			buff.Write([]byte("alias kube-ns=\"kubectl $KUBECTL_BASE\"\n"))
			buff.Write([]byte("alias helm-ns=\"helm --kubeconfig=$KUBECONFIG --namespace=$NAMESPACE --registry-config $HELM_DIR/registry.json --repository-cache $HELM_DIR/repository --repository-config $HELM_DIR/repositories.yaml\"\n"))
			buff.Write([]byte(fmt.Sprintf("CLUSTER_INDEX=%v\n", readyClusterIdx)))
			buff.Write([]byte(fmt.Sprintf("NODE_INDEX=%v\n", readyIdx)))
			buff.Write([]byte(fmt.Sprintf("NODE_HOST=%s\n", nd.Host)))
			err := ioutil.WriteFile(path, buff.Bytes(), 0664)
			if err != nil {
				return err
			}
			if inst.Status == model.Created {
				inst.Status = model.Ready
				_ = nd.UpdateInstance(inst)
				_ = cl.UpdateNode(nd)
				c.internal.UpdateCluster(cl)
				if err := c.commit(); err != nil {
					return err
				}
			}
			c.print(model.SuccessTypeResponse{
				Type:    "Instance",
				SubType: "Preparation",
				Content: path,
			})
		}
	}
	return nil
}
