// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package runners implements API commands
package runners

import (
	"fmt"
	"net"
	"net/url"
)

type Runner struct {
	baseUrl string
}

func New(schema, host, port string) *Runner {
	return &Runner{baseUrl: fmt.Sprintf("%s://%s", schema, net.JoinHostPort(host, port))}
}

func (r *Runner) apiUrl(route string, queryParameters url.Values) string {
	if len(queryParameters) == 0 {
		return r.baseUrl + route
	}
	return r.baseUrl + route + "?" + queryParameters.Encode()
}
