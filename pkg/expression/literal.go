package expression

import "fmt"

// NUMBER | STRING | TRUE | FALSE | NIL
type Literal struct {
	Value any
}

func (l Literal) String() string {
	switch v := l.Value.(type) {
	case string:
		return fmt.Sprintf("<\"%s\">", v)
	case float64:
		return fmt.Sprintf("<%f>", v)
	case bool:
		if v {
			return "<TRUE>"
		}
		return "<FALSE>"
	case nil:
		return "<NIL>"
	default:
		return fmt.Sprintf("<%v>", v)
	}
}
