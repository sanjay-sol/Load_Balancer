package balancer

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sanjay-sol/Load_Balancer/pkg/context"
)

//*  balancing incoming requests
func LoadBalancer(np *NodePool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		attempts := contextutil.GetAttemptsFromContext(r)
		if attempts > 3 {
			log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
			http.Error(w, "Service not available", http.StatusServiceUnavailable)
			return
		}
		node := np.NextNode()
		if node != nil {
			node.ReverseProxy.ServeHTTP(w, r)
			np.Heapify(0, true)
			return
		}
		//* 0 active available nodes
		http.Error(w, "Downtime: No nodes available", http.StatusServiceUnavailable)
	}
}

//* Adds a new node to the pool
func AddNode(np *NodePool, nodeURL string) {
	nodeURLParsed, err := url.Parse(nodeURL)
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(nodeURLParsed)
	proxy.ErrorHandler = func(w http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", nodeURLParsed.Host, e.Error())
		type RetryKey int

		const Retry RetryKey = 0

		retries := contextutil.GetRetryFromContext(request)
		if retries < 3 {
			<-time.After(10 * time.Millisecond)
			ctx := context.WithValue(request.Context(), Retry, retries+1)
			proxy.ServeHTTP(w, request.WithContext(ctx))
			return
		}

		//* Try different node
		attempts := contextutil.GetAttemptsFromContext(request)
		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
		ctx := context.WithValue(request.Context(), Retry, attempts+1)

		//* After 3 retries, set the current node as dead
		if attempts >= 3 {
			np.SetNodeStatus(nodeURLParsed, false)
		}

		LoadBalancer(np)(w, request.WithContext(ctx))
	}

	np.AddNode(&Node{
		URL:          nodeURLParsed,
		Active:       true,
		weight:       1,
		ReverseProxy: proxy,
	})

	log.Printf("Configured node: %s\n", nodeURLParsed)
}
