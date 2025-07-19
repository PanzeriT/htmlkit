package htmlkit

import (
	"bytes"
	"fmt"
	"html"
)

type Element interface {
	Render() []byte
}

type Tag struct {
	Name             string
	Attributes       map[string]string
	Children         []Element
	AllowedAttrs     map[string]struct{}
	AllowedChildTags map[string]struct{}
}

func (t *Tag) Render() []byte {
	var b bytes.Buffer

	// < ... > first part of the tag
	b.WriteString("<" + t.Name)
	for k, v := range t.Attributes {
		b.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
	}
	b.WriteString(">")

	// > ... < children of the tag
	for _, child := range t.Children {
		b.Write(child.Render())
	}

	// </ ... > end of the tag
	b.WriteString("</" + t.Name + ">")

	return b.Bytes()
}

func (t *Tag) AddChild(child Element) {
	if _, ok := t.AllowedChildTags[child.(*Tag).Name]; !ok {
		fmt.Printf("Warning: Child tag '%s' is not allowed for tag '%s'\nStop rendering.\n", child.(*Tag).Name, t.Name)
		return
	}
	t.Children = append(t.Children, child)
}

func (t *Tag) AddAttribute(attr Attr) {
	if _, ok := t.AllowedAttrs[attr.Key]; !ok {
		fmt.Printf("Warning: Attribute '%s' is not allowed for tag '%s'\nStop rendering.\n", attr.Key, t.Name)
	}
	t.Attributes[attr.Key] = attr.Val
}

type TextNode struct {
	Content string
}

func (n *TextNode) Render() []byte {
	return []byte(html.EscapeString(n.Content))
}

func Text(s string) Element {
	return &TextNode{
		Content: s,
	}
}

type RawNode struct {
	Content string
}

func (n *RawNode) Render() []byte {
	return []byte(n.Content)
}

func Raw(s string) Element {
	return &RawNode{
		Content: s,
	}
}

type Attr struct {
	Key string
	Val string
}

func Class(name string) Attr {
	return Attr{Key: "class", Val: name}
}

func Id(id string) Attr {
	return Attr{Key: "id", Val: id}
}

func tag(name string, allowedAttrs, allowedChildTags []string, args ...any) *Tag {
	t := &Tag{
		Name:             name,
		Attributes:       make(map[string]string),
		AllowedAttrs:     make(map[string]struct{}),
		AllowedChildTags: make(map[string]struct{}),
	}
	for _, attr := range allowedAttrs {
		t.AllowedAttrs[attr] = struct{}{}
	}
	for _, child := range allowedChildTags {
		t.AllowedChildTags[child] = struct{}{}
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case Attr:
			t.AddAttribute(v)
		case *Tag:
			t.AddChild(v)
		case Element:
			t.Children = append(t.Children, v)
		}
	}
	return t
}

func Table(args ...any) *Tag { return tag("table", []string{"class", "id"}, []string{"tr"}, args...) }
func Tr(args ...any) *Tag    { return tag("tr", []string{"class", "id"}, []string{"td"}, args...) }
func Td(args ...any) *Tag    { return tag("td", []string{"class", "id"}, nil, args...) }
