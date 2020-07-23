package main

type BodyResponse struct {
	NFCs []NFC `json:"data"`
	Page Page  `json:"page"`
}

type NFC struct {
	AccessKey string `json:"access_key"`
	XML       string `json:"xml"`
}

type Page struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
