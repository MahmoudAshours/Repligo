package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"repligo/internal/config"
	"repligo/internal/network"
	"repligo/internal/node"
	"repligo/internal/replication"
	"repligo/internal/storage"
)

func main() {
	cfg, err := config.Load("config.json")

	if err != nil {
		fmt.Println("Couldn't load slaves.")
	}

	storage := storage.NewInMemoryStore()
	role := flag.String("role", "slave", "Role of the node: master or slave")
	port := flag.Int("port", 8080, "Port to run the HTTP server on")
	flag.Parse()

	nodeNetwork := network.Network{
		IP:   "localhost.com",
		Port: *port,
	}
	node := &node.Node{
		NodeType: *role,
		Network:  nodeNetwork,
		Store:    storage,
	}

	fmt.Printf("🚀 Starting %s node on port %d...\n", *role, *port)

	// Set up HTTP handlers based on role
	if *role == "master" {
		http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
			masterWriteHandler(node, w, r, cfg.Cluster.Slaves)
		})
	} else {
		http.HandleFunc("/replicate", func(w http.ResponseWriter, r *http.Request) {
			slaveReplicateHandler(node, w, r)
		})
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Mock master write handler
func masterWriteHandler(n *node.Node, w http.ResponseWriter, r *http.Request, slaves []network.Network) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	fmt.Fprintf(w, "📝 Master received key=%s, value=%s\n", key, value)
	n.Store.Set(key, value)
	replicator := replication.NewReplicator(slaves)
	replicator.Replicate(key, value)
}

// Mock slave replicate handler
func slaveReplicateHandler(n *node.Node, w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	log.Printf("📥 Slave stored key=%s, value=%s\n", key, value)
	n.Store.Set(key, value)
	log.Println(n.Store.GetAll())
}
