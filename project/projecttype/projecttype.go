// Package projecttype defines the Arduino project types.
package projecttype

// Type is the type for Arduino project types.
//go:generate stringer -type=Type -linecomment
type Type int

const (
	Sketch       Type = iota // sketch
	Library                  // library
	Platform                 // boards platform
	PackageIndex             // Boards Manager package index
	All                      // any project type
	Not                      // N/A
)
