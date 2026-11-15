package expression

import "glox/pkg/token"

// "this"
type This struct {
	ID      uint64
	Keyword token.Token
}

func NewThis(keyword token.Token) This {
	return This{
		ID:      NextNodeID(),
		Keyword: keyword,
	}
}

func (t This) String() string {
	return "this"
}
