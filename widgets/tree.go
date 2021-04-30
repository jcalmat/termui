package widgets

import (
	"fmt"
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
)

const treeIndent = "  "

// TreeNode is a tree node.
type TreeNode struct {
	Value    fmt.Stringer
	Expanded bool
	Nodes    []*TreeNode

	// level stores the node level in the tree.
	level int
}

// TreeWalkFn is a function used for walking a Tree.
// To interrupt the walking process function should return false.
type TreeWalkFn func(*TreeNode) bool

func (self *TreeNode) parseStyles(style Style) []Cell {
	var sb strings.Builder
	if len(self.Nodes) == 0 {
		sb.WriteString(strings.Repeat(treeIndent, self.level+1))
	} else {
		sb.WriteString(strings.Repeat(treeIndent, self.level))
		if self.Expanded {
			sb.WriteRune(Theme.Tree.Expanded)
		} else {
			sb.WriteRune(Theme.Tree.Collapsed)
		}
		sb.WriteByte(' ')
	}
	sb.WriteString(self.Value.String())
	return ParseStyles(sb.String(), style)
}

// Tree is a tree widget.
type Tree struct {
	Block
	TextStyle        Style
	SelectedRowStyle Style
	WrapText         bool
	selectedRow      int

	nodes []*TreeNode
	// rows is flatten nodes for rendering.
	rows   []*TreeNode
	topRow int
}

// NewTree creates a new Tree widget.
func NewTree() *Tree {
	return &Tree{
		Block:            *NewBlock(),
		TextStyle:        Theme.Tree.Text,
		SelectedRowStyle: Theme.Tree.Text,
		WrapText:         true,
	}
}

func (self *Tree) SetNodes(nodes []*TreeNode) {
	self.nodes = nodes
	self.prepareNodes()
}

func (self *Tree) prepareNodes() {
	self.rows = make([]*TreeNode, 0)
	for _, node := range self.nodes {
		self.prepareNode(node, 0)
	}
}

func (self *Tree) prepareNode(node *TreeNode, level int) {
	self.rows = append(self.rows, node)
	node.level = level

	if node.Expanded {
		for _, n := range node.Nodes {
			self.prepareNode(n, level+1)
		}
	}
}

func (self *Tree) Walk(fn TreeWalkFn) {
	for _, n := range self.nodes {
		if !self.walk(n, fn) {
			break
		}
	}
}

func (self *Tree) walk(n *TreeNode, fn TreeWalkFn) bool {
	if !fn(n) {
		return false
	}

	for _, node := range n.Nodes {
		if !self.walk(node, fn) {
			return false
		}
	}

	return true
}

func (self *Tree) Draw(buf *Buffer) {
	self.Block.Draw(buf)
	point := self.Inner.Min

	// adjusts view into widget
	if self.selectedRow >= self.Inner.Dy()+self.topRow {
		self.topRow = self.selectedRow - self.Inner.Dy() + 1
	} else if self.selectedRow < self.topRow {
		self.topRow = self.selectedRow
	}

	// draw rows
	for row := self.topRow; row < len(self.rows) && point.Y < self.Inner.Max.Y; row++ {
		cells := self.rows[row].parseStyles(self.TextStyle)
		if self.WrapText {
			cells = WrapCells(cells, uint(self.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < self.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == self.selectedRow {
				style = self.SelectedRowStyle
			}
			if point.X+1 == self.Inner.Max.X+1 && len(cells) > self.Inner.Dx() {
				buf.SetCell(NewCell(ELLIPSES, style), point.Add(image.Pt(-1, 0)))
			} else {
				buf.SetCell(NewCell(cells[j].Rune, style), point)
				point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
			}
		}
		point = image.Pt(self.Inner.Min.X, point.Y+1)
	}

	// draw UP_ARROW if needed
	if self.topRow > 0 {
		buf.SetCell(
			NewCell(UP_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(self.rows) > int(self.topRow)+self.Inner.Dy() {
		buf.SetCell(
			NewCell(DOWN_ARROW, NewStyle(ColorWhite)),
			image.Pt(self.Inner.Max.X-1, self.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set self.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (self *Tree) ScrollAmount(amount int) {
	if len(self.rows)-int(self.selectedRow) <= amount {
		self.selectedRow = len(self.rows) - 1
	} else if int(self.selectedRow)+amount < 0 {
		self.selectedRow = 0
	} else {
		self.selectedRow += amount
	}
}

func (self *Tree) SelectedNode() *TreeNode {
	if len(self.rows) == 0 {
		return nil
	}
	return self.rows[self.selectedRow]
}

func (self *Tree) ScrollUp() {
	self.ScrollAmount(-1)
}

func (self *Tree) ScrollDown() {
	self.ScrollAmount(1)
}

func (self *Tree) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if self.selectedRow > self.topRow {
		self.selectedRow = self.topRow
	} else {
		self.ScrollAmount(-self.Inner.Dy())
	}
}

func (self *Tree) ScrollPageDown() {
	self.ScrollAmount(self.Inner.Dy())
}

func (self *Tree) ScrollHalfPageUp() {
	self.ScrollAmount(-int(FloorFloat64(float64(self.Inner.Dy()) / 2)))
}

func (self *Tree) ScrollHalfPageDown() {
	self.ScrollAmount(int(FloorFloat64(float64(self.Inner.Dy()) / 2)))
}

func (self *Tree) ScrollTop() {
	self.selectedRow = 0
}

func (self *Tree) ScrollBottom() {
	self.selectedRow = len(self.rows) - 1
}

func (self *Tree) Collapse() {
	self.rows[self.selectedRow].Expanded = false
	self.prepareNodes()
}

func (self *Tree) Expand() {
	node := self.rows[self.selectedRow]
	if len(node.Nodes) > 0 {
		self.rows[self.selectedRow].Expanded = true
	}
	self.prepareNodes()
}

func (self *Tree) ToggleExpand() {
	node := self.rows[self.selectedRow]
	if len(node.Nodes) > 0 {
		node.Expanded = !node.Expanded
	}
	self.prepareNodes()
}

func (self *Tree) ExpandAll() {
	self.Walk(func(n *TreeNode) bool {
		if len(n.Nodes) > 0 {
			n.Expanded = true
		}
		return true
	})
	self.prepareNodes()
}

func (self *Tree) CollapseAll() {
	self.Walk(func(n *TreeNode) bool {
		n.Expanded = false
		return true
	})
	self.prepareNodes()
}
