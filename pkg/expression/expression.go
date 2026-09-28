package expression

import (
	"fmt"
	"sync/atomic"
)

var nextNodeID uint64

type Expression interface {
	fmt.Stringer
}

func NextNodeID() uint64 {
	return atomic.AddUint64(&nextNodeID, 1)
}
