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

	lineCompares           map[int]interface{}
	lineIgnores            map[int]interface{}
	lineCompletionHandlers map[int]interface{}

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
		wr:                     os.Stdout,
		lineCompares:           make(map[int]interface{}, 0),
		lineIgnores:            make(map[int]interface{}, 0),
		lineCompletionHandlers: make(map[int]interface{}, 0),
	}, nil
}

// AddLineCompare adds a line compare rule
func (chk *SingleChecker) AddLineCompare(handler TwoLineCheckHandler, ft FailType) (id int) {
	return insertAtRandomKey(
		chk.lineCompares,
		LineCompare{handler: handler, ft: ft},
	)
}

// RemoveLineCompare removes a line compare rule
func (chk *SingleChecker) RemoveLineCompare(id int) {
	delete(chk.lineCompares, id)
}

// AddLineIgnore adds a line ignore rule
func (chk *SingleChecker) AddLineIgnore(handler LineCheckHandler, ft FailType) (id int) {
	return insertAtRandomKey(
		chk.lineIgnores,
		LineIgnore{handler: handler, ft: ft},
	)
}

// RemoveLineIgnore removes a line ignore rule
func (chk *SingleChecker) RemoveLineIgnore(id int) {
	delete(chk.lineIgnores, id)
}

// AddLineCompletionHandler adds a line handler
func (chk *SingleChecker) AddLineCompletionHandler(handler TwoLineNumHandler) (id int) {
	return insertAtRandomKey(
		chk.lineCompletionHandlers,
		handler,
	)
}

// RemoveLineCompletionHandler removes a line handler
func (chk *SingleChecker) RemoveLineCompletionHandler(id int) {
	delete(chk.lineCompletionHandlers, id)
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
		// check line compares
		for _, lcInt := range chk.lineCompares {
			lineComp := lcInt.(LineCompare)
			if !lineComp.handler(rcData1.line, rcData2.line) {
				return false
			}
		}
		// TODO: line completion handler
	}
}

func insertAtRandomKey(m map[int]interface{}, val interface{}) (id int) {
	ok := true
	for ok {
		id, ok = func() (int, bool) {
			id = rand.Int()
			_, ok = m[id]
			return id, ok
		}()
	}
	m[id] = val
	return id
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
