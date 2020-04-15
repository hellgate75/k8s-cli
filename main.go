package main

import (
	"flag"
	"fmt"
	"strings"
)

var command string
var subcommand string
var clusterName string
var clusterConfigYaml string
var nodeName string
var nodeSlots int
var installationName string
var namespace string

func init()  {
	flag.StringVar(&command, "command", "help", "Required command action (show, add, remove, check, install, help)")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (cluster, node, instance)")
}


func showCommandInit(subCommand string) {
	flag.StringVar(&command, "command", "help", "Required command action : show")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (cluster, node, instance)")
	flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	flag.StringVar(&clusterConfigYaml, "kubectl-yaml-file", "", "Kubectl Yaml file")
	if subCommand == "node" || subCommand == "installation" {
		flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
		flag.IntVar(&nodeSlots, "node-slots", 2, "Cluster node max number of installations")
		if subCommand == "instance" {
			flag.StringVar(&installationName, "installation-name", "", "Cluster node installation name")
			flag.StringVar(&namespace, "namespace", "", "Cluster node installation namespace name")
		}
	}
}


func addCommandInit(subCommand string) {
	flag.StringVar(&command, "command", "help", "Required command action : add")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (cluster, node, instance)")
	flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	flag.StringVar(&clusterConfigYaml, "kubectl-yaml-file", "", "Kubectl Yaml file")
	if subCommand == "node" || subCommand == "installation" {
		flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
		flag.IntVar(&nodeSlots, "node-slots", 2, "Cluster node max number of installations")
		if subCommand == "instance" {
			flag.StringVar(&installationName, "installation-name", "", "Cluster node installation name")
			flag.StringVar(&namespace, "namespace", "", "Cluster node installation namespace name")
		}
	}
}

func removeCommandInit(subCommand string) {
	flag.StringVar(&command, "command", "help", "Required command action : remove")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (cluster, node, instance)")
	flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	if subCommand == "node" || subCommand == "installation" {
		flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
		if subCommand == "instance" {
			flag.StringVar(&installationName, "installation-name", "", "Cluster node installation name")
		}
	}
}

func checkCommandInit(subCommand string) {
	flag.StringVar(&command, "command", "help", "Required command action : check")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (cluster, node, instance)")
	flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
	if subCommand == "node" || subCommand == "installation" {
		flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
		if subCommand == "instance" {
			flag.StringVar(&installationName, "installation-name", "", "Cluster node installation name")
		}
	}
}

func installCommandInit(subCommand string) {
	flag.StringVar(&command, "command", "help", "Required command action : add")
	flag.StringVar(&subcommand, "subject", "cluster", "Required command action subject (clusters, nodes, instances)")
	if subCommand == "nodes" || subCommand == "installation" {
		flag.StringVar(&clusterName, "cluster-name", "default", "Cluster name")
		if subCommand == "instances" {
			flag.StringVar(&nodeName, "node-name", "", "Cluster node name")
		}
	}
}

func main()  {
	flag.Parse()
	switch strings.ToLower(command) {
	case "help":
		if "" == subcommand {
			flag.Usage()
		} else {
			switch strings.ToLower(subcommand) {
			case "show":
				showCommandInit(subcommand)
				flag.Parse()
				flag.Usage()
			case "add":
				addCommandInit(subcommand)
				flag.Parse()
				flag.Usage()
			case "remove":
				removeCommandInit(subcommand)
				flag.Parse()
				flag.Usage()
			case "check":
				checkCommandInit(subcommand)
				flag.Parse()
				flag.Usage()
			case "install":
				installCommandInit(subcommand)
				flag.Parse()
				flag.Usage()
			default:
				flag.Usage()
			}
		}
	default:
		fmt.Printf("Unknown command : %s\n", command)
		flag.Usage()
	}
}
