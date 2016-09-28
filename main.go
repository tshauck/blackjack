package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {
	L := log.New()
	a := NewAgent(L)
	a.PlayGames(100000)

	a.SavePolicy("policy.json")
}
