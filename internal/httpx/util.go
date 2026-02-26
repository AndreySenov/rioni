// Copyright (c) 2026 Andrey Senov
// SPDX-License-Identifier: Apache-2.0

package httpx

import "net/http"

const (
	ContentTypeApplicationDnsMessage = "application/dns-message"
	HeaderAccept                     = "Accept"
	HeaderAllow                      = "Allow"
	HeaderContentType                = "Content-Type"
)

func IsStatusOk(r *http.Response) bool {
	return r != nil && r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices
}

func IsStatusNotOk(r *http.Response) bool {
	return !IsStatusOk(r)
}
