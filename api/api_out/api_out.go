package api_out

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/mortenterhart/trivial-tickets/globals"
    "github.com/mortenterhart/trivial-tickets/structs"
    "github.com/mortenterhart/trivial-tickets/util/filehandler"
    "github.com/mortenterhart/trivial-tickets/util/httptools"
    "github.com/mortenterhart/trivial-tickets/util/random"
)

// Construct a mail and save it to cache
func SendMail(email, subject, message string) {
    mail := structs.Mail{
        Id:      random.CreateRandomId(10),
        Email:   email,
        Subject: subject,
        Message: message,
    }

    writeErr := filehandler.WriteMailFile(globals.ServerConfig.Mails, &mail)
    if writeErr != nil {
        log.Printf("unable to send mail to '%s': %s\n", email, writeErr)
    }
}

// Output all cached mails
func FetchMails(writer http.ResponseWriter, request *http.Request) {

    if request.Method == "GET" {
        mails, readErr := filehandler.ReadMailFiles(globals.ServerConfig.Mails)
        if readErr != nil {
            httptools.StatusCodeError(writer, readErr.Error(), http.StatusInternalServerError)
            return
        }

        jsonResponse, marshalErr := json.MarshalIndent(mails, "", "    ")
        if marshalErr != nil {
            httptools.StatusCodeError(writer, marshalErr.Error(), http.StatusInternalServerError)
            return
        }

        writer.Write(append(jsonResponse, '\n'))
        return
    }

    httptools.JsonError(writer, map[string]interface{}{
        "status": http.StatusMethodNotAllowed,
        "message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
    }, http.StatusMethodNotAllowed)
}

// VerifyMailSent gets a mail id and verifies the mail with the id is cached
// and exists, then deletes it because the sending is verified, otherwise sending
// is retried on next call to FetchMails
func VerifyMailSent(writer http.ResponseWriter, request *http.Request) {

    if request.Method == "POST" {
        return
    }

    // http.Error()
}
