// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"log"
	"net/http"
)

/*
[
  {
    "id": "123",
    "name": "Turn 14",
    "type": "report",
    "status": "parsed",
    "owned": true,
    "shared": false
  }
]

*/

func (s *Server) handleGetDocuments() http.Handler {
	type document_t struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Owned  bool   `json:"owned"`
		Shared bool   `json:"shared"`
	}
	type response []document_t
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var payload response
			// Check if the "me" query parameter is set
			if r.URL.Query().Get("me") == "true" {
				payload = response{
					{
						Id:     "123",
						Name:   "Turn 14",
						Type:   "report",
						Status: "parsed",
						Owned:  true,
						Shared: false,
					},
				}
			}

			if payload == nil {
				payload = []document_t{}
			}
			err := encode(w, r, http.StatusOK, payload)
			if err != nil {
				log.Printf("%s %s: encode %v\n", r.Method, r.URL.Path, err)
			}
		})
}
