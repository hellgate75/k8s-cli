package main

import "flag"

var command string
var nodeName string

func init()  {
	flag.StringVar(&command, "command", "help", "Required command action")
}

func main()  {

}
