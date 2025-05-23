package replication

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"repligo/internal/network"
)

type Replicator struct {
	Slaves []network.Network
}

func NewReplicator(slaves []network.Network) *Replicator {
	return &Replicator{
		Slaves: slaves,
	}

}
func (r *Replicator) Replicate(key, value string) error {
	params := url.Values{}
	params.Add("key", key)
	params.Add("value", value)

	// Combine base URL with encoded query string
	for i, slave := range r.Slaves {
		url := fmt.Sprintf("http://%s:%d/replicate", slave.IP, slave.Port)
		fullURL := fmt.Sprintf("%s?%s", url, params.Encode())

		resp, err := http.Get(fullURL)
		if err != nil {
			log.Printf("❌ Slave %d (%s) failed: %v", i, url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("⚠️ Slave %d (%s) responded with status: %d", i, url, resp.StatusCode)
		} else {
			log.Printf("✅ Successfully replicated to slave %d (%s)", i, url)
		}
	}
	return nil
}
