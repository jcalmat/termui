package widgets

import (
	"strings"
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

func (c *checkbox) handleInput(e formEvent) {
	if e == enter {
		c.checked = !c.checked
	}
}

func (c *checkbox) setVisible(visible bool) {
	c.visible = visible
}

func (c *checkbox) selectable() bool { return true }

func (c *checkbox) Answer() bool {
	return c.checked && c.visible
}
