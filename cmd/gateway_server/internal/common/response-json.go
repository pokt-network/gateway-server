// Package common provides common utilities and structures used across the application.

//go:generate ffjson $GOFILE
package common

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
)

// ErrorResponse represents a JSON-formatted error response.
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// JSONError creates a JSON-formatted error response and sends it to the client.
// It takes the fasthttp.RequestCtx, error message, and HTTP status code as parameters.
func JSONError(ctx *fasthttp.RequestCtx, message string, statusCode int) {
	// Create an ErrorResponse instance with the provided message and status code.
	errorResponse := ErrorResponse{
		Message: message,
		Status:  statusCode,
	}

	// Marshal the ErrorResponse instance into JSON format.
	jsonData, err := ffjson.Marshal(errorResponse)
	if err != nil {
		// If there's an error during JSON marshaling, log it and set a generic internal server error response.
		ctx.Error(fmt.Sprintf("Error marshaling JSON: %s", err), fasthttp.StatusInternalServerError)
		return
	}

	// Set the response headers and body with the JSON data.
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(jsonData)
	// Set the HTTP status code for the response.
	ctx.SetStatusCode(statusCode)
}

// JSONError creates a JSON-formatted error response and sends it to the client.
// It takes the fasthttp.RequestCtx, error message, and HTTP status code as parameters.
func JSONSuccess(ctx *fasthttp.RequestCtx, data any, statusCode int) {

	// Marshal the ErrorResponse instance into JSON format.
	jsonData, err := ffjson.Marshal(data)
	if err != nil {
		// If there's an error during JSON marshaling, log it and set a generic internal server error response.
		ctx.Error(fmt.Sprintf("Error marshaling JSON: %s", err), fasthttp.StatusInternalServerError)
		return
	}

	// Set the response headers and body with the JSON data.
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(jsonData)
	// Set the HTTP status code for the response.
	ctx.SetStatusCode(statusCode)
}
