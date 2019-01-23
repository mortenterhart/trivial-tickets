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

// Package api_out implements a web interface for outgoing mails
// to be fetched and verified to be sent
package api_out

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger"
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

// jsonContentType is used as a constant content type for
// json responses
const jsonContentType = "application/json; charset=utf-8"

// SendMail takes a mail event and a specified ticket and constructs
// a new mail which is then saved into its own file. The message of
// the mail is wrapped inside a mail template depending on the event.
func SendMail(mailEvent mail_events.Event, ticket structs.Ticket) {
	newMail := structs.Mail{
		ID:      random.CreateRandomID(10),
		From:    "no-reply@trivial-tickets.com",
		To:      ticket.Customer,
		Subject: fmt.Sprintf("[trivial-tickets] %s", ticket.Subject),
		Message: mail_events.NewMailBody(mailEvent, ticket),
	}

	logger.Infof(`Composing notification mail (id "%s") to '%s' for %s`,
		newMail.ID, newMail.To, mailEvent.String())

	globals.Mails[newMail.ID] = newMail

	logger.Info("Saving new mail as", globals.ServerConfig.Mails+"/"+newMail.ID+".json")
	writeErr := filehandler.WriteMailFile(globals.ServerConfig.Mails, &newMail)
	if writeErr != nil {
		logger.Errorf("unable to send mail to '%s': %v", ticket.Customer, writeErr)
	}
}

// FetchMails is an endpoint to the outgoing mail API and sends all
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
	logger.APIRequest(request)

	if request.Method == "GET" {

		mails := globals.Mails

		jsonResponse, marshalErr := json.MarshalIndent(&mails, "", "    ")
		if marshalErr != nil {
			httptools.StatusCodeError(writer, marshalErr.Error(), http.StatusInternalServerError)
			return
		}

		logger.Infof("%d %s: Delivering %d mail(s) as response to client", http.StatusOK,
			http.StatusText(http.StatusOK), len(mails))
		writer.Header().Set("Content-Type", jsonContentType)
		fmt.Fprintln(writer, string(jsonResponse))
		return
	}

	httptools.JSONError(writer, structs.JSONMap{
		"status":  http.StatusMethodNotAllowed,
		"message": fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
	logger.Errorf("%d %s: request sent with wrong method '%s', expecting 'GET'", http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed), request.Method)
}

const idParameter = "id"

// VerifyMailSent can be used by an external service to verify that a mail was sent.
// It requests an unique mail id and checks if the corresponding mail exists inside
// the cache. If it does, the mail can be safely deleted and the API returns a verified
// JSON object. If the mail does not exist, the API returns an unverified object.
func VerifyMailSent(writer http.ResponseWriter, request *http.Request) {
	logger.APIRequest(request)

	if request.Method == "POST" {

		var jsonProperties structs.JSONMap
		decodeErr := json.NewDecoder(request.Body).Decode(&jsonProperties)
		if decodeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("could not decode request body: %v", decodeErr),
				http.StatusBadRequest)
			return
		}

		if propErr := verifyMailCheckRequiredPropertiesSet(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("missing required property: %v", propErr),
				http.StatusBadRequest)
			return
		}

		if propErr := verifyMailCheckAdditionalProperties(jsonProperties); propErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("too many properties set: %v", propErr),
				http.StatusBadRequest)
			return
		}

		if typeErr := verifyMailCheckPropertyTypes(jsonProperties); typeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("properties have invalid data types: %v", typeErr),
				http.StatusBadRequest)
			return
		}

		mailID := jsonProperties[idParameter].(string)
		if _, mailExists := globals.Mails[mailID]; !mailExists {
			writer.Header().Set("Content-Type", jsonContentType)
			httptools.JSONResponse(writer, structs.JSONMap{
				"verified": false,
				"status":   http.StatusOK,
				"message":  fmt.Sprintf("mail '%s' does not exist or has already been deleted", mailID),
			})
			logger.Infof("%d %s: Verification of mail '%s' failed: mail does not exist or has already been deleted",
				http.StatusOK, http.StatusText(http.StatusOK), mailID)
			return
		}

		logger.Infof("Removing mail '%s' from global mail storage", mailID)
		delete(globals.Mails, mailID)

		logger.Info("Deleting mail file", globals.ServerConfig.Mails+"/"+mailID+".json")
		if removeErr := filehandler.RemoveMailFile(globals.ServerConfig.Mails, mailID); removeErr != nil {
			httptools.StatusCodeError(writer, fmt.Sprintf("error while trying to remove mail: %v", removeErr),
				http.StatusInternalServerError)
			return
		}

		httptools.JSONResponse(writer, structs.JSONMap{
			"verified": true,
			"status":   http.StatusOK,
			"message":  fmt.Sprintf("mail '%s' was successfully sent and deleted from server cache", mailID),
		})
		logger.Infof("%d %s: Verified sending of mail '%s' successfully and deleted from server cache",
			http.StatusOK, http.StatusText(http.StatusOK), mailID)
		return
	}

	httptools.JSONError(writer, structs.JSONMap{
		"verified": false,
		"status":   http.StatusMethodNotAllowed,
		"message":  fmt.Sprintf("METHOD_NOT_ALLOWED (%s)", request.Method),
	}, http.StatusMethodNotAllowed)
	logger.Errorf("%d %s: request sent with wrong method '%s', expecting 'POST'", http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed), request.Method)
}

// verifyMailCheckRequiredPropertiesSet takes JSON properties and checks if
// the required property "id" is set
func verifyMailCheckRequiredPropertiesSet(jsonProperties structs.JSONMap) error {
	if _, idPropertySet := jsonProperties[idParameter]; !idPropertySet {
		return fmt.Errorf("required JSON property '%s' not defined", idParameter)
	}

	return nil
}

// verifyMailCheckAdditionalProperties checks if any other properties than "id" are
// set within the JSON request
func verifyMailCheckAdditionalProperties(jsonProperties structs.JSONMap) error {
	for property := range jsonProperties {
		if property != idParameter {
			return fmt.Errorf("invalid additional property '%s' defined", property)
		}
	}

	return nil
}

// verifyMailCheckPropertyTypes examines the type of the value of the property "id"
// and verifies its correctness
func verifyMailCheckPropertyTypes(jsonProperties structs.JSONMap) error {
	idProperty := jsonProperties[idParameter]
	if _, idIsString := idProperty.(string); !idIsString {
		return fmt.Errorf("property '%s' has invalid type: expected string, "+
			"instead got %T (located in %s)", idParameter, idProperty, convertToJSON(jsonProperties))
	}

	return nil
}

// convertToJSON converts a json map into a json string and returns it
// as string.
func convertToJSON(properties structs.JSONMap) string {
	jsonString := jsontools.MapToJSON(properties)
	return string(jsonString)
}
