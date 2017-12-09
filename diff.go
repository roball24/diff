package diff

import "io"

type Writer = io.Writer  // Writer is an io.Writer
type Reader = io.Reader  // Reader is an io.Reader

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

type WriteMode int64  // WriteMode configures Checker writes to Writer
const (
    BasicWriteMode WriteMode = 1
)

type FailType int64  // FailType configures Checker handler failure responses
const (
    SoftFail     FailType = 1  // SoftFail will do nothing and continue to Run Checker
    WarningFail  FailType = 2  // WarningFail will write errors to current Writer and continue to Run Checker
    ErrorFail    FailType = 3  // ErrorFail will write to current Writer and terminate Checker's Run
)

type LineCheckHandler = func(line []byte) bool  // LineCheckHandler receives a line, returns a bool
type TwoLineCheckHandler = func(baseLine, diffLine []byte) bool  // TwoLineCheckHandler receives two lines, returns a bool
type TwoLineHandler = func(baseLine, diffLine []byte)  // TwoLineHandler receives two lines
type TwoLineNumHandler = func(baseLineNum, diffLineNum int64)  // TwoLineNumHandler receives to line numbers
