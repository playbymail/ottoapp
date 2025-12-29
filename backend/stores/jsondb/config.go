// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package jsondb

import (
	"encoding/json"
	"os"
	"strings"
)

type OttoAppConfig struct {
	Mailgun struct {
		Domain  string `json:"domain"`
		From    string `json:"from"`
		ApiBase string `json:"api-base"`
		ApiKey  string `json:"api-key"`
	} `json:"mailgun"`
}

func LoadOttoAppConfig(path string) (*OttoAppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var oac OttoAppConfig
	err = json.Unmarshal(data, &oac)
	if err != nil {
		return nil, err
	}
	oac.Mailgun.Domain = strings.ToLower(oac.Mailgun.Domain)
	oac.Mailgun.From = strings.ToLower(oac.Mailgun.From)
	return &oac, nil
}
