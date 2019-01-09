package api_in

import (
	"encoding/json"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/pkg/errors"
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
 *                 Subject: [Ticket "<ID>"] <Ticket subject>
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

type jsonMap map[string]interface{}

var answerSubjectRegex = regexp.MustCompile(`\[Ticket "([A-Za-z0-9]+)"\]\s*[A-Za-z0-9_\s"'\.,+=\[\]()@/&$§!#-]*`)

var stringType = reflect.TypeOf("")

var apiParameters = map[string]reflect.Type{
	"email":   stringType,
	"subject": stringType,
	"message": stringType,
}

var jsonProperties = make(jsonMap)

func ReceiveMail(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		body, readErr := ioutil.ReadAll(request.Body)
		if readErr != nil {
			http.Error(writer, fmt.Sprintf("unable to read request body: %s", readErr), http.StatusInternalServerError)
			return
		}

		parseErr := json.Unmarshal(body, &jsonProperties)
		if parseErr != nil {
			http.Error(writer, fmt.Sprintf("unable to parse JSON due to invalid syntax: %s", parseErr), http.StatusBadRequest)
			return
		}

		propErr := checkRequiredPropertiesSet()
		if propErr != nil {
			http.Error(writer, fmt.Sprintf("missing required properties in JSON: %s", propErr), http.StatusBadRequest)
			return
		}

		propErr = checkAdditionalPropertiesSet()
		if propErr != nil {
			http.Error(writer, fmt.Sprintf("too many JSON properties given: %s", propErr), http.StatusBadRequest)
			return
		}

		typeErr := checkCorrectPropertyTypes()
		if typeErr != nil {
			http.Error(writer, fmt.Sprintf("properties have invalid data types: %s", typeErr), http.StatusBadRequest)
			return
		}

		// Extract mail from JSON body
		mail := structs.Mail{
			Email:   jsonProperties["email"].(string),
			Subject: jsonProperties["subject"].(string),
			Message: jsonProperties["message"].(string),
		}

		var createdTicket structs.Ticket

		isAnswerMail := false

		if ticketId, matchesAnswerRegex := matchSubject(mail.Subject); matchesAnswerRegex {
			if existingTicket, ticketExists := globals.Tickets[ticketId]; ticketExists {
				isAnswerMail = true

				if existingTicket.Status == structs.CLOSED {
					existingTicket.Status = structs.OPEN
				}

				log.Printf(`"Attaching new answer from '%s' to ticket '%s' (id "%s")`+"\n",
					mail.Email, existingTicket.Subject, existingTicket.Id)
				createdTicket = ticket.UpdateTicket(convertStatusToString(existingTicket.Status),
					mail.Email, mail.Message, "extern", existingTicket)
			} else {
				log.Printf("WARNING: ticket id '%s' does not belong to an existing ticket, creating "+
					"new ticket out of mail\n", ticketId)
			}
		}

		if !isAnswerMail {
			createdTicket = ticket.CreateTicket(mail.Email, mail.Subject, mail.Message)
		}

		globals.Tickets[createdTicket.Id] = createdTicket
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &createdTicket)

		writer.Write([]byte(buildJSONResponseStatus(http.StatusOK, "OK") + "\n"))
		return
	}

	http.Error(writer, buildJSONResponseStatus(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED"), http.StatusMethodNotAllowed)
}

func convertStatusToString(status structs.State) string {
	return strconv.Itoa(int(status))
}

func matchSubject(subject string) (string, bool) {
	if answerSubjectRegex.Match([]byte(subject)) {
		ticketIdMatches := answerSubjectRegex.FindStringSubmatch(subject)
		ticketId := ticketIdMatches[0]
		return ticketId, true
	}

	return "", false
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

func checkRequiredPropertiesSet() (returnErr error) {
	defer func() {
		propError := recover()
		if errValue, isPropError := propError.(propertyNotDefinedError); isPropError {
			returnErr = errors.Wrap(errValue, "missing properties in JSON body")
		}
	}()

	propsSet := checkPropertySet(jsonProperties, "email")
	propsSet = propsSet && checkPropertySet(jsonProperties, "subject")
	propsSet = propsSet && checkPropertySet(jsonProperties, "message")

	if propsSet {
		return nil
	}

	return errors.New("missing properties in JSON body")
}

func checkPropertySet(props jsonMap, propName string) bool {
	if _, defined := props[propName]; defined {
		return true
	}

	panic(newPropertyNotDefinedError(propName))
}

func checkAdditionalPropertiesSet() error {
	permittedKeys := newStringList("email", "subject", "message")
	for key := range jsonProperties {
		if !permittedKeys.contains(key) {
			return fmt.Errorf("JSON contains illegal additional property: '%s'", key)
		}
	}

	return nil
}

func checkCorrectPropertyTypes() error {
	for parameter, parameterType := range apiParameters {
		if property, propertyGiven := jsonProperties[parameter]; reflect.TypeOf(property) != parameterType {
			if !propertyGiven {
				return newPropertyNotDefinedError(parameter)
			}

			return fmt.Errorf("type mismatch in property '%s': expected %s, instead got %T "+
				`(located in %s)`,
				parameter, parameterType.Name(), property, writeJSONProperty(parameter, property))
		}
	}

	return nil
}

func writeJSONProperty(key, value interface{}) string {
	var jsonBuilder strings.Builder
	jsonBuilder.WriteString(enquote(key))
	jsonBuilder.WriteString(":")
	jsonBuilder.WriteString(writeJSONValue(value))

	return jsonBuilder.String()
}

func writeJSONValue(value interface{}) string {
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

func buildJSONResponseStatus(statusCode int, message string) string {
	return fmt.Sprintf(`{"status":%d,"message":"%s"}`, statusCode, message)
}

func checkEmailSyntax(email string) error {
	return nil
}
