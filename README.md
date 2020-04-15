<p align="right">
 <img src="https://github.com/hellgate75/k8s-cli/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/k8s-cli.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-tcp-server">Check last build on Travis-CI</a>
 </p>
<p align="center">
<image width="260" height="410" src="images/k8s-logo.png">
</p>
<br/>
<br/>

# Kubernetes Clusters Client
Go-Language CI-CD Kubernetes Cluster Client

## Purposes of the project

Provide easy to use infrastructure tools that collect information about K8S clusters as json database and order/execute helm charts installation.


## Command Options

Command provides some fetaures, as follow:

* `-command show` (cluster, nodes, history, instances) - Provides the list of available cluster, clusters, nodes, installations or the history of commands

* `-command add` (cluster, node, instance) Add  a node or an instance to a node

* `-command remove` (cluster, node, instance) Remove  a node or an instance from a node

* `-command check` (cluster, node, instance) Verifies with kube-ctl commands the availability of the cluster, node cluster, etc...

* `-command install` Install charts from a list to the node

* `-command help` Shows commands or command details (if used help <command>)


## Build the project

Build command sample :
```
go build -buildmode=exe github.com/hellgate75/k8s-cli
```

## Get the executable

Build command sample :
```
go get -u github.com/hellgate75/k8s-cli
```

Enjoy the experience.

## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)

