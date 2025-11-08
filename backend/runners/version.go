// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package runners

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (r *Runner) GetVersion(showBuildInfo bool) error {
	var params url.Values
	hasParameters := showBuildInfo == true
	if hasParameters { // Encode the parameters and append them to the route
		params = url.Values{}
		params.Add("show-build-info", "true")
	}
	apiUrl := r.apiUrl("/api/version", params)
	log.Printf("url %s\n", apiUrl)

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return err
	}

	fmt.Println(string(body))
	return nil
}
