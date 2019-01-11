package api_out

import (
	"encoding/json"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"log"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/httptools"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"github.com/mortenterhart/trivial-tickets/util/random"
)

// Construct a mail_events and save it to cache
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
		log.Printf("unable to send newMail to '%s': %s\n", ticket.Customer, writeErr)
	}
}

// Output all cached mails
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

// Input: JSON {"id":"<mail_id>"}
// Output: JSON {"sent":true/false,"message":"mail_events id does not exist"}
// VerifyMailSent gets a mail_events id and verifies the mail_events with the id is cached
// and exists, then deletes it because the sending is verified, otherwise sending
// is retried on next call to FetchMails
func VerifyMailSent(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {

		var jsonProperties structs.JsonMap
		decodeErr := json.NewDecoder(request.Body).Decode(&jsonProperties)
		if decodeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("could not decode request body: %s", decodeErr),
				http.StatusInternalServerError)
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
				"message":  fmt.Sprintf("mail_events '%s' does not exist or has already been deleted", mailId),
			})
			return
		}

		delete(globals.Mails, mailId)

		if removeErr := filehandler.RemoveMailFile(mailId); removeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("error while trying to remove mail_events: %s", removeErr),
				http.StatusInternalServerError)
			return
		}

		httptools.JsonResponse(writer, structs.JsonMap{
			"verified": true,
			"message":  fmt.Sprintf("mail_events '%s' was successfully sent and deleted from server cache", mailId),
		})
		return
	}

	httptools.JsonError(writer, structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
}

func verifyMailCheckRequiredPropertiesSet(jsonProperties structs.JsonMap) error {
	if _, idPropertySet := jsonProperties["id"]; !idPropertySet {
		return fmt.Errorf("required JSON property '%s' not defined", "id")
	}

	return nil
}

func verifyMailCheckAdditionalProperties(jsonProperties structs.JsonMap) error {
	for property := range jsonProperties {
		if property != "id" {
			return fmt.Errorf("invalid additional property '%s' defined", property)
		}
	}

	return nil
}

func verifyMailCheckPropertyTypes(jsonProperties structs.JsonMap) error {
	if idContent, idIsString := jsonProperties["id"].(string); !idIsString {
		return fmt.Errorf("property '%s' has invalid type: expected string, "+
			"instead got %T (located in %s)", "id", idContent, convertToJson(jsonProperties))
	}

	return nil
}

func convertToJson(properties structs.JsonMap) string {
	jsonString, decodeErr := jsontools.MapToJson(properties)
	if decodeErr != nil {
		log.Println(decodeErr)
		return ""
	}

	return string(jsonString)
}
