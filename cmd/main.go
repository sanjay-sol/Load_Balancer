package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/sanjay-sol/Load_Balancer/internal/balancer"
	"github.com/sanjay-sol/Load_Balancer/internal/health"
)

var nodePool balancer.NodePool

func main() {
	var nodeList string
	var port int
	flag.StringVar(&nodeList, "nodeList", "", "List the nodes comma-separated")
	flag.IntVar(&port, "port", 3030, "Port to serve load-balancer")
	flag.Parse()

	if len(nodeList) == 0 {
		log.Fatal("Please provide one or more nodes to start load balancing...")
	}

	for _, nodeURL := range strings.Split(nodeList, ",") {
		balancer.AddNode(&nodePool, nodeURL)
	}

	//* Create server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(balancer.LoadBalancer(&nodePool)),
	}

	go health.StartHealthCheck(&nodePool)

	log.Printf("Load Balancer started on port: %d", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
