package balancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Node holds the data about a backend server
type Node struct {
	URL          *url.URL
	Active       bool
	weight       float64
	mutex        sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// isActive returns whether node is active or dead
func (n *Node) isActive() bool {
	var active bool
	n.mutex.RLock()
	active = n.Active
	n.mutex.RUnlock()
	return active
}

// getWeight returns the weight of the node
func (n *Node) getWeight() float64 {
	n.mutex.RLock()
	weight := n.weight
	n.mutex.RUnlock()
	return weight
}

// SetProps sets node's status and changes node's weight
func (n *Node) SetProps(status bool) {
	n.mutex.Lock()
	n.Active = status
	if !status {
		n.weight /= 3.0
	} else if n.weight < 1 {
		n.weight *= 2.0
	}
	n.mutex.Unlock()
}
