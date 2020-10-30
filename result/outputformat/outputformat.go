// Package projecttype defines the output formats
package outputformat

// Type is the type for output formats
//go:generate stringer -type=Type -linecomment
type Type int

const (
	Text Type = iota // text
	JSON             // JSON
)
