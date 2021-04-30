package widgets

import (
	"fmt"
	"unicode/utf8"
)

// textField implements item interface
type textField struct {
	question          string
	input             string
	cursorPosition    int
	minCursorPosition int
	visible           bool
}

var _ FormItem = (*textField)(nil)

// NewTextField creates a new instance of textField object
func NewTextField(question string) *textField {
	return &textField{
		question:          question,
		input:             "",
		minCursorPosition: 0,
		cursorPosition:    0,
	}
}

func (t *textField) setVisible(visible bool) {
	t.visible = visible
}

func (t *textField) string() string {
	return fmt.Sprintf("%s %s", t.question, t.input)
}

func (t *textField) setCursorPosition() {
	if t.cursorPosition > utf8.RuneCountInString(t.input) {
		t.cursorPosition = utf8.RuneCountInString(t.input)
	}
	if t.cursorPosition < 0 {
		t.cursorPosition = 0
	}
}

func (t *textField) handleInput(e formEvent) {
	if e == right {
		t.cursorPosition++
		t.setCursorPosition()
		return
	} else if e == left {
		t.cursorPosition--
		t.setCursorPosition()
		return
	}

	if e == del {
		if t.cursorPosition > 0 {
			t.input = t.input[:t.cursorPosition-1] + t.input[t.cursorPosition:]
			t.cursorPosition--
		}
		t.setCursorPosition()
		return
	}

	// unhandled special char
	if len(e) > 1 {
		return
	}

	for _, c := range e {
		if c >= 32 && c <= 126 {
			t.input = fmt.Sprintf("%s%s%s", t.input[:t.cursorPosition], string(c), t.input[t.cursorPosition:])
			t.cursorPosition++
		}
	}
	t.setCursorPosition()
}

func (t *textField) selectable() bool { return true }

func (t *textField) Answer() string {
	return t.input
}
