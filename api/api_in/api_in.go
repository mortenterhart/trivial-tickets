// Web API for incoming mails to create new tickets or answers
package api_in

import (
	"encoding/json"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/logger"
	"io/ioutil"
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
 * Matriculation numbers: 3040018, 6694964, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package api_in
 * Web API for incoming mails to create new tickets or answers
 */

// Regex defining the syntax of an answer subject
var answerSubjectRegex = regexp.MustCompile(`\[Ticket "([A-Za-z0-9]+)"\].*`)

// Regex defining the syntax of valid email addresses
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")

// Type constant for string in order to check the
// request parameter's type for validity
var stringType = reflect.TypeOf("")

// API parameters for the handler ReceiveMail.
// Parameter names are mapped to their expected type and are
// used for parameter existence and type checking.
var apiParameters = map[string]reflect.Type{
	"from":    stringType,
	"subject": stringType,
	"message": stringType,
}

// ReceiveMail serves as the uniform interface for creating new tickets and answers
// out of mails. The mail is passed as JSON to this handler and requires the exact
// properties "from" (the sender's email address), "subject" (the ticket subject)
// and "message" (the ticket's message body).
func ReceiveMail(writer http.ResponseWriter, request *http.Request) {
	logger.ApiRequest(request)

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
					logger.Infof(`Reopened ticket '%s' (subject "%s") because it was closed`,
						existingTicket.Id, existingTicket.Subject)
				}

				// Update the ticket with a new comment consisting of the
				// email address and message from the mail
				logger.Infof(`Attaching new answer from '%s' to ticket '%s' (subject "%s")`+"\n",
					mail.From, existingTicket.Id, existingTicket.Subject)
				createdTicket = ticket.UpdateTicket(convertStatusToString(existingTicket.Status),
					mail.From, mail.Message, "extern", existingTicket)

				// Send mail notification to customer that a new answer
				// has been created
				api_out.SendMail(mail_events.NewAnswer, createdTicket)
			} else {
				// The subject is formatted like an answering mail, but the
				// ticket id does not exist
				logger.Warnf("Ticket id '%s' does not belong to an existing ticket, creating "+
					"new ticket out of mail\n", ticketId)
			}
		}

		// If the mail is not an answer create a new ticket in every other case
		if !isAnswerMail {
			createdTicket = ticket.CreateTicket(mail.From, mail.Subject, mail.Message)
			logger.Infof(`Creating new ticket "%s" (id '%s') out of mail from '%s'`,
				createdTicket.Subject, createdTicket.Id, mail.From)

			// Send mail notification to customer that a new ticket
			// has been created
			api_out.SendMail(mail_events.NewTicket, createdTicket)
		}

		// Push the created or updated ticket to the ticket storage and write
		// it into its own file
		globals.Tickets[createdTicket.Id] = createdTicket
		if writeErr := filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &createdTicket); writeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("failed to write file for ticket '%s'", createdTicket.Id),
				http.StatusInternalServerError)
			return
		}

		// Construct a JSON response with successful status and message
		// and write it into the response writer
		httptools.JsonResponse(writer, structs.JsonMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		})
		logger.Infof("%d %s: Mail request was processed successfully", http.StatusOK, http.StatusText(http.StatusOK))
		return
	}

	// The handler does not accept any other method than POST
	httptools.JsonError(writer, structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
	logger.Infof("%d %s: request sent with wrong method '%s', expecting 'POST'", http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed), request.Method)
}

// convertStatusToString converts a status enum constant which is
// an integer to a string considering correct string conversion.
// Casting the integer to a string is not an option since the string
// will consist of the character at the Unicode index that the integer
// infers.
func convertStatusToString(status structs.Status) string {
	return strconv.Itoa(int(status))
}

// matchAnswerSubject matches the given subject against the syntax
// of a subject which causes a new answer to be created instead of
// a new ticket. If the subject is conform to this pattern, the
// function returns the contained ticket id as string and true,
// otherwise an empty string and false.
func matchAnswerSubject(subject string) (string, bool) {
	if answerSubjectRegex.Match([]byte(subject)) {
		ticketIdMatches := answerSubjectRegex.FindStringSubmatch(subject)
		ticketId := ticketIdMatches[1]
		return ticketId, true
	}

	return "", false
}

// validEmailAddress checks the given email address against the email
// regular expression and examines if the supplied email address is
// valid or not.
func validEmailAddress(email string) bool {
	return emailRegex.Match([]byte(email))
}

// Error denoting a missing property in the JSON request.
// propertyName holds the name of the missing property.
type propertyNotDefinedError struct {
	propertyName string
}

// Error returns a standard error message for missing required properties
// and the property name.
func (err propertyNotDefinedError) Error() string {
	return fmt.Sprintf("required JSON property not defined: '%s'", err.propertyName)
}

// newPropertyNotDefinedError creates a new error object with the property
// name that is missing in case one is missing.
func newPropertyNotDefinedError(propertyName string) propertyNotDefinedError {
	return propertyNotDefinedError{propertyName}
}

// checkRequiredPropertiesSet checks if the properties sent within the request contain
// all required property names that the API expects. If all required properties are
// defined, the result is nil, otherwise an error is returned.
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

// checkPropertySet is a helper function of checkRequiredPropertiesSet. It checks if
// a single property name is defined in the json properties map. The result is true,
// if propName is defined in the map props. If it is not defined, a panic will be thrown
// with a propertyNotDefinedError and the corresponding property name. The panic is
// recovered in the parent function.
func checkPropertySet(props structs.JsonMap, propName string) bool {
	if _, defined := props[propName]; defined {
		return true
	}

	panic(newPropertyNotDefinedError(propName))
}

// checkAdditionalPropertiesSet checks if any other than the required properties are defined
// in the json properties map. If there are additional properties, an error with the name
// of that property is returned, otherwise nil.
func checkAdditionalPropertiesSet(jsonProperties structs.JsonMap) error {
	permittedKeys := newStringList("from", "subject", "message")
	for key := range jsonProperties {
		if !permittedKeys.contains(key) {
			return fmt.Errorf("JSON contains illegal additional property: '%s'", key)
		}
	}

	return nil
}

// checkCorrectPropertyTypes reports whether all given properties have the correct data type.
// The property values are type checked using reflections and compared to the defined types
// in the apiParameters map. If there is a type mismatch, an exhaustive error with expected
// and given type and location is returned, otherwise nil.
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

// writeJsonProperty writes a key-value-pair in correct JSON format
// and returns it as string
func writeJsonProperty(key, value interface{}) string {
	jsonKey := enquote(key) + ":"
	return jsonKey + writeJsonValue(value)
}

// writeJsonValue writes a value in correct JSON format and returns it
// as a string
func writeJsonValue(value interface{}) string {
	if stringValue, isString := value.(string); isString {
		return enquote(stringValue)
	}

	return fmt.Sprintf("%v", value)
}

// enquote surrounds a given potion with double quotes
func enquote(potion interface{}) string {
	return fmt.Sprintf(`"%v"`, potion)
}

// stringList is a type for a list of strings
type stringList []string

// newStringList returns a new list of strings with
// the given strings as initial values
func newStringList(values ...string) stringList {
	return stringList(values)
}

// contains searches after a value inside a string list
// and returns true if it found, otherwise false.
func (slice stringList) contains(value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}

	return false
}
