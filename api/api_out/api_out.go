// Web API for outgoing mails to be fetched and verified to be sent
package api_out

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/httptools"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"github.com/mortenterhart/trivial-tickets/util/random"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 6694964, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package api_out
 * Web API for outgoing mails to be fetched and verified to be sent
 */

// SendMail takes a mail event and a specified ticket and constructs
// a new mail which is then saved into its own file. The message of
// the mail is wrapped inside a mail template depending on the event.
func SendMail(mailEvent mail_events.Event, ticket structs.Ticket) {
	newMail := structs.Mail{
		Id:      random.CreateRandomId(10),
		From:    "no-reply@trivial-tickets.com",
		To:      ticket.Customer,
		Subject: fmt.Sprintf("[trivial-tickets] %s", ticket.Subject),
		Message: mail_events.NewMailBody(mailEvent, ticket),
	}

	globals.Mails[newMail.Id] = newMail

	writeErr := filehandler.WriteMailFile(globals.ServerConfig.Mails, &newMail)
	if writeErr != nil {
		log.Printf("unable to send mail to '%s': %s\n", ticket.Customer, writeErr)
	}
}

// FetchMail is an endpoint to the outgoing mail API and sends all
// mails which are currently cached and ready to be sent. The response
// is in JSON format.
//
//     Takes: no parameters
//     Returns: {
//         "<mail_id>": {
//             "from": "",
//             "id": "",
//             "message": "",
//             "subject": "",
//             "to": ""
//         }
//     }
func FetchMails(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {

		mails := globals.Mails

		jsonResponse, marshalErr := json.MarshalIndent(&mails, "", "    ")
		if marshalErr != nil {
			httptools.StatusCodeError(writer, marshalErr.Error(), http.StatusInternalServerError)
			return
		}

		writer.Write(append(jsonResponse, '\n'))
		return
	}

	httptools.JsonError(writer, structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
}

// VerifyMailSent can be used by an external service to verify that a mail was sent.
// It requests an unique mail id and checks if the corresponding mail exists inside
// the cache. If it does, the mail can be safely deleted and the API returns a verified
// JSON object. If the mail does not exist, the API returns an unverified object.
func VerifyMailSent(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {

		var jsonProperties structs.JsonMap
		decodeErr := json.NewDecoder(request.Body).Decode(&jsonProperties)
		if decodeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("could not decode request body: %s", decodeErr),
				http.StatusBadRequest)
			return
		}

		if propErr := verifyMailCheckRequiredPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("missing required property: %s", propErr),
				http.StatusBadRequest)
			return
		}

		if propErr := verifyMailCheckAdditionalProperties(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("too many properties set: %s", propErr),
				http.StatusBadRequest)
			return
		}

		if typeErr := verifyMailCheckPropertyTypes(jsonProperties); typeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("properties have invalid data types: %s", typeErr),
				http.StatusBadRequest)
			return
		}

		mailId := jsonProperties["id"].(string)
		if _, mailExists := globals.Mails[mailId]; !mailExists {
			httptools.JsonResponse(writer, structs.JsonMap{
				"verified": false,
				"message":  fmt.Sprintf("mail '%s' does not exist or has already been deleted", mailId),
			})
			return
		}

		delete(globals.Mails, mailId)

		if removeErr := filehandler.RemoveMailFile(globals.ServerConfig.Mails, mailId); removeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("error while trying to remove mail: %s", removeErr),
				http.StatusInternalServerError)
			return
		}

		httptools.JsonResponse(writer, structs.JsonMap{
			"verified": true,
			"message":  fmt.Sprintf("mail '%s' was successfully sent and deleted from server cache", mailId),
		})
		return
	}

	httptools.JsonError(writer, structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
}

// verifyMailCheckRequiredPropertiesSet takes JSON properties and checks if
// the required property "id" is set
func verifyMailCheckRequiredPropertiesSet(jsonProperties structs.JsonMap) error {
	if _, idPropertySet := jsonProperties["id"]; !idPropertySet {
		return fmt.Errorf("required JSON property '%s' not defined", "id")
	}

	return nil
}

// verifyMailCheckAdditionalProperties checks if any other properties than "id" are
// set within the JSON request
func verifyMailCheckAdditionalProperties(jsonProperties structs.JsonMap) error {
	for property := range jsonProperties {
		if property != "id" {
			return fmt.Errorf("invalid additional property '%s' defined", property)
		}
	}

	return nil
}

// verifyMailCheckPropertyTypes examines the type of the value of the property "id"
// and verifies its correctness
func verifyMailCheckPropertyTypes(jsonProperties structs.JsonMap) error {
	if idContent, idIsString := jsonProperties["id"].(string); !idIsString {
		return fmt.Errorf("property '%s' has invalid type: expected string, "+
			"instead got %T (located in %s)", "id", idContent, convertToJson(jsonProperties))
	}

	return nil
}

// convertToJson converts a json map into a json string and logs an error if it failed
func convertToJson(properties structs.JsonMap) string {
	jsonString, decodeErr := jsontools.MapToJson(properties)
	if decodeErr != nil {
		log.Println(decodeErr)
		return ""
	}

	return string(jsonString)
}
