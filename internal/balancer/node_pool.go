package balancer

import (
	"net/url"
	"log"
	"sync/atomic"
	"net"
	"time"
)

//* NodePool holds most recently used node index and slice of nodes 
type NodePool struct {
	nodes   []*Node
	current uint64
}

//* AddNode to the Pool
func (np *NodePool) AddNode(n *Node) {
	np.nodes = append(np.nodes, n)
}

//* NextIdx atomically increases the counter and returns an index
func (np *NodePool) NextIdx() int {
	return int(atomic.AddUint64(&np.current, uint64(1)) % uint64(len(np.nodes)))
}

func (np *NodePool) Swap(i uint64, j uint64) {
	temp := np.nodes[i]
	np.nodes[i] = np.nodes[j]
	np.nodes[j] = temp
}

//*  Will rearrange the max heap based on weights, takes index and if the node is root
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

//!  finds the next active node - currently based on max Heap , TODO - RoundRobin 
func (np *NodePool) NextNode() *Node {
	return np.nodes[0]
}

//* Sets status of the given nodeURL
func (np *NodePool) SetNodeStatus(url *url.URL, status bool) {
	for _, n := range np.nodes {
		if n.URL.String() == url.String() {
			n.SetProps(status)
			break
		}
	}
}

func (n *Node) Status() bool {
	conn, err := net.DialTimeout("tcp", n.URL.Host, 2*time.Second)
	if err != nil {
		log.Println("Node unreachable: ", err)
		return false
	}
	_ = conn.Close()
	return true
}

//* Pings the node and updates status
func (np *NodePool) HealthCheck() {
	for i, n := range np.nodes {
		status := n.Status()
		n.SetProps(status)
		msg := "active"
		if !status {
			msg = "dead"
			np.Heapify(uint64(i), false)
		}
		log.Printf("%s [%s] [%0.2g]\n", n.URL, msg, n.weight)
	}
}

