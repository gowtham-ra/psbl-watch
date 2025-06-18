package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
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
	log.Println("Starting PSBL Game Watcher (scheduled every 3 minutes)...")

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

	// Wait for SIGTERM or SIGINT and shutdown gracefully.
	shutdownGracefully(c)
}

// watchOnce fetches the Hoops-on-Demand page, parses it and sends a push notification.
func watchOnce() {
	// Fetch the Hoops-on-Demand data.
	data, err := fetch.HoopsOnDemandData()
	if err != nil {
		log.Printf("failed to fetch: %v", err)
		return
	}

	// Parse the Game Status from the fetched HTML data.
	gameStatus, err := parse.GameStatusData(data, targetGame)
	if err != nil {
		log.Printf("failed to parse: %v", err)
		return
	}

	// Only send a push notification if the game status has changed.
	prev := store.GetGameStatus()
	if prev != nil && !gameStatusChanged(prev, gameStatus) {
		log.Println("No change in game status - skipping push notification.")
		return
	}

	// Format the push notification message.
	message := fmt.Sprintf("%s\n"+
		"%s\n"+
		"%s\n"+
		"%s\n"+
		"%s\n",
		gameStatus.Target.Gym,
		gameStatus.Target.Type,
		gameStatus.Target.Level,
		gameStatus.Target.DateTime.In(time.FixedZone("PST", -7*3600)).Format("Monday, January 2 at 3:04 PM"),
		formatPlayers(gameStatus.Players))

	// Send a push notification using Pushover.
	notify.SendPushover(message)

	// Update the cache with the latest game status.
	store.SaveGameStatus(gameStatus)
}

func formatPlayers(players map[string][]string) string {
	formatted := ""
	for team, players := range players {
		formatted += fmt.Sprintf("%s: %d\n", team, len(players))
		formatted += fmt.Sprintf("%s\n", strings.Join(players, ", "))
	}
	return formatted
}

func shutdownGracefully(c *cron.Cron) {
	// Set up a channel to listen for SIGTERM and SIGINT.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Wait for SIGTERM or SIGINT.
	<-sigChan
	log.Println("Received SIGTERM or SIGINT, stopping scheduler...")

	// Stop the scheduler gracefully.
	ctx := c.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped.")
}

func gameStatusChanged(prev, curr *store.GameStatus) bool {
	return prev.Found != curr.Found ||
		prev.IsFull != curr.IsFull ||
		prev.TotalPlayers != curr.TotalPlayers ||
		!reflect.DeepEqual(prev.Players, curr.Players)
}
