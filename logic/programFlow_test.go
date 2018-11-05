package logic

import (
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
	"testing"
)

var output []string
var com structs.Command
var err error

func testOutputFunction(out string) {
	output = append(output, out)
}

func clearOutput() {
	output = make([]string, 0)
}

func getNextCommandReplacement() (structs.Command, error) {
	return com, err
}

func TestRequestCommandOutput(t *testing.T) {
	Output = testOutputFunction
	NextCommand = getNextCommandReplacement
	com = structs.FETCH
	err = nil
	requestCommand()
	assert.Equal(t, string(structs.REQUEST_COMMAND_INPUT), output[0])
	clearOutput()
}
