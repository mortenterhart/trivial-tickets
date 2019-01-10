package httptools

import (
    "fmt"
    "net/http"

    "github.com/mortenterhart/trivial-tickets/util/jsontools"
)

func StatusCodeError(writer http.ResponseWriter, cause string, statusCode int) {
    http.Error(writer, fmt.Sprintf("%d %s: %s", statusCode, http.StatusText(statusCode), cause), statusCode)
}

// successful response 200 OK with appended newline
func JsonResponse(writer http.ResponseWriter, jsonProperties map[string]interface{}) {
    writer.Write(append(responseToJson(writer, jsonProperties), '\n'))
}

// erroneous response
func JsonError(writer http.ResponseWriter, jsonProperties map[string]interface{}, statusCode int) {
    http.Error(writer, string(responseToJson(writer, jsonProperties)), statusCode)
}

func responseToJson(writer http.ResponseWriter, jsonProperties map[string]interface{}) []byte {
    jsonResponse, conversionErr := jsontools.MapToJson(jsonProperties)
    if conversionErr != nil {
        StatusCodeError(writer, fmt.Sprintf("error building JSON response: %s", conversionErr),
            http.StatusInternalServerError)
        return []byte{}
    }

    return jsonResponse
}
