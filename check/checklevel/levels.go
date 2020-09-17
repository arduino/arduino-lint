package checklevel

//go:generate stringer -type=Level -linecomment
type Level int

// Line comments set the string for each level
const (
	Info    Level = iota // info
	Warning              // warning
	Error                // error
	Pass                 // pass
)
