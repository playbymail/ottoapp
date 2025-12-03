// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cst2

//go:generate stringer --type Kind

import (
	"bytes"
	"fmt"

	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

type Node interface {
	Kind() Kind
	Pos() Position
	String() string
	Source() []byte // for reconstructing source
}

type Position struct {
	Line, Col int
}

type Kind int

const (
	Base Kind = iota
	TurnReport
	PrologueSection // lines before first unit section
	Line
	UnitSection
	UnitLine
	BadUnitSection
	EpilogueSection
)

type BaseNode struct {
	kind     Kind
	position Position
	errors   []error
	tokens   []*lexers.Token
}

func (n BaseNode) Kind() Kind     { return n.kind }
func (n BaseNode) Pos() Position  { return n.position }
func (n BaseNode) String() string { return fmt.Sprintf("<%s>", n.Kind()) }
func (n BaseNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, token := range n.tokens {
		b.Write(token.Source())
	}
	return b.Bytes()
}

// TurnReportNode is the root of the CST.
type TurnReportNode struct {
	BaseNode
	prologue *PrologueSectionNode
	sections []Node
	epilogue *EpilogueSectionNode
}

func newTurnReportNode() *TurnReportNode {
	return &TurnReportNode{BaseNode: BaseNode{kind: TurnReport}}
}

func (n *TurnReportNode) Source() []byte {
	b := &bytes.Buffer{}
	if n.prologue != nil {
		b.Write(n.prologue.Source())
		b.WriteString("8<------- prologue ------->8\n")
	}
	for _, section := range n.sections {
		b.WriteString("8<----- unit section ----->8\n")
		b.Write(section.Source())
	}
	if n.epilogue != nil {
		b.WriteString("8<------- epilogue ------->8\n")
		b.Write(n.epilogue.Source())
	}
	return b.Bytes()
}

// PrologueSectionNode collects the lines prior to the first unit section.
type PrologueSectionNode struct {
	BaseNode
	lines []*LineNode
}

func newPreambleSectionNode() *PrologueSectionNode {
	return &PrologueSectionNode{BaseNode: BaseNode{kind: PrologueSection}}
}

func (n *PrologueSectionNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, line := range n.lines {
		b.Write(line.Source())
	}
	return b.Bytes()
}

// LineNode captures unknown lines
type LineNode struct {
	BaseNode
	line []*lexers.Token
}

func newLineNode() *LineNode {
	return &LineNode{BaseNode: BaseNode{kind: Line}}
}

func (n *LineNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, token := range n.tokens {
		b.Write(token.Source())
	}
	return b.Bytes()
}

// UnitSectionNode is
type UnitSectionNode struct {
	BaseNode
	nodes    []Node
	epilogue []*LineNode
}

func newUnitSectionNode() *UnitSectionNode {
	return &UnitSectionNode{BaseNode: BaseNode{kind: UnitSection}}
}

func (n *UnitSectionNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, node := range n.nodes {
		b.Write(node.Source())
	}
	if n.epilogue != nil {
		b.WriteString("  8<----- epilogue ----->8\n")
		for _, line := range n.epilogue {
			b.Write(line.Source())
		}
	}
	return b.Bytes()
}

type UnitLineNode struct {
	BaseNode
	UnitKeyword    *lexers.Token
	UnitId         *lexers.Token
	Comma1         *lexers.Token
	Note           *lexers.Token
	Comma2         *lexers.Token
	Current        *lexers.Token
	Hex1           *lexers.Token
	Equals         *lexers.Token
	CurrentCoords  *lexers.Token
	Comma3         *lexers.Token
	LeftParen      *lexers.Token
	Previous       *lexers.Token
	Hex2           *lexers.Token
	Equals2        *lexers.Token
	PreviousCoords *lexers.Token
	RightParen     *lexers.Token
}

func newUnitLineNode() *UnitLineNode {
	return &UnitLineNode{BaseNode: BaseNode{kind: UnitLine}}
}

// BadUnitSectionNode is generated from the top of the report until we encounter our first unit section
type BadUnitSectionNode struct {
	BaseNode
	lines []*LineNode
}

func newBadUnitSectionNode() *BadUnitSectionNode {
	return &BadUnitSectionNode{BaseNode: BaseNode{kind: BadUnitSection}}
}

func (n *BadUnitSectionNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, line := range n.lines {
		b.Write(line.Source())
	}
	return b.Bytes()
}

type EpilogueSectionNode struct {
	BaseNode
	lines []*LineNode
}

func newEpilogueSectionNode() *EpilogueSectionNode {
	return &EpilogueSectionNode{BaseNode: BaseNode{kind: EpilogueSection}}
}

func (n *EpilogueSectionNode) Source() []byte {
	b := &bytes.Buffer{}
	for _, line := range n.lines {
		b.Write(line.Source())
	}
	return b.Bytes()
}
