package widgets

import (
	"fmt"
	"unicode/utf8"
)

// TextField implements item interface
type TextField struct {
	question          string
	input             string
	cursorPosition    int
	minCursorPosition int
	visible           bool
}

var _ FormItem = (*TextField)(nil)

// NewTextField creates a new instance of TextField object
func NewTextField(question string) *TextField {
	return &TextField{
		question:          question,
		input:             "",
		minCursorPosition: 0,
		cursorPosition:    0,
	}
}

func (t *TextField) setVisible(visible bool) {
	t.visible = visible
}

func (t *TextField) string() string {
	return fmt.Sprintf("%s %s", t.question, t.input)
}

func (t *TextField) setCursorPosition() {
	if t.cursorPosition > utf8.RuneCountInString(t.input) {
		t.cursorPosition = utf8.RuneCountInString(t.input)
	}
	if t.cursorPosition < 0 {
		t.cursorPosition = 0
	}
}

func (t *TextField) handleInput(e formEvent) {
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

func (t *TextField) selectable() bool { return true }

func (t *TextField) Answer() string {
	return t.input
}
