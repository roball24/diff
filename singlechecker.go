package diff

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"os"
	"time"
)

// SingleChecker diff checks input from two readers
type SingleChecker struct {
	rch1 *readerChan
	rch2 *readerChan
	wr   Writer

	conditions  map[int]interface{}
	ignores     map[int]interface{}
	breakpoints map[int]interface{}

	diffs []*diff
}

// NewSingleChecker initializes a SingleChecker from io.Readers
func NewSingleChecker(rd1, rd2 Reader) (*SingleChecker, error) {
	if rd1 == nil || rd2 == nil {
		return nil, errors.New("Readers must not be nil")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	return &SingleChecker{
		rch1: &readerChan{
			reader: bufio.NewReader(rd1),
			delim:  '\n',
			ch:     nil,
		},
		rch2: &readerChan{
			reader: bufio.NewReader(rd2),
			delim:  '\n',
			ch:     nil,
		},
		wr:          os.Stdout,
		conditions:  make(map[int]interface{}, 0),
		ignores:     make(map[int]interface{}, 0),
		breakpoints: make(map[int]interface{}, 0),
	}, nil
}

func insertAtRandomKey(m map[int]interface{}, val interface{}) (id int) {
	ok := true
	for ok {
		id = rand.Int()
		_, ok = m[id]
	}
	m[id] = val
	return id
}

// AddCondition adds a line compare rule
func (chk *SingleChecker) AddCondition(handler TwoLineCheckHandler, ft FailType) (id int) {
	return insertAtRandomKey(
		chk.conditions,
		Condition{handler: handler, ft: ft},
	)
}

// RemoveCondition removes a line compare rule
func (chk *SingleChecker) RemoveCondition(id int) {
	delete(chk.conditions, id)
}

// AddIgnore adds a line ignore rule
func (chk *SingleChecker) AddIgnore(handler LineCheckHandler, ft FailType) (id int) {
	return insertAtRandomKey(
		chk.ignores,
		Ignore{handler: handler, ft: ft},
	)
}

// RemoveIgnore removes a line ignore rule
func (chk *SingleChecker) RemoveIgnore(id int) {
	delete(chk.ignores, id)
}

// AddBreakpoint adds a line handler
func (chk *SingleChecker) AddBreakpoint(handler TwoLineNumHandler) (id int) {
	return insertAtRandomKey(
		chk.breakpoints,
		handler,
	)
}

// RemoveBreakpoint removes a line handler
func (chk *SingleChecker) RemoveBreakpoint(id int) {
	delete(chk.breakpoints, id)
}

// SetWriter assigns output writer
func (chk *SingleChecker) SetWriter(wr Writer) { chk.wr = wr }

// Delimiters overwrites default reader line delimiters
func (chk *SingleChecker) Delimiters(delim1, delim2 byte) {
	chk.rch1.delim = delim1
	chk.rch2.delim = delim2
}

// Run begins diff checking reader lines
func (chk *SingleChecker) Run() (equal bool) {
	chk.rch1.ch = make(chan readerChanData)
	chk.rch2.ch = make(chan readerChanData)
	// read lines to reader channels
	// TODO: Panic on error, handle recovery
	go chk.readLines()
	// handle data from reader channels
	for {
		rcData1, ok1 := <-chk.rch1.ch
		rcData2, ok2 := <-chk.rch2.ch
		// handled EOF
		if !ok1 || !ok2 {
			return true
		}
		// TODO: breakpoint
		// check line compares
		for _, condInt := range chk.conditions {
			cond := condInt.(Condition)
			if !cond.handler(rcData1.line, rcData2.line) {
				return false
			}
		}
	}
}

func (chk *SingleChecker) readLines() {
	defer close(chk.rch1.ch)
	defer close(chk.rch2.ch)
	lnum1 := 0
	lnum2 := 0
	for {
		_, err1 := chk.read(chk.rch1, &lnum1)
		_, err2 := chk.read(chk.rch2, &lnum2)
		if err1 == io.EOF && err2 == io.EOF {
			return
		}
		// files not same length TODO: Panic
		if err1 == io.EOF {
			return
		}
		if err2 == io.EOF {
			return
		}
	}
}

func (chk *SingleChecker) read(rchn *readerChan, lineNum *int) (data []byte, err error) {
	// TODO: line ignores
	data, err = rchn.reader.ReadBytes(rchn.delim)
	(*lineNum)++
	rchn.ch <- readerChanData{line: data, num: *lineNum}
	return data, err
}
