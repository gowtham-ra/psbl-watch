package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/gowtham-ra/psbl-watch/internal/fetch"
	"github.com/gowtham-ra/psbl-watch/internal/notify"
	"github.com/gowtham-ra/psbl-watch/internal/parse"
	"github.com/gowtham-ra/psbl-watch/internal/store"
)

// targetGame represents the game we are interested in tracking.
var targetGame = store.TargetGame{
	Gym:      "Seattle Central College #1",
	Type:     "Saturday Morning Hoops",
	Level:    "Recreational-CoEd",
	DateTime: time.Date(2025, 6, 21, 10, 0, 0, 0, time.Local),
}

func main() {
	log.Println("Starting psbl-watcher (scheduled every 3 minutes)...")

	// Create a new cron scheduler in the local time-zone.
	c := cron.New(cron.WithLocation(time.Local))

	// Schedule the watcher to run every 3 minutes.
	_, err := c.AddFunc("@every 3m", watchOnce)
	if err != nil {
		log.Fatalf("failed to schedule cron job: %v", err)
	}

	// Start the scheduler.
	c.Start()

	// Watch once immediately, so that we do not have to wait 3 minutes.
	watchOnce()

	// Block forever for the application does not exit.
	select {}
}

// watchOnce fetches the Hoops-on-Demand page, parses it and sends a push notification.
func watchOnce() {
	data, err := fetch.HoopsOnDemandData()
	if err != nil {
		log.Printf("failed to fetch: %v", err)
		return
	}

	gameStatus, err := parse.GameStatusData(data, targetGame)
	if err != nil {
		log.Printf("failed to parse: %v", err)
		return
	}

	message := fmt.Sprintf("Observed At: %v\n"+
		"Gym: %s\n"+
		"Type: %s\n"+
		"DateTime: %v\n"+
		"Players:\n%s",
		gameStatus.ObservedAt.Format(time.RFC3339),
		gameStatus.Target.Gym,
		gameStatus.Target.Type,
		gameStatus.Target.DateTime.Format(time.RFC3339),
		formatPlayers(gameStatus.Players))

	notify.SendPushover(message)
}

func formatPlayers(players map[string][]string) string {
	formatted := ""
	for team, players := range players {
		formatted += fmt.Sprintf("%s: %s\n\n", team, strings.Join(players, ", "))
	}
	return formatted
}
