package widgets

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
)

const formIndent = "  "

type formEvent string

const (
	left  formEvent = "left"
	right formEvent = "right"
	del   formEvent = "del"
	enter formEvent = "enter"
	space formEvent = " "
)

var formEventMap = map[string]formEvent{
	"<Left>":      left,
	"<Right>":     right,
	"<Backspace>": del,
	"<Enter>":     enter,
	"<Space>":     space,
}

func (f formEvent) String() string {
	return string(f)
}

type FormItem interface {
	// handleInput forwards events to formItems
	handleInput(formEvent)

	// selectable indicates if the item should be selectable of if it should
	// be skipped when navigating in the item list
	selectable() bool

	// setVisible set the item visibility
	setVisible(bool)

	// string returns the stringified item
	string() string
}

// FormNode is a form node.
type FormNode struct {
	Item     FormItem
	Expanded bool
	Nodes    []*FormNode

	// level stores the node level in the form.
	level int
}

// FormWalkFn is a function used for walking a Form.
// To interrupt the walking process function should return false.
type FormWalkFn func(*FormNode) bool

func (self *FormNode) parseStyles(style Style) []Cell {
	var sb strings.Builder

	sb.WriteString(strings.Repeat(formIndent, self.level))

	if len(self.Nodes) == 0 {
	} else {
		// sb.WriteString(strings.Repeat(formIndent, self.level))
		if self.Expanded {
			sb.WriteRune(Theme.Form.Expanded)
		} else {
			sb.WriteRune(Theme.Form.Collapsed)
		}
		sb.WriteByte(' ')
	}
	sb.WriteString(self.Item.string())
	return ParseStyles(sb.String(), style)
}

// Form is a form widget.
type Form struct {
	Block
	TextStyle         Style
	SelectedTextStyle Style
	WrapText          bool
	selectedRow       int

	nodes []*FormNode
	// rows is flatten nodes for rendering.
	rows []*FormNode
	// visibleRows is flatten nodes used for visibility assignment
	visibleRows map[*FormNode]bool
	topRow      int
}

// NewForm creates a new Form widget.
func NewForm() *Form {
	return &Form{
		Block:             *NewBlock(),
		TextStyle:         Theme.Form.Text,
		SelectedTextStyle: Theme.Form.Text,
		WrapText:          true,
	}
}

func (self *Form) initVisibilityMap(node *FormNode) {
	self.visibleRows[node] = false

	for _, node := range node.Nodes {
		self.initVisibilityMap(node)
	}
}

func (self *Form) SetNodes(nodes []*FormNode) {
	self.nodes = nodes
	self.visibleRows = make(map[*FormNode]bool)

	for _, node := range self.nodes {
		self.initVisibilityMap(node)
	}
	self.prepareNodes()
}

func (self *Form) prepareNodes() {
	self.rows = make([]*FormNode, 0)

	// reset visibility for every node
	for row := range self.visibleRows {
		self.visibleRows[row] = false
		row.Item.setVisible(false)
	}

	for _, node := range self.nodes {
		self.prepareNode(node, 0)
	}
}

func (self *Form) prepareNode(node *FormNode, level int) {
	self.rows = append(self.rows, node)
	node.level = level
	node.Item.setVisible(true)

	if node.Expanded {
		for _, n := range node.Nodes {
			self.prepareNode(n, level+1)
		}
	}
}

func (self *Form) Walk(fn FormWalkFn) {
	for _, n := range self.nodes {
		if !self.walk(n, fn) {
			break
		}
	}
}

func (self *Form) walk(n *FormNode, fn FormWalkFn) bool {
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

func (self *Form) Draw(buf *Buffer) {
	self.Block.Draw(buf)
	point := self.Inner.Min

	if !self.rows[self.selectedRow].Item.selectable() {
		self.ScrollAmount(1)
	}

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
				style = self.SelectedTextStyle
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
func (self *Form) ScrollAmount(amount int) {
	for {
		if len(self.rows)-int(self.selectedRow) <= amount {
			self.selectedRow = 0
		} else if int(self.selectedRow)+amount < 0 {
			self.selectedRow = len(self.rows) - 1
		} else {
			self.selectedRow += amount
		}
		if self.rows[self.selectedRow].Item.selectable() {
			break
		}
	}
}

func (self *Form) SelectedNode() *FormNode {
	if len(self.rows) == 0 {
		return nil
	}
	return self.rows[self.selectedRow]
}

func (self *Form) ScrollUp() {
	self.ScrollAmount(-1)
}

func (self *Form) ScrollDown() {
	self.ScrollAmount(1)
}

func (self *Form) ToggleExpand() {
	node := self.rows[self.selectedRow]
	if len(node.Nodes) > 0 {
		node.Expanded = !node.Expanded
	}
	self.prepareNodes()
}

// HandleKeyboard handle special events that don't need to mapped by hand.
func (self *Form) HandleKeyboard(e Event) {
	if e.Type != KeyboardEvent {
		return
	}

	s := e.ID
	node := self.rows[self.selectedRow]
	if e, ok := formEventMap[s]; ok {
		node.Item.handleInput(e)
		return
	}

	node.Item.handleInput(formEvent(s))
}
