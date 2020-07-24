package main

import "github.com/KairoBoni/boltons/pkg/kafka"

type BodyResponse struct {
	NFCs []kafka.WorkerMessage `json:"data"`
	Page Page                  `json:"page"`
}

type Page struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
