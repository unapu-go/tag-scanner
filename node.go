package tag_scanner

import "fmt"

const (
	KeyValue NodeType = iota + 1
	Tags
	Flag
)

type NodeType uint8

type Node interface {
	fmt.Stringer
	Type() NodeType
}

type NodeKeyValue struct {
	Key, Value string
	KeyArgs []string
}

func (NodeKeyValue) Type() NodeType {
	return KeyValue
}

func (this NodeKeyValue) String() string {
	return this.Key + ":" + this.Value
}

type NodeTags string

func (NodeTags) Type() NodeType {
	return Tags
}

func (this NodeTags) String() string {
	return string(this)
}

func (this NodeTags) Value() string {
	return string(this)
}

type NodeFlag string

func (NodeFlag) Type() NodeType {
	return Flag
}

func (this NodeFlag) String() string {
	return string(this)
}