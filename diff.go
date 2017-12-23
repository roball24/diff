package diff

import (
	"bufio"
	"bytes"
	"io"
)

// Writer is an io.Writer
type Writer = io.Writer

// Reader is an io.Reader
type Reader = io.Reader

type bufReader = bufio.Reader

type readerChan struct {
	reader *bufReader
	ch     chan readerChanData
	delim  byte
}

type readerChanData struct {
	line []byte
	num  int
}

// FullEqual checks condition lines exactly the same
var FullEqual = bytes.Equal

// Checker can DiffCheck delimited io.Reader streams
type Checker interface {
	AddLineCompare(handler TwoLineCheckHandler, ft FailType) (id int)
	RemoveLineCompare(id int)
	AddLineIgnore(handler LineCheckHandler, ft FailType) (id int)
	RemoveLineIgnore(id int)
	AddLineCompletionHandler(handler TwoLineNumHandler) (id int)
	RemoveLineCompletionHandler(id int)
	SetWriter(w Writer)
	Delimiters(delim1, delim2 byte)
	Run() (equal bool)
}

// WriteMode configures Checker writes to Writer
type WriteMode int

const (
	// BasicWriteMode stub
	BasicWriteMode WriteMode = 1
)

// FailType configures Checker handler failure responses
type FailType int

const (
	// SoftFail will do nothing and continue to Run Checker
	SoftFail FailType = 1
	// WarningFail will write all errors to current Writer and continue to Run Checker
	WarningFail FailType = 2
	// ErrorFail will write first failure to current Writer and terminate Checker's Run
	ErrorFail FailType = 3
)

// LineCheckHandler receives a line, returns a bool
type LineCheckHandler = func(line []byte) bool

// TwoLineCheckHandler receives two lines, returns a bool
type TwoLineCheckHandler = func(line1, line2 []byte) bool

// TwoLineHandler receives two lines
type TwoLineHandler = func(chk Checker, line1, line2 []byte)

// TwoLineNumHandler receives to line numbers
type TwoLineNumHandler = func(chk Checker, lineNum1, lineNum2 int)

// LineCompare is a comparision handler and FailType
type LineCompare struct {
	handler TwoLineCheckHandler
	ft      FailType
}

// LineIgnore is an ignore handler and FailType
type LineIgnore struct {
	handler LineCheckHandler
	ft      FailType
}

type diff struct {
	description  string
	lineNum      int
	line1, line2 []byte
}
