package projecttype

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
