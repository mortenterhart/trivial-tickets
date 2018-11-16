package communicationToServer

import "github.com/mortenterhart/trivial-tickets/structs"

var serverConfig structs.CLIConfig
var send = sendPost

func FetchEmails() (mails []structs.Mail, err error) {
	return
}

func SubmitEmail(mail structs.Mail) (err error) {
	return
}

func sendPost(payload string, path string) (response string, err error) {
	return
}

func SetServerConfig(config structs.CLIConfig) {

}
