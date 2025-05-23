package node

import (
	"repligo/internal/network"
	"repligo/internal/storage"
)

type Node struct {
	Network  network.Network
	Store    storage.Store
	NodeType string
}
