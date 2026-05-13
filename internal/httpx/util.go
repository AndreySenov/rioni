// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"errors"
	"mime"
	"net/http"
)

const (
	ContentTypeApplicationDnsMessage = "application/dns-message"
	HeaderAccept                     = "Accept"
	HeaderAllow                      = "Allow"
	HeaderContentType                = "Content-Type"
)

func IsStatusOk(response *http.Response) bool {
	return response != nil && response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
}

func IsStatusNotOk(response *http.Response) bool {
	return !IsStatusOk(response)
}

func GetResponseContentType(response *http.Response) (mediatype string, params map[string]string, err error) {
	if response == nil {
		return "", nil, errors.New("response is nil")
	}
	return mime.ParseMediaType(response.Header.Get(HeaderContentType))
}

func GetRequestAttrs(request *http.Request) []any {
	if request == nil {
		return nil
	}
	return []any{
		"method", request.Method,
		"host", request.Host,
		"remote_addr", request.RemoteAddr,
		"user_agent", request.UserAgent(),
	}
}
