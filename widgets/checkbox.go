package widgets

import (
	"strings"
)

const (
	checkbox_uncheck string = "☐"
	checkbox_check   string = "☑"
)

// Checkbox implements item interface
type Checkbox struct {
	question string
	checked  bool
	visible  bool
}

var _ FormItem = (*Checkbox)(nil)

// NewCheckbox creates a new instance of Checkbox object
func NewCheckbox(question string, checked bool) *Checkbox {
	return &Checkbox{
		question: question,
		checked:  checked,
	}
}

func (c *Checkbox) string() string {
	var sb strings.Builder
	if c.checked {
		sb.WriteString(checkbox_check)
	} else {
		sb.WriteString(checkbox_uncheck)
	}
	sb.WriteString(" ")
	sb.WriteString(c.question)
	return sb.String()
}

func (c *Checkbox) handleInput(e formEvent) {
	if e == enter {
		c.checked = !c.checked
	}
}

func (c *Checkbox) setVisible(visible bool) {
	c.visible = visible
}

func (c *Checkbox) selectable() bool { return true }

func (c *Checkbox) Answer() bool {
	return c.checked && c.visible
}
