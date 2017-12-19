package diff

import (
	"errors"
	"math/rand"
	"os"
	"time"
)

// SingleChecker diff checks input from two readers
type SingleChecker struct {
	rd1 Reader
	rd2 Reader
	wr  Writer

	baseDelim byte
	diffDelim byte

	lineCompares           map[int]interface{}
	lineIgnores            map[int]interface{}
	lineCompletionHandlers map[int]interface{}
}

// NewSingleChecker initializes a SingleChecker from io.Readers
func NewSingleChecker(rd1, rd2 Reader) (*SingleChecker, error) {
	if rd1 == nil || rd2 == nil {
		return nil, errors.New("Readers must not be nil")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	return &SingleChecker{
		rd1:                    rd1,
		rd2:                    rd2,
		wr:                     os.Stdout,
		baseDelim:              '\n',
		diffDelim:              '\n',
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
func (chk *SingleChecker) SetWriter(wr Writer) {
	chk.wr = wr
}

// Delimiters overwrites default reader line delimiters
func (chk *SingleChecker) Delimiters(baseDelim, diffDelim byte) {
	chk.baseDelim = baseDelim
	chk.diffDelim = diffDelim
}

// Run begins diff checking reader lines
func (chk *SingleChecker) Run() (equal bool, err error) {
	return true, nil
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
