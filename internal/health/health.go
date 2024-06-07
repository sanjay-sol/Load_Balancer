package health

import (
	"log"
	"time"

	"github.com/sanjay-sol/Load_Balancer/internal/balancer"
)

// StartHealthCheck starts periodic health checks
func StartHealthCheck(np *balancer.NodePool) {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			log.Printf("Starting health check...")
			np.HealthCheck()
		}
	}
}
