package widgets

// Label implements item interface
type Label struct {
	s string
}

var _ FormItem = (*Label)(nil)

// NewLabel creates a new instance of Label object
func NewLabel(s string) *Label {
	return &Label{
		s: s,
	}
}

func (l *Label) string() string {
	return l.s
}

func (l *Label) handleInput(e formEvent) {}

func (l *Label) selectable() bool { return false }

func (l *Label) setVisible(v bool) {}
