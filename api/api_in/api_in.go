package api_in

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/httptools"
	"github.com/pkg/errors"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package api_in
 * Web API for incoming mails to create new tickets or answers
 */

var answerSubjectRegex = regexp.MustCompile(`\[Ticket "([A-Za-z0-9]+)"\].*`)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")

var stringType = reflect.TypeOf("")

var apiParameters = map[string]reflect.Type{
	"from":    stringType,
	"subject": stringType,
	"message": stringType,
}

func ReceiveMail(writer http.ResponseWriter, request *http.Request) {

	// Only accept POST requests
	if request.Method == "POST" {

		// Read the request body
		body, readErr := ioutil.ReadAll(request.Body)
		if readErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("unable to read request body: %s", readErr),
				http.StatusInternalServerError)
			return
		}

		// Decode JSON message and save it in jsonProperties map
		// for further investigation
		var jsonProperties structs.JsonMap
		if parseErr := json.Unmarshal(body, &jsonProperties); parseErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("unable to parse JSON due to invalid syntax: %s", parseErr),
				http.StatusBadRequest)
			return
		}

		// Check if all JSON properties required by the API are set
		if propErr := checkRequiredPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("missing required properties in JSON: %s", propErr.Error()),
				http.StatusBadRequest)
			return
		}

		// Check if no additional JSON properties are defined
		if propErr := checkAdditionalPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("too many JSON properties given: %s", propErr),
				http.StatusBadRequest)
			return
		}

		// If all required properties are given, check further if
		// the properties are of the correct data types
		if typeErr := checkCorrectPropertyTypes(jsonProperties); typeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("properties have invalid data types: %s", typeErr),
				http.StatusBadRequest)
			return
		}

		// Populate the mail struct with the previously parsed JSON properties
		mail := structs.Mail{
			From:    jsonProperties["from"].(string),
			Subject: jsonProperties["subject"].(string),
			Message: jsonProperties["message"].(string),
		}

		// Validate the email address syntax using the above regular expression
		if !validEmailAddress(mail.From) {
			httptools.StatusCodeError(writer, fmt.Sprintf("invalid email address given: '%s'", mail.From),
				http.StatusBadRequest)
			return
		}

		// Container for the created or updated ticket
		var createdTicket structs.Ticket

		// Flag indicating that an incoming request belongs to an answer
		isAnswerMail := false

		// Determine if the email's subject is compliant to the answer
		// regular expression
		if ticketId, matchesAnswerRegex := matchAnswerSubject(mail.Subject); matchesAnswerRegex {

			// If so lookup the subject's ticket id in the ticket storage
			// and check if this ticket exists
			if existingTicket, ticketExists := globals.Tickets[ticketId]; ticketExists {
				isAnswerMail = true

				// If the ticket status was already closed, open it again
				if existingTicket.Status == structs.CLOSED {
					existingTicket.Status = structs.OPEN
				}

				// Update the ticket with a new comment consisting of the
				// email address and message from the mail
				log.Printf(`Attaching new answer from '%s' to ticket '%s' (id "%s")`+"\n",
					mail.From, existingTicket.Subject, existingTicket.Id)
				createdTicket = ticket.UpdateTicket(convertStatusToString(existingTicket.Status),
					mail.From, mail.Message, "extern", existingTicket)

				api_out.SendMail(mail_events.NewAnswer, createdTicket)
			} else {
				// The subject is formatted like an answering mail, but the
				// ticket id does not exist
				log.Printf("WARNING: ticket id '%s' does not belong to an existing ticket, creating "+
					"new ticket out of mail\n", ticketId)
			}
		}

		// If the mail is not an answer create a new ticket in every other case
		if !isAnswerMail {
			createdTicket = ticket.CreateTicket(mail.From, mail.Subject, mail.Message)

			api_out.SendMail(mail_events.NewTicket, createdTicket)
		}

		// Push the created or updated ticket to the ticket storage and write
		// it into its own file
		globals.Tickets[createdTicket.Id] = createdTicket
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &createdTicket)

		// Construct a JSON response with successful status and message
		// and write it into the response writer
		httptools.JsonResponse(writer, structs.JsonMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		})
		return
	}

	// The handler does not accept any other method than POST
	httptools.JsonError(writer, structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
}

func convertStatusToString(status structs.Status) string {
	return strconv.Itoa(int(status))
}

func matchAnswerSubject(subject string) (string, bool) {
	if answerSubjectRegex.Match([]byte(subject)) {
		ticketIdMatches := answerSubjectRegex.FindStringSubmatch(subject)
		ticketId := ticketIdMatches[1]
		return ticketId, true
	}

	return "", false
}

func validEmailAddress(email string) bool {
	return emailRegex.Match([]byte(email))
}

type propertyNotDefinedError struct {
	propertyName string
}

func (err propertyNotDefinedError) Error() string {
	return fmt.Sprintf("required JSON property not defined: '%s'", err.propertyName)
}

func newPropertyNotDefinedError(propertyName string) propertyNotDefinedError {
	return propertyNotDefinedError{propertyName}
}

func checkRequiredPropertiesSet(jsonProperties structs.JsonMap) (returnErr error) {
	defer func() {
		propError := recover()
		if errValue, isPropError := propError.(propertyNotDefinedError); isPropError {
			returnErr = errors.Wrap(errValue, "missing properties in JSON body")
		}
	}()

	propsSet := checkPropertySet(jsonProperties, "from")
	propsSet = propsSet && checkPropertySet(jsonProperties, "subject")
	propsSet = propsSet && checkPropertySet(jsonProperties, "message")

	if propsSet {
		return nil
	}

	return errors.New("missing properties in JSON body")
}

func checkPropertySet(props structs.JsonMap, propName string) bool {
	if _, defined := props[propName]; defined {
		return true
	}

	panic(newPropertyNotDefinedError(propName))
}

func checkAdditionalPropertiesSet(jsonProperties structs.JsonMap) error {
	permittedKeys := newStringList("from", "subject", "message")
	for key := range jsonProperties {
		if !permittedKeys.contains(key) {
			return fmt.Errorf("JSON contains illegal additional property: '%s'", key)
		}
	}

	return nil
}

func checkCorrectPropertyTypes(jsonProperties structs.JsonMap) error {
	for parameter, parameterType := range apiParameters {
		if property, propertyGiven := jsonProperties[parameter]; reflect.TypeOf(property) != parameterType {
			if !propertyGiven {
				return newPropertyNotDefinedError(parameter)
			}

			return fmt.Errorf("type mismatch in property '%s': expected %s, instead got %T "+
				"(located in %s)",
				parameter, parameterType.Name(), property, writeJsonProperty(parameter, property))
		}
	}

	return nil
}

func writeJsonProperty(key, value interface{}) string {
	jsonKey := enquote(key) + ":"
	return jsonKey + writeJsonValue(value)
}

func writeJsonValue(value interface{}) string {
	if stringValue, isString := value.(string); isString {
		return enquote(stringValue)
	}

	return fmt.Sprintf("%v", value)
}

func enquote(potion interface{}) string {
	return fmt.Sprintf(`"%v"`, potion)
}

type stringList []string

func newStringList(values ...string) stringList {
	return stringList(values)
}

func (slice stringList) contains(value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}

	return false
}
