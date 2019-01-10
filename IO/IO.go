package IO

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "regexp"
    "strconv"

    "github.com/mortenterhart/trivial-tickets/structs"
)

var reader = io.Reader(os.Stdin)
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

func PrintEmail(mail structs.Mail) error {
    _, err := fmt.Fprintf(writer, "Email Addresse: %s\n\n"+
        "Betreff: %s\n\n"+
        "%s", mail.Email, mail.Subject, mail.Message)
    return err
}

func GetEmailAddress() (result string, err error) {
    bufReader := bufio.NewReader(reader)
    input, err := bufReader.ReadString('\n')
    // gets rid of the delimiter if there was no error
    if err == nil {
        input = input[:(len(input) - 1)]
    } else if err == io.EOF {
        err = nil
    } else {
        return
    }
    r, _ := regexp.Compile("^(\\w*|\\.*)+@(\\w*)(\\.\\w+)+$")
    result = r.FindString(input)
    if result == "" {
        err = errors.New("not a valid email address")
    }
    return
}

func GetString() (result string, err error) {
    bufReader := bufio.NewReader(reader)
    result, err = bufReader.ReadString('\n')
    // gets rid of the delimiter if there was no error
    if err == nil {
        result = result[:(len(result) - 1)]
    } else if err == io.EOF {
        err = nil
    } else {
        return
    }
    if result == "" {
        err = errors.New("string empty")
    }
    return
}
