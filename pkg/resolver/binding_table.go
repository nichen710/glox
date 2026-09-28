package resolver

import (
	"fmt"
	"sort"
	"strings"
)

type BindingTable struct {
	depths map[uint64]int
	names  map[uint64]string
}

func NewBindingTable() *BindingTable {
	return &BindingTable{
		depths: make(map[uint64]int),
		names:  make(map[uint64]string),
	}
}

func (t *BindingTable) Resolve(nodeID uint64, depth int) {
	t.depths[nodeID] = depth
}

func (t *BindingTable) ResolveNamed(nodeID uint64, name string, depth int) {
	t.depths[nodeID] = depth
	t.names[nodeID] = name
}

func (t *BindingTable) Depth(nodeID uint64) (int, bool) {
	depth, ok := t.depths[nodeID]
	return depth, ok
}

func (t *BindingTable) String() string {
	if len(t.depths) == 0 {
		return "{}"
	}
	ids := make([]uint64, 0, len(t.depths))
	for id := range t.depths {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	var parts []string
	for _, id := range ids {
		name := t.names[id]
		if name != "" {
			parts = append(parts, fmt.Sprintf("%s (id %d): depth %d", name, id, t.depths[id]))
		} else {
			parts = append(parts, fmt.Sprintf("id %d: depth %d", id, t.depths[id]))
		}
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

