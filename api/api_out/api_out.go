package api_out

import (
    "net/http"
)

// Construct a mail and save it to cache
func SendMail(email, subject, message string) {
    /*mail := structs.Mail{
        Id:      random.CreateRandomId(10),
        Email:   email,
        Subject: subject,
        Message: message,
    }*/

}

// Output all cached mails
func FetchMails(writer http.ResponseWriter, request *http.Request) {

}

// VerifyMailSent gets a mail id and verifies the mail with the id is cached
// and exists, then deletes it because the sending is verified, otherwise sending
// is retried on next call to FetchMails
func VerifyMailSent(writer http.ResponseWriter, request *http.Request) {

}
