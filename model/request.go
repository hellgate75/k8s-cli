package model

type CommandRequest struct{
	Command			string
	SubCommand		string
	ClusterName 	string
	KubeCtlFile		string
	NodeName		string
	HostName		string
	NodeSlots		int
	Instance		string
	Namespace		string
	Format			string
}

