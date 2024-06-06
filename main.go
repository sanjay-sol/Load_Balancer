package main 


import (
  "flag"
  "fmt"
  "log"
  "net/url"
)

// Node - has the data about the backend server

type Node struct {
  URL *url.URL
  Active bool 
  weight float64
}

// NodePool - Has slice of the nodes and most recently used node indexx

type NodePool struct {
  nodes []*Node 
  current uint64
}

func main() {
	var nodeList string
	var port int
	flag.StringVar(&nodeList, "nodeList", "", "List of available nodes comma-separated")
	flag.IntVar(&port, "port", 3030, "Port to serve load-balancer")
	flag.Parse()

	fmt.Println("Node List:", nodeList)
	fmt.Println("Port:", port)
	log.Println("Load Balancer started")
}


