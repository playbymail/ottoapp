// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package status

//go:generate pigeon -o grammar.go grammar.peg

type Status struct {
	UnitId  string
	Terrain string
	Rest    string
}
