package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gowtham-ra/psbl-watch/internal/fetch"
	"github.com/gowtham-ra/psbl-watch/internal/parse"
	"github.com/gowtham-ra/psbl-watch/internal/store"
)

func main() {
	fmt.Println("Starting psbl-watcher...")
	data, err := fetch.HoopsOnDemandData()
	if err != nil {
		log.Fatalf("failed to fetch: %v", err)
	}

	targetGame := store.TargetGame{
		Gym:      "Seattle Central College #1",
		Type:     "Saturday Morning Hoops",
		Level:    "Recreational-CoEd",
		DateTime: time.Date(2025, 6, 21, 10, 0, 0, 0, time.Local),
	}

	gameStatus, err := parse.GameStatusData(data, targetGame)
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}

	fmt.Printf("Game Status:\n")
	fmt.Printf("\tFound: %v\n", gameStatus.Found)
	fmt.Printf("\tIs Full: %v\n", gameStatus.IsFull)
	fmt.Printf("\tObserved At: %v\n", gameStatus.ObservedAt.Format(time.RFC3339))
	fmt.Printf("\tTarget Game:\n")
	fmt.Printf("\tGym: %s\n", gameStatus.Target.Gym)
	fmt.Printf("\tType: %s\n", gameStatus.Target.Type)
	fmt.Printf("\tDateTime: %v\n", gameStatus.Target.DateTime.Format(time.RFC3339))
	fmt.Printf("\tPlayers:\n")
	for team, players := range gameStatus.Players {
		fmt.Printf("\t\t%s:\n", team)
		for _, player := range players {
			fmt.Printf("\t\t\t- %s\n", player)
		}
	}
}
