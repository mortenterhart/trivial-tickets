package IO

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/structs"
	"io"
	"os"
	"strconv"
)

var reader io.Reader = io.Reader(os.Stdin)
var writer = io.Writer(os.Stdout)

func GetNextCommand() (structs.Command, error) {
	bufReader := bufio.NewReader(reader)
	var ret structs.Command
	var asInt int
	input, err := bufReader.ReadString('\n')
	// gets rid of the delimiter if there was no error
	if err == nil {
		input = input[:(len(input) - 1)]
	}
	asInt, err = strconv.Atoi(input)
	switch structs.Command(asInt) {
	case structs.FETCH,
		structs.SUBMIT,
		structs.EXIT:
		ret = structs.Command(asInt)
	default:
		err = errors.New("not within range of valid options")

	}
	return ret, err
}

func OutputStringToCommandLine(output string) {
	fmt.Fprintf(writer, "%s", output)
}
