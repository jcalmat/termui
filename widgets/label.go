package widgets

// label implements item interface
type label struct {
	s string
}

var _ FormItem = (*label)(nil)

// NewLabel creates a new instance of label object
func NewLabel(s string) *label {
	return &label{
		s: s,
	}
}

func (l *label) string() string {
	return l.s
}

func (l *label) handleInput(e formEvent) {}

func (l *label) selectable() bool { return false }

func (l *label) setVisible(v bool) {}
