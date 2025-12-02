// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cst

import "github.com/playbymail/ottoapp/backend/services/reports/lexers"

// Node is the interface implemented by all CST nodes.
type Node interface {
	Errors() []error
	Tokens() []*lexers.Token // tokens consumed by this node (for source reconstruction)
}

// TurnReportNode is the root of the CST.
type TurnReportNode struct {
	Sections []*UnitSectionNode
	EOF      *lexers.Token
	errors   []error
	tokens   []*lexers.Token
}

func (n *TurnReportNode) Errors() []error         { return n.errors }
func (n *TurnReportNode) Tokens() []*lexers.Token { return n.tokens }

// UnitSectionNode represents a unit section (currently just a unit line).
type UnitSectionNode struct {
	UnitLine *UnitLineNode
	errors   []error
	tokens   []*lexers.Token
}

func (n *UnitSectionNode) Errors() []error         { return n.errors }
func (n *UnitSectionNode) Tokens() []*lexers.Token { return n.tokens }

// UnitLineNode represents a unit line.
// Example: Tribe 0987, , Current Hex = QQ 1315, (Previous Hex = N/A)
type UnitLineNode struct {
	Keyword     *lexers.Token // Tribe | Courier | Element | Fleet | Garrison
	UnitID      *lexers.Token // Number or UnitId
	Comma1      *lexers.Token
	Note        *lexers.Token // optional
	Comma2      *lexers.Token
	Current     *lexers.Token
	Hex1        *lexers.Token
	Equals1     *lexers.Token
	CurrentHex  CoordsNode
	Comma3      *lexers.Token
	LeftParen   *lexers.Token
	Previous    *lexers.Token
	Hex2        *lexers.Token
	Equals2     *lexers.Token
	PreviousHex CoordsNode
	RightParen  *lexers.Token
	EOL         *lexers.Token
	errors      []error
	tokens      []*lexers.Token
}

func (n *UnitLineNode) Errors() []error         { return n.errors }
func (n *UnitLineNode) Tokens() []*lexers.Token { return n.tokens }

// CoordsNode is the interface for coordinate nodes.
type CoordsNode interface {
	Node
	coordsNode() // marker method
}

// GridCoordsNode represents grid coordinates (e.g., "QQ 1315").
type GridCoordsNode struct {
	Grid   *lexers.Token // e.g., "QQ"
	Number *lexers.Token // e.g., "1315"
	errors []error
	tokens []*lexers.Token
}

func (n *GridCoordsNode) Errors() []error         { return n.errors }
func (n *GridCoordsNode) Tokens() []*lexers.Token { return n.tokens }
func (n *GridCoordsNode) coordsNode()             {}

// NACoordsNode represents "N/A" coordinates (not available).
type NACoordsNode struct {
	Text   *lexers.Token // "N/A"
	errors []error
	tokens []*lexers.Token
}

func (n *NACoordsNode) Errors() []error         { return n.errors }
func (n *NACoordsNode) Tokens() []*lexers.Token { return n.tokens }
func (n *NACoordsNode) coordsNode()             {}

// ObscuredCoordsNode represents obscured coordinates (e.g., "## 1315").
type ObscuredCoordsNode struct {
	Grid   *lexers.Token // e.g., "##"
	Number *lexers.Token // e.g., "1315"
	errors []error
	tokens []*lexers.Token
}

func (n *ObscuredCoordsNode) Errors() []error         { return n.errors }
func (n *ObscuredCoordsNode) Tokens() []*lexers.Token { return n.tokens }
func (n *ObscuredCoordsNode) coordsNode()             {}

// ErrorNode represents a parsing error with recovery.
// It captures tokens that couldn't be parsed.
type ErrorNode struct {
	Message string
	errors  []error
	tokens  []*lexers.Token // tokens consumed during error recovery
}

func (n *ErrorNode) Errors() []error         { return n.errors }
func (n *ErrorNode) Tokens() []*lexers.Token { return n.tokens }

// ErrorCoordsNode is an error node that satisfies CoordsNode.
type ErrorCoordsNode struct {
	Message string
	errors  []error
	tokens  []*lexers.Token
}

func (n *ErrorCoordsNode) Errors() []error         { return n.errors }
func (n *ErrorCoordsNode) Tokens() []*lexers.Token { return n.tokens }
func (n *ErrorCoordsNode) coordsNode()             {}

// TurnLineNode represents a turn line.
// Example: Current Turn 899-12 (#0), Winter, FINE Next Turn 900-01 (#1), 28/11/2025
// Example: Current Turn 900-01 (#1), Spring, FINE
type TurnLineNode struct {
	Current1       *lexers.Token   // Current
	Turn1          *lexers.Token   // Turn
	TurnYearMonth1 *lexers.Token   // e.g., 899-12
	TurnNumber1    *TurnNumberNode // e.g., (#0)
	Comma1         *lexers.Token
	Season         *lexers.Token // e.g., Winter, Spring
	Comma2         *lexers.Token
	Weather        *lexers.Token // e.g., FINE
	// Optional next turn info
	Next           *lexers.Token   // Next (optional)
	Turn2          *lexers.Token   // Turn (optional)
	TurnYearMonth2 *lexers.Token   // e.g., 900-01 (optional)
	TurnNumber2    *TurnNumberNode // e.g., (#1) (optional)
	Comma3         *lexers.Token   // (optional)
	ReportDate     *ReportDateNode // e.g., 28/11/2025 (optional)
	EOL            *lexers.Token
	errors         []error
	tokens         []*lexers.Token
}

func (n *TurnLineNode) Errors() []error         { return n.errors }
func (n *TurnLineNode) Tokens() []*lexers.Token { return n.tokens }

// TurnNumberNode represents a turn number (e.g., (#0), (#1)).
type TurnNumberNode struct {
	LeftParen  *lexers.Token
	Hash       *lexers.Token
	Number     *lexers.Token
	RightParen *lexers.Token
	errors     []error
	tokens     []*lexers.Token
}

func (n *TurnNumberNode) Errors() []error         { return n.errors }
func (n *TurnNumberNode) Tokens() []*lexers.Token { return n.tokens }

// ReportDateNode represents a report date (DD/MM/YYYY).
type ReportDateNode struct {
	Day    *lexers.Token
	Slash1 *lexers.Token
	Month  *lexers.Token
	Slash2 *lexers.Token
	Year   *lexers.Token
	errors []error
	tokens []*lexers.Token
}

func (n *ReportDateNode) Errors() []error         { return n.errors }
func (n *ReportDateNode) Tokens() []*lexers.Token { return n.tokens }
