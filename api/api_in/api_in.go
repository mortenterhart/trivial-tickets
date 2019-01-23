// Trivial Tickets Ticketsystem
// Copyright (C) 2019 The Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package api_in implements a web interface for incoming mails
// to create new tickets or answers
package api_in

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/httptools"
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

// answerSubjectRegex is a regular expression
// defining the syntax of a subject that creates
// a new answer to an existing ticket.
var answerSubjectRegex = regexp.MustCompile(`\[Ticket "([A-Za-z0-9]+)"\].*`)

// emailRegex defines the syntax of valid email
// addresses.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$")

// stringType is the type for string in order to
// check the request parameter's type for validity.
var stringType = reflect.TypeOf("")

// parameterMap is used as type for the expected
// API parameters and their types to check for
// existence of all required and no additional
// given properties
type parameterMap map[string]reflect.Type

// contains checks if the parameter map contains
// a given key. It returns true if the key is found
// otherwise false.
func (m parameterMap) contains(key string) bool {
	_, found := m[key]
	return found
}

// apiParameters defines the required parameters and their
// types for the handler ReceiveMail. Parameter names are
// mapped to their expected type and are used for parameter
// existence and type checking.
var apiParameters = parameterMap{
	"from":    stringType,
	"subject": stringType,
	"message": stringType,
}

// ReceiveMail serves as the uniform interface for creating new
// tickets and answers out of mails. The mail is passed as JSON
// to this handler and requires the exact properties "from" (the
// sender's email address), "subject" (the ticket subject) and
// "message" (the ticket's message body).
func ReceiveMail(writer http.ResponseWriter, request *http.Request) {
	log.APIRequest(request)

	// Only accept POST requests
	if request.Method == "POST" {

		// Read the request body
		body, readErr := ioutil.ReadAll(request.Body)
		if readErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("unable to read request body: %v", readErr),
				http.StatusInternalServerError)
			return
		}

		// Decode JSON message and save it in jsonProperties map
		// for further investigation
		var jsonProperties structs.JSONMap
		if parseErr := json.Unmarshal(body, &jsonProperties); parseErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("unable to parse JSON due to invalid syntax: %v", parseErr),
				http.StatusBadRequest)
			return
		}

		// Check if all JSON properties required by the API are set
		if propErr := checkRequiredPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("missing required properties in JSON: %v", propErr.Error()),
				http.StatusBadRequest)
			return
		}

		// Check if no additional JSON properties are defined
		if propErr := checkAdditionalPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("too many JSON properties given: %v", propErr),
				http.StatusBadRequest)
			return
		}

		// If all required properties are given, check further if
		// the properties are of the correct data types
		if typeErr := checkCorrectPropertyTypes(jsonProperties); typeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("properties have invalid data types: %v", typeErr),
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
		if ticketID, matchesAnswerRegex := matchAnswerSubject(mail.Subject); matchesAnswerRegex {

			// If so lookup the subject's ticket id in the ticket storage
			// and check if this ticket exists
			if existingTicket, ticketExists := globals.Tickets[ticketID]; ticketExists {
				isAnswerMail = true

				// If the ticket status was already closed, open it again
				if existingTicket.Status == structs.StatusClosed {
					existingTicket.Status = structs.StatusOpen
					log.Infof(`Reopened ticket '%s' (subject "%s") because it was closed`,
						existingTicket.ID, existingTicket.Subject)
				}

				// Update the ticket with a new comment consisting of the
				// email address and message from the mail
				log.Infof(`Attaching new answer from '%s' to ticket '%s' (subject "%s")`,
					mail.From, existingTicket.ID, existingTicket.Subject)
				createdTicket = ticket.UpdateTicket(convertStatusToString(existingTicket.Status),
					mail.From, mail.Message, "extern", existingTicket)

				// Send mail notification to customer that a new answer
				// has been created
				api_out.SendMail(mail_events.NewAnswer, createdTicket)
			} else {
				// The subject is formatted like an answering mail, but the
				// ticket id does not exist
				log.Warnf("Ticket id '%s' does not belong to an existing ticket, creating "+
					"new ticket out of mail", ticketID)
			}
		}

		// If the mail is not an answer create a new ticket in
		// every other case
		if !isAnswerMail {
			createdTicket = ticket.CreateTicket(mail.From, mail.Subject, mail.Message)
			log.Infof(`Creating new ticket "%s" (id '%s') out of mail from '%s'`,
				createdTicket.Subject, createdTicket.ID, mail.From)

			// Send mail notification to customer that a new ticket
			// has been created
			api_out.SendMail(mail_events.NewTicket, createdTicket)
		}

		// Push the created or updated ticket to the ticket storage
		// and write it into its own file
		globals.Tickets[createdTicket.ID] = createdTicket
		if writeErr := filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &createdTicket); writeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("failed to write file for ticket '%s'", createdTicket.ID),
				http.StatusInternalServerError)
			return
		}

		// Construct a JSON response with successful status and message
		// and write it into the response writer
		httptools.JSONResponse(writer, structs.JSONMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		})
		log.Infof("%d %s: Mail request was processed successfully", http.StatusOK, http.StatusText(http.StatusOK))
		return
	}

	// The handler does not accept any other method than POST
	httptools.JSONError(writer, structs.JSONMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
	log.Errorf("%d %s: request sent with wrong method '%s', expecting 'POST'", http.StatusMethodNotAllowed,
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
		ticketIDMatches := answerSubjectRegex.FindStringSubmatch(subject)
		ticketID := ticketIDMatches[1]
		return ticketID, true
	}

	return "", false
}

// validEmailAddress checks the given email address against the
// email regular expression and examines if the supplied email
// address is valid or not.
func validEmailAddress(email string) bool {
	return emailRegex.Match([]byte(email))
}

// propertyNotDefinedError denotes a missing required property
// in the JSON request. It creates an error message telling
// which property name is not defined.
type propertyNotDefinedError struct {
	// propertyName is the name of the missing property.
	propertyName string
}

// Error returns a standard error message for a missing required
// property and the corresponding property name.
func (err propertyNotDefinedError) Error() string {
	return fmt.Sprintf("required JSON property not defined: '%s'", err.propertyName)
}

// newPropertyNotDefinedError creates a new error object with
// the property name that is missing in case one is missing.
func newPropertyNotDefinedError(propertyName string) propertyNotDefinedError {
	return propertyNotDefinedError{propertyName}
}

// checkRequiredPropertiesSet checks if the properties sent
// within the request contain all required property names that
// the API expects. If all required properties are defined
// the result is nil, otherwise an error is returned.
func checkRequiredPropertiesSet(jsonProperties structs.JSONMap) error {
	for requiredProperty := range apiParameters {
		if propErr := checkPropertySet(jsonProperties, requiredProperty); propErr != nil {
			return errors.Wrap(propErr, "missing properties in JSON body")
		}
	}

	return nil
}

// checkPropertySet is a helper function of checkRequiredPropertiesSet.
// It checks if a single property name is defined in the json properties
// map. In case the property is defined it returns a nil error, otherwise
// a new propertyNotDefinedError with the property name that is missing
// is returned.
func checkPropertySet(jsonProperties structs.JSONMap, propName string) error {
	if _, defined := jsonProperties[propName]; defined {
		return nil
	}

	return newPropertyNotDefinedError(propName)
}

// checkAdditionalPropertiesSet checks if any other than the required
// properties are defined in the json properties map. If there are
// additional properties, an error with the name of that property is
// returned, otherwise nil.
func checkAdditionalPropertiesSet(jsonProperties structs.JSONMap) error {
	for key := range jsonProperties {
		if !apiParameters.contains(key) {
			return fmt.Errorf("JSON contains illegal additional property: '%s'", key)
		}
	}

	return nil
}

// checkCorrectPropertyTypes reports whether all given properties have
// the correct data type. The property values are type checked using
// reflections and compared to the defined types in the apiParameters
// map. If there is a type mismatch, an exhaustive error with expected
// and given type and location is returned, otherwise nil.
func checkCorrectPropertyTypes(jsonProperties structs.JSONMap) error {
	for parameter, parameterType := range apiParameters {
		if property, propertyGiven := jsonProperties[parameter]; reflect.TypeOf(property) != parameterType {
			if !propertyGiven {
				return newPropertyNotDefinedError(parameter)
			}

			return fmt.Errorf("type mismatch in property '%s': expected %s, instead got %T "+
				"(located in %s)",
				parameter, parameterType.Name(), property, writeJSONProperty(parameter, property))
		}
	}

	return nil
}

// writeJSONProperty writes a key-value-pair in correct JSON format
// and returns it as string.
func writeJSONProperty(key string, value interface{}) string {
	return fmt.Sprint(enquote(key), ":", writeJSONValue(value))
}

// writeJSONValue writes a value in correct JSON format and returns
// it as a string.
func writeJSONValue(value interface{}) string {
	if stringValue, isString := value.(string); isString {
		return enquote(stringValue)
	}

	return fmt.Sprintf("%v", value)
}

// enquote surrounds a given potion with double quotes to be used
// as JSON key or string value.
func enquote(potion interface{}) string {
	return fmt.Sprintf(`"%v"`, potion)
}
