package health

import (
	"log"
	"time"

	"github.com/sanjay-sol/Load_Balancer/internal/balancer"
)

//* StartHealthCheck starts periodic health checks ( for every 10 seconds )
func StartHealthCheck(np *balancer.NodePool) {
	t := time.NewTicker(time.Second * 10)
	for range t.C {
		log.Printf("### Starting health check ###")
		np.HealthCheck()
	}
}
