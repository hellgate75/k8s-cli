package main

import (
	"flag"
	"fmt"
	"github.com/hellgate75/k8s-cli/common"
	"github.com/hellgate75/k8s-cli/executor"
	"github.com/hellgate75/k8s-cli/io"
	"github.com/hellgate75/k8s-cli/model"
	"os"
	"strings"
	"time"
)

var format string
var command string
var subcommand string
var subsubcommand string
var clusterName string
var clusterConfigYaml string
var nodeName string
var hostName string
var nodeSlots int
var verifySlots bool
var prepareationName string
var namespace string

var dataDir string

func initHelp() {
	flag.StringVar(&command, "command", "help", "Required executor action (show, add, remove, verify, prepare. ensure, help)")
	flag.StringVar(&subcommand, "subject", "", "Required executor action subject (cluster, node, instance) or executor in case of help")
	flag.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	flag.StringVar(&subsubcommand, "details", "", "Required executor action subject (cluster, node, instance) only in case of help")
	flag.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	flag.BoolVar(&verifySlots, "verify-slots", false, "Retrun information about Free slots for nodes and clusters")
	flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	flag.StringVar(&clusterConfigYaml, "kubectl-yaml-file", "", "Kubectl Yaml file")
	flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
	flag.StringVar(&hostName, "node-host-name", "", "Cluster node host name")
	flag.IntVar(&nodeSlots, "node-slots", 2, "Cluster node max number of prepareations")
	flag.StringVar(&prepareationName, "instance-name", "", "Cluster node instance name")
	flag.StringVar(&namespace, "namespace", "", "Cluster node instance namespace name")
}

func showCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Show clusters, nodes or instances details:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: show %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "show", "Required executor action : add")
	fset.StringVar(&subcommand, "subject", "cluster", "Required executor action subject (clusters, nodes, instances)")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	fset.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	fset.BoolVar(&verifySlots, "verify-slots", false, "Retrun information about Free slots for nodes and clusters")
	if subCommand == "nodes" || subCommand == "instances" {
		fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
		if subCommand == "instances" {
			fset.StringVar(&nodeName, "node-name", "", "Cluster node name")
		}
	}
	return fset
}
func ensureCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Calculate first node available for a new instance:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: ensure %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "ensure", "Required executor action : ensure")
	fset.StringVar(&subcommand, "subject", "cluster", "Optional executor action subject (instance)")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	fset.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	return fset
}

func addCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Add a new cluster, node or instance:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: add %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "add", "Required executor action : show")
	fset.StringVar(&subcommand, "subject", "cluster", "Required executor action subject (cluster, node, instance)")
	fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	fset.StringVar(&clusterConfigYaml, "kubectl-yaml-file", "", "Kubectl Yaml file")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	fset.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	if subCommand == "node" || subCommand == "instance" {
		fset.StringVar(&nodeName, "node-name", "", "Cluster node name")
		fset.StringVar(&hostName, "node-host-name", "", "Cluster node host name")
		fset.IntVar(&nodeSlots, "node-slots", 2, "Cluster node max number of prepareations")
		if subCommand == "instance" {
			fset.StringVar(&prepareationName, "instance-name", "", "Cluster node instance name")
			fset.StringVar(&namespace, "namespace", "", "Cluster node instance namespace name")
		}
	}
	return fset
}

func removeCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Remove an existing cluster, node or instance:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: remove %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "remove", "Required executor action : remove")
	fset.StringVar(&subcommand, "subject", "cluster", "Required executor action subject (cluster, node, instance)")
	fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	flag.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	if subCommand == "node" || subCommand == "instance" {
		fset.StringVar(&nodeName, "node-name", "", "Cluster node name")
		if subCommand == "instance" {
			fset.StringVar(&prepareationName, "instance-name", "", "Cluster node instance name")
		}
	}
	return fset
}

func verifyCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Verify healthy state of an existing cluster, node or instance:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: verify %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "verify", "Required executor action : verify")
	fset.StringVar(&subcommand, "subject", "cluster", "Required executor action subject (cluster, node, instance)")
	fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	flag.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	if subCommand == "node" || subCommand == "instance" {
		fset.StringVar(&nodeName, "node-name", "", "Cluster node name")
		if subCommand == "instance" {
			fset.StringVar(&prepareationName, "instance-name", "", "Cluster node instance name")
		}
	}
	return fset
}

func prepareCommandInit(subCommand string) *flag.FlagSet {
	fmt.Println("Prepare a deployment environment file:")
	fset := flag.NewFlagSet(fmt.Sprintf("k8s-cli (cmd: prepare %s)", subCommand), flag.ContinueOnError)
	fset.StringVar(&command, "command", "prepare", "Required executor action : prepare")
	fset.StringVar(&subcommand, "subject", "cluster", "Required executor action subject (instance)")
	fset.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	fset.StringVar(&dataDir, "config-dir", common.ConfigDir(), "Configuration folder")
	fset.StringVar(&format, "format", "json", "Required output format (json, yaml), in case of error or missing will be used JSON")
	if subCommand == "instance" {
		fset.StringVar(&nodeName, "node-name", "", "Cluster node name")
		fset.StringVar(&prepareationName, "instance-name", "", "Cluster node instance name")
	}
	return fset
}

func waitApp(folder string) {
	if _, err := os.Stat(folder); err != nil {
		_ = os.MkdirAll(folder, 0660)
	}
	var lockFile = fmt.Sprintf("%s%c%s", folder, os.PathSeparator, ".lock")
	_, err := os.Stat(lockFile)
	for err == nil {
		time.Sleep(5 * time.Second)
		_, err = os.Stat(lockFile)
	}
}

func unlockApp(folder string) bool {
	if _, err := os.Stat(folder); err != nil {
		_ = os.MkdirAll(folder, 0660)
	}
	var lockFile = fmt.Sprintf("%s%c%s", folder, os.PathSeparator, ".lock")
	if _, err := os.Stat(lockFile); err == nil {
		err := os.Remove(lockFile)
		return err == nil
	}
	return false
}

func lockApp(folder string) bool {
	if _, err := os.Stat(folder); err != nil {
		_ = os.MkdirAll(folder, 0660)
	}
	var lockFile = fmt.Sprintf("%s%c%s", folder, os.PathSeparator, ".lock")
	if _, err := os.Stat(lockFile); err != nil {
		_, err := os.Create(lockFile)
		return err == nil
	}
	return false
}

func main() {
	initHelp()
	var args = []string(os.Args)
	flag.Parse()
	switch strings.ToLower(command) {
	case "help":
		fmt.Println(command)
		fmt.Println(subcommand)
		fmt.Println(subsubcommand)
		switch strings.ToLower(subcommand) {
		case "show":
			fset := showCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		case "add":
			fset := addCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		case "ensure":
			fset := ensureCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		case "remove":
			fset := removeCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		case "verify":
			fset := verifyCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		case "prepare":
			fset := prepareCommandInit(subsubcommand)
			fset.Parse(args)
			fset.Usage()
		default:
			fmt.Printf("Requested help of unknown subject: %s\n", subcommand)
			flag.Usage()
		}
	case "show", "add", "remove", "verify", "ensure", "prepare":
		waitApp(dataDir)
		_ = lockApp(dataDir)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Error:", r)
			}
			_ = unlockApp(dataDir)
		}()
		exec := executor.New(dataDir, model.CommandRequest{
			Command:     command,
			SubCommand:  subcommand,
			ClusterName: clusterName,
			KubeCtlFile: clusterConfigYaml,
			NodeName:    nodeName,
			HostName:    hostName,
			NodeSlots:   nodeSlots,
			Instance:    prepareationName,
			Namespace:   namespace,
			Format:      common.FixOutputType(format),
			VerifySlots: verifySlots,
		})
		err := exec.Init()
		if err != nil {
			errResp := model.ErrorResponse{
				Code:    401,
				Message: fmt.Sprintf("%v", err),
			}
			var dt []byte
			if common.FixOutputType(format) == "yaml" {
				dt, _ = io.ToYaml(errResp)
			} else {
				dt, _ = io.ToJson(errResp)
			}
			fmt.Sprintf("%s", string(dt))
		}
		err = exec.Execute()
		if err != nil {
			errResp := model.ErrorResponse{
				Code:    402,
				Message: fmt.Sprintf("%v", err),
			}
			var dt []byte
			if common.FixOutputType(format) == "yaml" {
				dt, _ = io.ToYaml(errResp)
			} else {
				dt, _ = io.ToJson(errResp)
			}
			fmt.Printf("%s\n", string(dt))
		}
	default:
		fmt.Printf("Unknown executor : %s\n", command)
		flag.Usage()
	}
}
