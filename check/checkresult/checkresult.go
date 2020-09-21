package checkresult

//go:generate stringer -type=Type -linecomment
type Type int

const (
	Pass Type = iota // pass
	Fail             // fail
	// The check is configured to be skipped in the current tool configuration mode
	Skip // skipped
	// An unrelated error prevented the check from running
	NotRun // unable to run
)
