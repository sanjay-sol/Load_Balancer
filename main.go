package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type Node struct {
	URL    *url.URL
	Active bool
	weight float64
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

func (np *NodePool) Swap(i uint64, j uint64) {
	temp := np.nodes[i]
	np.nodes[i] = np.nodes[j]
	np.nodes[j] = temp
}

func (np *NodePool) Heapify(idx uint64, root bool) {
	largest := idx
	left := 2*idx + 1
	right := 2*idx + 2

	if root {
		np.nodes[idx].weight /= 2
	}

	if left < uint64(len(np.nodes)) && np.nodes[left].isActive() && np.nodes[left].getWeight() > np.nodes[largest].getWeight() {
		largest = left
	}

	if right < uint64(len(np.nodes)) && np.nodes[right].isActive() && np.nodes[right].getWeight() > np.nodes[largest].getWeight() {
		largest = right
	}

	if largest != idx {
		if root {
			np.nodes[idx].weight *= 2
		}
		np.Swap(largest, idx)
		np.Heapify(largest, false)
	}

	if left < uint64(len(np.nodes)) && np.nodes[left].getWeight() < 1 {
		np.Heapify(left, false)
	}

	if right < uint64(len(np.nodes)) && np.nodes[right].getWeight() < 1 {
		np.Heapify(right, false)
	}
}

func (np *NodePool) NextNode() *Node {
	return np.nodes[0]
}

func (n *Node) isActive() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.Active
}

func (n *Node) getWeight() float64 {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.weight
}

func loadBalancer(w http.ResponseWriter, r *http.Request) {
	node := nodePool.NextNode()
	if node != nil {
		// Dummy reverse proxy logic (to be replaced with actual proxy logic)
		w.Write([]byte(fmt.Sprintf("Forwarding to node: %s", node.URL)))
		nodePool.Heapify(0, true)
		return
	}
	http.Error(w, "No nodes available", http.StatusServiceUnavailable)
}

var nodePool NodePool

func main() {
	var nodeList string
	var port int
	flag.StringVar(&nodeList, "nodeList", "", "List of available nodes comma-separated")
	flag.IntVar(&port,"port", 3030, "Port to serve load-balancer")
	flag.Parse()

	for _, nodeURL := range strings.Split(nodeList, ",") {
		nodeURLParsed, err := url.Parse(nodeURL)
		if err != nil {
			log.Fatal(err)
		}
		nodePool.AddNode(&Node{
			URL:    nodeURLParsed,
			Active: true,
			weight: 1,
		})
	}

	http.HandleFunc("/", loadBalancer)
	fmt.Println("Node List:", nodeList)
	fmt.Println("Port:", port)
	log.Println("Load Balancer started")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
