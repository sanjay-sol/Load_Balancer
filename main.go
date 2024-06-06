package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type Node struct {
	URL    *url.URL
	Active bool
	mutex  sync.RWMutex
}

type NodePool struct {
	nodes   []*Node
	current uint64
}

func (np *NodePool) AddNode(n *Node) {
	np.nodes = append(np.nodes, n)
}

func (np *NodePool) NextIdx() int {
	return int(atomic.AddUint64(&np.current, uint64(1)) % uint64(len(np.nodes)))
}
// IsActive checks if the node is active
func (n *Node) isActive() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.Active
}

func main() {
	var nodeList string
	var port int
	flag.StringVar(&nodeList, "nodeList", "", "List of available nodes comma-separated")
	flag.IntVar(&port,"port", 3030, "Port to serve load-balancer")
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
