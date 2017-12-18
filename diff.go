package diff

import "io"

// Writer is an io.Writer
type Writer = io.Writer

// Reader is an io.Reader
type Reader = io.Reader

// Checker can DiffCheck delimited io.Reader streams
type Checker interface {
	AddLineCompare(handler TwoLineCheckHandler, ft FailType) (id int64)
	RemoveLineCompare(id int64)
	AddLineIgnore(handler LineCheckHandler, ft FailType) (id int64)
	RemoveLineIgnore(id int64)
	Writer(w Writer)
	LineCompletionHandler(handler TwoLineNumHandler)
	Delimiters(baseDelim, diffDelim byte)
	Run() error
	Reader
}

// WriteMode configures Checker writes to Writer
type WriteMode int64

const (
	// BasicWriteMode stub
	BasicWriteMode WriteMode = 1
)

// FailType configures Checker handler failure responses
type FailType int64

const (
	// SoftFail will do nothing and continue to Run Checker
	SoftFail FailType = 1
	// WarningFail will write errors to current Writer and continue to Run Checker
	WarningFail FailType = 2
	// ErrorFail will write to current Writer and terminate Checker's Run
	ErrorFail FailType = 3
)

// LineCheckHandler receives a line, returns a bool
type LineCheckHandler = func(line []byte) bool

// TwoLineCheckHandler receives two lines, returns a bool
type TwoLineCheckHandler = func(baseLine, diffLine []byte) bool

// TwoLineHandler receives two lines
type TwoLineHandler = func(baseLine, diffLine []byte)

// TwoLineNumHandler receives to line numbers
type TwoLineNumHandler = func(baseLineNum, diffLineNum int64)
