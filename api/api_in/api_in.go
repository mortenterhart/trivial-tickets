package api_in

import (
	"encoding/json"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
)

/*
 * Matrikelnummern
 * 3040018
 * 6694964
 * 3478222
 *
 * TODO:
 *           - Currently, JSON body is decoded in ReceiveMail function, this
 *             can be extracted to a central ticket creation / answer creation
 *             function.
 * 			 - Check subject of incoming mails for existing ticket id in the
 *             correct format:
 *
 *                 Subject: [Ticket <ID>] <Ticket subject>
 *
 *             Check for this format and if <ID> already exists as ticket, append
 *             message as new entry to this ticket. Every other case causes a new
 *             ticket to be created.
 *           - Implement noreply@trivial-tickets.de mails:
 *             If a new ticket is created, send mail with permalink to ticket page
 *             and an additional "mailto:<email>?subject=<notification>" (consider
 *             encoding notification subject to URL encoding (Example: <Space> = %20))
 *			   If a new answer is created, send mail with permalink to ticket page
 *             without additional mailto-Link.
 *           - Command-line tool for creating mails to be sent to the server should
 *             define flags -email and -subject to apply email and subject, message
 *             is given through concatenation of command-line arguments or stdin
 *             Example calls:
 *
 *                $ ./send_mail -subject "Hello" -email "example@example.org" "My mail goes here"
 *                $ echo "My mail goes here" | ./send_mail -subject "Hello" -email "example@example.org" (optional)
 *
 *             The tool has to care about converting the parameters into valid JSON with "email",
 *             "subject" and "message" properties and then make an API call to the mail api
 *             using a POST request with the JSON in its body.
 *           - Save sent mails to be requested by an external service
 *
 * NOTE: 1. How shall we test handler methods / following api methods with
 *          http.ResponseWriter and http.Request as parameters? They needed
 *          to be created in the tests artificially (problem of self-created
 *          requests in tests).
 *       2. How shall we simulate an email address for creating tickets or answering
 * 			to tickets? Emails need to be forwarded to the mail api and make
 *          an API call with mail contents in the correct JSON format.
 */

func ReceiveMail(writer http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		// curl -X POST -H 'Content-Type: application/json' --insecure -d '{"email": "example@example.org", "subject": "Test", "message": "Another test"}' https://127.0.0.1:443/api/create_ticket

		var newTicket structs.Mail
		err := json.NewDecoder(req.Body).Decode(&newTicket)

		defer req.Body.Close()

		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		createdTicket := ticket.CreateTicket(newTicket.Email, newTicket.Subject, newTicket.Message)

		globals.Tickets[createdTicket.Id] = createdTicket
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &createdTicket)
	}
}
