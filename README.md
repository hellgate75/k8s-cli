<p align="right">
 <img src="https://github.com/hellgate75/k8s-cli/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/k8s-cli.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/k8s-cli">Check last build on Travis-CI</a>
 </p>
<p align="center">
<image width="450" height="260" src="images/k8s-logo.png">
</p>
<br/>
<br/>

# Kubernetes Clusters Client
Go-Language CI-CD Kubernetes Cluster Client

## Purposes of the project

Provide easy to use infrastructure tools that collect information about K8S clusters as json database and order/execute helm charts installation.


## Command Options

Command provides some fetaures, as follow:

* `-command show` *(clusters, nodes, instances)* - Provides the list of available cluster, clusters, nodes, installations

* `-command details` *(cluster, node, instance)* - Provides the list of specific cluster, cluster, node, installation

* `-command discover` *(nodes)* - Discover New Kubernetes *Added* and *Ready* nodes and creates eventually nodes in the local configuration

* `-command add` *(cluster, node, instance)* Add  a node or an instance to a node

* `-command remove` *(cluster, node, instance)* Remove  a node or an instance from a node

* `-command check` *(cluster, node, instance)* Verifies with kube-ctl commands the availability of the cluster, node cluster, etc...

* `-command prepare` *(instance)*  Prepare a spacific  instance environment for installation purposes

* `-command ensure` *(instance)* Verify first node with availability to deploy an helm instance, not taken yet

* `-command help` *(show, details, discover, add, remove, verify, prepare. ensure)* Shows commands or command details (if used help <command>)


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

