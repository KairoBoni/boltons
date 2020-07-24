package main

import "github.com/KairoBoni/boltons/pkg/kafka"

//BodyResponse Arquivei API response model
type BodyResponse struct {
	NFCs []kafka.WorkerMessage `json:"data"`
	Page Page                  `json:"page"`
}

//Page Aquivei API paging structure
type Page struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
