// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package make implements a dependency management service.
package make

import (
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

type Node struct {
	ID         domains.ID
	ModifiedAt time.Time

	// Outgoing edges: this node depends on these nodes.
	DependsOn []domains.ID
}

type Graph struct {
	Nodes map[domains.ID]*Node
}
