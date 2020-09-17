package checkresult

//go:generate stringer -type=Result -linecomment
type Result int

const (
	Pass Result = iota // pass
	Fail               // fail
	// The check is configured to be skipped in the current tool configuration mode
	Skipped // skipped
	// An error prevented the check from running
	NotRun // not run
)
