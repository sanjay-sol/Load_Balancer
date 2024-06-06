package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
  "strings"
)

type Node struct {
	URL    *url.URL
	Active bool
}

type NodePool struct {
	nodes []*Node
}

func (np *NodePool) AddNode(n *Node) {
	np.nodes = append(np.nodes, n)
}

func main() {
	var nodeList string
	var port int
	flag.StringVar(&nodeList, "nodeList", "", "List of available nodes comma-separated")
  flag.IntVar(&port, "port", 8080, "Port to run the load balancer")
	flag.Parse()

	nodePool := &NodePool{}
	for _, nodeURL := range strings.Split(nodeList, ",") {
		nodeURLParsed, err := url.Parse(nodeURL)
		if err != nil {
			log.Fatal(err)
		}
		nodePool.AddNode(&Node{
			URL:    nodeURLParsed,
			Active: true,
		})
	}

	fmt.Println("Node List:", nodeList)
	fmt.Println("Port:", port)
	log.Println("Load Balancer started")
}
