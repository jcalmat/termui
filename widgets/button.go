package widgets

import "strings"

// Button implements item interface
type Button struct {
	label    string
	visible  bool
	callback func()
}

var _ FormItem = (*Button)(nil)

// NewButton creates a new instance of Button object
func NewButton(label string, callback func()) *Button {
	return &Button{
		label:    label,
		callback: callback,
	}
}

func (c *Button) string() string {
	var sb strings.Builder
	sb.WriteRune('[')
	sb.WriteString(c.label)
	sb.WriteRune(']')
	sb.WriteRune(' ')
	return sb.String()
}

func (c *Button) handleInput(e formEvent) {
	if e == enter {
		c.callback()
	}
}

func (c *Button) setVisible(visible bool) {
	c.visible = visible
}

func (c *Button) selectable() bool { return true }
