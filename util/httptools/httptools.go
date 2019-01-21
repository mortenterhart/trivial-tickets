// Useful tools for HTTP handlers
package httptools

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/logger"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
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
 * Package httptools
 * Useful tools for HTTP handlers
 */

// StatusCodeError writes an error message given by the cause parameter alongside
// with a HTTP status code and its description to the given response writer. The
// status code should be an erroneous one as the function implies it.
func StatusCodeError(writer http.ResponseWriter, cause string, statusCode int) {
	errorMessage := fmt.Sprintf("%d %s: %s", statusCode, http.StatusText(statusCode), cause)
	http.Error(writer, errorMessage, statusCode)
	logger.Error(errorMessage)
}

// JsonResponse
func JsonResponse(writer http.ResponseWriter, jsonProperties structs.JsonMap) {
	writer.Write(append(responseToJson(writer, jsonProperties), '\n'))
}

// JsonError
func JsonError(writer http.ResponseWriter, jsonProperties structs.JsonMap, statusCode int) {
	http.Error(writer, string(responseToJson(writer, jsonProperties)), statusCode)
}

func responseToJson(writer http.ResponseWriter, jsonProperties structs.JsonMap) []byte {
	jsonResponse, conversionErr := jsontools.MapToJson(jsonProperties)
	if conversionErr != nil {
		StatusCodeError(writer, fmt.Sprintf("error building JSON response: %s", conversionErr),
			http.StatusInternalServerError)
		return []byte{}
	}

	return jsonResponse
}
