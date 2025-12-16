// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cst

import (
	"fmt"
	"io"

	"github.com/playbymail/ottoapp/backend/services/reports/lexers"
)

type PrintOption func(*printer)

func WithColumnNo() PrintOption {
	return func(p *printer) {
		p.showColNo = true
	}
}

func WithLineNo() PrintOption {
	return func(p *printer) {
		p.showLineNo = true
	}
}

func WithName(name string) PrintOption {
	return func(p *printer) {
		p.name = name
	}
}

func WithSource(src []byte) PrintOption {
	return func(p *printer) {
		p.src = src
	}
}

func WithTokenKind() PrintOption {
	return func(p *printer) {
		p.showTokenKind = true
	}
}

func PrettyPrint(w io.Writer, t *TurnReportNode, opts ...PrintOption) {
	p := &printer{w: w}
	for _, opt := range opts {
		opt(p)
	}

	if p.name != "" {
		p.printf("// source %q\n", p.name)
	}

	p.printf("TurnReport {\n")
	p.indent++

	for i, section := range t.Sections {
		p.printUnitSection(i, section)
	}

	if t.EOF != nil {
		p.printToken("EOF", t.EOF)
	}
	p.printErrors(t.errors)

	p.indent--
	p.printf("}\n")
}

type printer struct {
	w             io.Writer
	indent        int
	name          string
	showColNo     bool
	showLineNo    bool
	showTokenKind bool
	src           []byte
}

func (p *printer) printf(format string, args ...any) {
	for i := 0; i < p.indent; i++ {
		_, _ = p.w.Write([]byte("  "))
	}
	_, _ = p.w.Write([]byte(fmt.Sprintf(format, args...)))
}

func (p *printer) writef(format string, args ...any) {
	_, _ = p.w.Write([]byte(fmt.Sprintf(format, args...)))
}

func (p *printer) printLineNo(lineNo int) {
	if p.showLineNo {
		p.printf("Line: %d\n", lineNo)
	}
}

func (p *printer) printUnitSection(index int, n *UnitSectionNode) {
	p.printf("UnitSection[%d] {\n", index)
	p.indent++
	p.printLineNo(n.LineNo)

	if n.UnitLine != nil {
		p.printUnitLine(n.UnitLine)
	} else {
		p.printf("UnitLine: <nil>\n")
	}

	if n.TurnLine != nil {
		p.printTurnLine(n.TurnLine)
	} else {
		p.printf("TurnLine: <nil>\n")
	}

	if n.UnitMovementLine != nil {
		p.printUnitMovementLine(n.UnitMovementLine)
	}

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printUnitLine(n *UnitLineNode) {
	p.printf("UnitLine {\n")
	p.indent++

	p.printToken("Keyword", n.Keyword)
	p.printToken("UnitID", n.UnitID)
	p.printToken("Note", n.Note)
	p.printCoords("CurrentHex", n.CurrentHex)
	p.printCoords("PreviousHex", n.PreviousHex)

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printTurnLine(n *TurnLineNode) {
	p.printf("TurnLine {\n")
	p.indent++

	p.printToken("TurnYearMonth", n.TurnYearMonth1)
	if n.TurnNumber1 != nil {
		p.printTurnNumber("TurnNumber", n.TurnNumber1)
	}
	p.printToken("Season", n.Season)
	p.printToken("Weather", n.Weather)

	if n.Next != nil {
		p.printToken("NextTurnYearMonth", n.TurnYearMonth2)
		if n.TurnNumber2 != nil {
			p.printTurnNumber("NextTurnNumber", n.TurnNumber2)
		}
		if n.ReportDate != nil {
			p.printReportDate(n.ReportDate)
		}
	}

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printTurnNumber(label string, n *TurnNumberNode) {
	if n.Number != nil {
		line, col := n.Number.Position()
		p.printLabelValue(line, col, label, fmt.Sprintf("#%s", string(n.Number.Bytes())), n.Number.Kind.String())
	}
}

func (p *printer) printReportDate(n *ReportDateNode) {
	var line, col int
	dayValue, dayKind := "<nil>", "<nil>"
	if n.Day != nil {
		line, col = n.Day.Position()
		dayValue, dayKind = string(n.Day.Bytes()), n.Day.Kind.String()
	}
	monthValue, monthKind := "<nil>", "<nil>"
	if n.Month != nil {
		if line == 0 && col == 0 {
			line, col = n.Month.Position()
		}
		monthValue, monthKind = string(n.Month.Bytes()), n.Month.Kind.String()
	}
	yearValue, yearKind := "<nil>", "<nil>"
	if n.Year != nil {
		if line == 0 && col == 0 {
			line, col = n.Year.Position()
		}
		yearValue, yearKind = string(n.Year.Bytes()), n.Year.Kind.String()
	}
	value := fmt.Sprintf("%s/%s/%s", dayValue, monthValue, yearValue)
	kind := fmt.Sprintf("%s %s %s", dayKind, monthKind, yearKind)
	p.printLabelValue(line, col, "Report Date", value, kind)
}

func (p *printer) printUnitMovementLine(n UnitMovementLineNode) {
	switch v := n.(type) {
	case *UnitGoesToLineNode:
		p.printUnitGoesToLine(v)
	case *LandMovementLineNode:
		p.printLandMovementLine(v)
	default:
		p.printf("UnitMovementLine: <unknown type>\n")
	}
}

func (p *printer) printUnitGoesToLine(n *UnitGoesToLineNode) {
	p.printf("UnitGoesToLine {\n")
	p.indent++
	p.printLineNo(n.LineNo)

	p.printToken("Tribe", n.Tribe)
	if n.Coords != nil {
		p.printCoords("Coords", n.Coords)
	}

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printLandMovementLine(n *LandMovementLineNode) {
	p.printf("LandMovementLine {\n")
	p.indent++

	p.printToken("Tribe", n.Tribe)
	if n.LandMovement != nil {
		p.printLandMovement(n.LandMovement)
	}

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printLandMovement(n *LandMovementNode) {
	p.printf("LandMovement {\n")
	p.indent++

	p.printToken("Move", n.Move)
	for i, step := range n.Steps {
		p.printLandStep(i, step)
	}

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printLandStep(index int, n *LandStepNode) {
	if n.IsEmpty() {
		p.printf("Step[%d]: <empty>\n", index)
		return
	}

	p.printf("Step[%d] {\n", index)
	p.indent++

	p.printToken("Direction", n.Direction)
	p.printToken("Terrain", n.Terrain)

	p.printErrors(n.errors)

	p.indent--
	p.printf("}\n")
}

func (p *printer) printLabelValue(line, col int, label, value, kind string) {
	p.printf("%s: %q", label, value)
	if p.showLineNo || p.showColNo || p.showTokenKind {
		p.writef("  (")
		if p.showLineNo || p.showColNo {
			if !p.showColNo {
				p.writef("%d", line)
			} else {
				p.writef("%d:%d", line, col)
			}
			if p.showTokenKind {
				p.writef(":%s", kind)
			}
		} else {
			p.writef("%s", kind)
		}
		p.writef(")")
	}
	p.writef("\n")
}

func (p *printer) printToken(label string, tok *lexers.Token) {
	if tok == nil {
		return
	}
	line, col := tok.Position()
	p.printLabelValue(line, col, label, string(tok.Bytes()), tok.Kind.String())
}

func (p *printer) printCoords(label string, c CoordsNode) {
	if c == nil {
		p.printf("%s: <nil>\n", label)
		return
	}
	var line, col int
	var value, kind string
	switch v := c.(type) {
	case *GridCoordsNode:
		if v.Grid != nil && v.Number != nil {
			line, col = v.Grid.Position()
			value = fmt.Sprintf("%s %s", string(v.Grid.Bytes()), string(v.Number.Bytes()))
			kind = fmt.Sprintf("%s %s", v.Grid.Kind.String(), v.Number.Kind.String())
		} else if v.Grid != nil {
			line, col = v.Grid.Position()
			value = fmt.Sprintf("%s <nil>", string(v.Grid.Bytes()))
			kind = fmt.Sprintf("%s <nil>", v.Grid.Kind.String())
		} else if v.Number != nil {
			line, col = v.Number.Position()
			value = fmt.Sprintf("<nil> %s", string(v.Number.Bytes()))
			kind = fmt.Sprintf("<nil> %s", v.Number.Kind.String())
		}
	case *NACoordsNode:
		if v.Text != nil {
			line, col = v.Text.Position()
			value = string(v.Text.Bytes())
			kind = v.Text.Kind.String()
		} else {
			value, kind = "<nil>", "<nil>"
		}
	case *ObscuredCoordsNode:
		if v.Grid != nil && v.Number != nil {
			line, col = v.Grid.Position()
			value = fmt.Sprintf("%s %s", string(v.Grid.Bytes()), string(v.Number.Bytes()))
			kind = fmt.Sprintf("%s %s", v.Grid.Kind.String(), v.Number.Kind.String())
		} else if v.Grid != nil {
			line, col = v.Grid.Position()
			value = fmt.Sprintf("%s <nil>", string(v.Grid.Bytes()))
			kind = fmt.Sprintf("%s <nil>", v.Grid.Kind.String())
		} else if v.Number != nil {
			line, col = v.Number.Position()
			value = fmt.Sprintf("<nil> %s", string(v.Number.Bytes()))
			kind = fmt.Sprintf("<nil> %s", v.Number.Kind.String())
		}
	case *ErrorCoordsNode:
		//p.printf("%s: <error: %s>\n", label, v.Message)
		value, kind = fmt.Sprintf("<error: %s>", v.Message), "<errorCoordsNode>"
	default:
		//p.printf("%s: <unknown>\n", label)
		value, kind = "<unknown>", fmt.Sprintf("<%T>", v)
	}
	p.printLabelValue(line, col, label, value, kind)
}

func (p *printer) printErrors(errors []error) {
	for _, err := range errors {
		p.printf("error: %q\n", err.Error())
	}
}
