package IO

import (
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type TestWriter struct {
	output *string
}

func NewTestWriter(out *string) (w *TestWriter) {
	return &TestWriter{out}
}
func (w *TestWriter) Write(p []byte) (n int, err error) {
	*w.output = string(p)
	return len(p), nil
}

type ITestReader interface {
	io.Reader
	setData(d string)
}

type TestReader struct {
	data []byte
}

func NewTestReader(data string) (r *TestReader) {
	return &TestReader{[]byte(data)}
}
func (r *TestReader) setData(d string) {
	r.data = []byte(d)
}
func (r *TestReader) readByte() byte {
	// this function assumes that eof() check was done before
	b := r.data[0]
	r.data = r.data[1:]
	return b
}
func (r *TestReader) eof() (eof bool) {
	return len(r.data) == 0
}
func (r *TestReader) Read(p []byte) (n int, err error) {
	if r.eof() {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			p[n] = r.readByte()
			n++
			if r.eof() {
				break
			}
		}
	}
	return
}

func TestGetNextCommandSuccess(t *testing.T) {
	reader = NewTestReader("0")
	command, _ := GetNextCommand()
	assert.Equal(t, structs.FETCH, command, "does not return the correct number defined in structs.FETCH.")
	r, ok := reader.(ITestReader)
	if ok {
		r.setData("1")
		command, _ = GetNextCommand()
		assert.Equal(t, structs.SUBMIT, command, "does not return the correct number defined in structs.SUBMIT.")
		r.setData("2")
		command, _ = GetNextCommand()
		assert.Equal(t, structs.EXIT, command, "does not return the correct number defined in structs.EXIT.")
	}
}

func TestGetNextCommandError(t *testing.T) {
	reader = NewTestReader("-1")
	_, err := GetNextCommand()
	assert.Error(t, err, "-1 should not be a valid argument.")
	r, ok := reader.(ITestReader)
	if ok {
		r.setData("abcd")
		_, err = GetNextCommand()
		assert.Error(t, err, "abcd should not be a valid argument.")
		r.setData("")
		_, err = GetNextCommand()
		assert.Error(t, err, "'' should not be a valid argument.")
	}
}

func TestOutputStringToCommandLine(t *testing.T) {
	var output string
	writer = NewTestWriter(&output)
	testString := "a String!!"
	OutputStringToCommandLine(testString)
	assert.Equal(t, testString, output, "string output failed.")
}
