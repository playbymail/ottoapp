// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package iana

import (
	"sort"
	"strings"
	"time"
)

type Zones []*Zone

type Zone struct {
	Name     string
	Location *time.Location
	SubZones []*Zone
}

// loadZones builds a tree of Zones from a list of IANA timezone names.
//
// It returns the structured Zones and the first error encountered while
// loading time zone data (e.g. if the local tzdata is out of sync).
func loadZones(names []string) (Zones, error) {
	root := make(map[string]*Zone)
	var firstErr error

	for _, full := range names {
		if full == "" {
			continue
		}
		parts := strings.Split(full, "/")

		parentMap := root
		var parent *Zone

		for i, part := range parts {
			var node *Zone

			if parent == nil {
				// top-level node
				node = parentMap[part]
				if node == nil {
					node = &Zone{Name: part}
					parentMap[part] = node
				}
			} else {
				// nested node
				node = findChild(parent.SubZones, part)
				if node == nil {
					node = &Zone{Name: part}
					parent.SubZones = append(parent.SubZones, node)
				}
			}

			if i == len(parts)-1 {
				loc, err := time.LoadLocation(full)
				if err != nil && firstErr == nil {
					firstErr = err
				} else {
					node.Location = loc
				}
			}

			parent = node
		}
	}

	// Convert map to sorted slice
	zs := make(Zones, 0, len(root))
	for _, z := range root {
		sortZoneTree(z)
		zs = append(zs, z)
	}
	sort.Slice(zs, func(i, j int) bool { return zs[i].Name < zs[j].Name })

	// Ensure parent nodes don't have a Location
	for _, z := range zs {
		ensureParentLocationNil(z)
	}

	return zs, firstErr
}

func findChild(children []*Zone, name string) *Zone {
	for _, c := range children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func sortZoneTree(z *Zone) {
	if len(z.SubZones) == 0 {
		return
	}
	sort.Slice(z.SubZones, func(i, j int) bool {
		return z.SubZones[i].Name < z.SubZones[j].Name
	})
	for _, c := range z.SubZones {
		sortZoneTree(c)
	}
}

func ensureParentLocationNil(z *Zone) {
	if len(z.SubZones) > 0 {
		z.Location = nil
		for _, c := range z.SubZones {
			ensureParentLocationNil(c)
		}
	}
}
