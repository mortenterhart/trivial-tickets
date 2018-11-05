package logic

import (
	"github.com/mortenterhart/trivial-tickets/IO"
	"github.com/mortenterhart/trivial-tickets/structs"
	"log"
	"net"
)

var Output = IO.OutputStringToCommandLine
var NextCommand = IO.GetNextCommand
var requestCom = requestCommand

func MainLoop() {
	ok := true
	for ok {
		com := requestCommand()
		var err error
		switch com {
		case structs.FETCH:
			err = fetchEmails()
		case structs.SUBMIT:
			err = submitEmail()
		case structs.EXIT:
			ok = false
		default:
			log.Fatal(com)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func requestCommand() (com structs.Command) {
	for ok := true; ok; {
		var err error
		Output(string(structs.REQUEST_COMMAND_INPUT))
		com, err = NextCommand()
		ok = err != nil
		if ok {
			Output(structs.COMMAND_NOT_ACCEPTED + err.Error())
		}
	}
	return
}

func submitEmail() error {
	var mail structs.Mail
	var err error
	Output(structs.REQUEST_EMAIL_ADDRESS)
	mail.Email, err = IO.GetEmailAddress()
	for err != nil {
		Output(structs.COMMAND_NOT_ACCEPTED + err.Error() + "\n" + structs.REQUEST_EMAIL_ADDRESS)
		mail.Email, err = IO.GetEmailAddress()
	}
	Output(structs.REQUEST_SUBJECT)
	mail.Subject, err = IO.GetString()
	for err != nil {
		Output(structs.COMMAND_NOT_ACCEPTED + err.Error() + "\n" + structs.REQUEST_SUBJECT)
		mail.Subject, err = IO.GetString()
	}
	Output(structs.REQUEST_MESSAGE)
	mail.Message, err = IO.GetString()
	for err != nil {
		Output(structs.COMMAND_NOT_ACCEPTED + err.Error() + "\n" + structs.REQUEST_MESSAGE)
		mail.Subject, err = IO.GetString()
	}
	err = send(mail, net.ParseIP("127.0.0.0"), 443)
	return err
}

func fetchEmails() error {
	mails, err := receive(net.ParseIP("127.0.0.0"), 443)
	for _, m := range mails {
		IO.PrintEmail(m)
	}
	return err
}

// send turns the email into a JSON string and sends it via POST to the specified address/api/create_ticket and port
func send(email structs.Mail, address net.IP, port uint16) error {
	// just a placeholder so far.
	return nil
}

func receive(address net.IP, port uint16) ([]structs.Mail, error) {
	//just a placeholder so far.
	return make([]structs.Mail, 0), nil
}
