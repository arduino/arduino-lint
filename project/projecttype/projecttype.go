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

// Matches returns whether the receiver project type matches the argument project type
func (projectTypeA Type) Matches(projectTypeB Type) bool {
	if projectTypeA == Not && projectTypeB == Not {
		return true
	} else if projectTypeA == Not || projectTypeB == Not {
		return false
	}
	return (projectTypeA == All || projectTypeB == All || projectTypeA == projectTypeB)
}
