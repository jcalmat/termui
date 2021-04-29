package widgets

import (
	"strings"

	ui "github.com/gizak/termui/v3"
)

const (
	checkbox_uncheck string = "☐"
	checkbox_check   string = "☑"
)

// checkbox implements item interface
type checkbox struct {
	question string
	prefix   string
	checked  bool
	visible  bool
}

var _ FormItem = (*checkbox)(nil)

// NewCheckbox creates a new instance of checkbox object
func NewCheckbox(question string, checked bool) *checkbox {
	return &checkbox{
		prefix:   "",
		question: question,
		checked:  checked,
	}
}

func (c *checkbox) string() string {
	var sb strings.Builder
	sb.WriteString(c.prefix)
	if c.checked {
		sb.WriteString(checkbox_check)
	} else {
		sb.WriteString(checkbox_uncheck)
	}
	sb.WriteString(" ")
	sb.WriteString(c.question)
	return sb.String()
}

func (c *checkbox) handleInput(e ui.Event) {
	if e.ID == "<Enter>" {
		c.toggle()
	}
}

func (c *checkbox) setVisible(visible bool) {
	c.visible = visible
}

func (c *checkbox) toggle() {
	c.checked = !c.checked
}

func (c *checkbox) selectable() bool { return true }

func (c *checkbox) Answer() bool {
	return c.checked && c.visible
}

// func (c *checkbox) setPrefix(prefix string) {
// 	c.prefix = prefix
// }

// func (c *checkbox) clearValue() {
// 	c.checked = false
// }

// func (c *checkbox) setCursorPosition() {}

// func (c *checkbox) ToFormItem() *FormItem {
// 	return &FormItem{item: c}
// }

// func (c *checkbox) write() {
// 	// var question string

// 	checkbox := checkbox_uncheck
// 	if c.checked {
// 		checkbox = "\u001b[32;1m" + checkbox_check
// 	}
// 	if c.selected {
// 		checkbox = "\u001b[7m" + checkbox
// 	}

// 	// question = fmt.Sprintf("%s%s %s\u001b[0m", c.prefix, checkbox, c.question)

// 	// clearLine()
// 	// write(question)
// 	// cursor.MoveColumn(1)
// }

// func (c *checkbox) pick() {
// 	c.selected = true
// 	// cursor.HideCursor()
// }

// func (c *checkbox) unpick() {
// 	c.selected = false
// 	// cursor.DisplayCursor()
// }

// func (c *checkbox) displayChildren() bool {
// 	return c.checked
// }
