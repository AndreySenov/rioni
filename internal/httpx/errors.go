// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package httpx

import "net/http"

func ErrorBadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func ErrorMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func ErrorNotAcceptable(w http.ResponseWriter) {
	http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
}

func ErrorLengthRequired(w http.ResponseWriter) {
	http.Error(w, "Length Required", http.StatusLengthRequired)
}

func ErrorRequestEntityTooLarge(w http.ResponseWriter) {
	http.Error(w, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
}

func ErrorUnsupportedMediaType(w http.ResponseWriter) {
	http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
}

func ErrorInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ErrorBadGateway(w http.ResponseWriter) {
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}
