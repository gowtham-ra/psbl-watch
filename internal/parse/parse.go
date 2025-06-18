package parse

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gowtham-ra/psbl-watch/internal/fetch"
	"github.com/gowtham-ra/psbl-watch/internal/store"
)

func GameStatusData(data *fetch.Result, targetGame store.TargetGame) ([]*store.GameStatus, error) {
	var statuses []*store.GameStatus

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data.Body))
	if err != nil {
		return nil, err
	}

	log.Printf("Total games: %d", doc.Find("div.mobilehod").Length())
	doc.Find("div.mobilehod").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		// Get the game info details from the header
		gameInfoHeader := strings.TrimSpace(s.Find("#gameinfo div").First().Text())
		lines := strings.Split(gameInfoHeader, "\n")

		timeStr := strings.TrimSpace(lines[0])
		gym := strings.TrimSpace(lines[1])
		gameType := strings.TrimSpace(lines[2])
		gameLevel := strings.TrimSpace(lines[3])

		// Return true will continue the searching
		if !strings.EqualFold(gameType, targetGame.Type) {
			return true
		}

		if !strings.EqualFold(gym, targetGame.Gym) {
			return true
		}

		if !strings.EqualFold(gameLevel, targetGame.Level) {
			return true
		}

		// Check if the game time is the same as the target game time
		t, err := getGameTime(timeStr)
		if err != nil {
			log.Printf("Error parsing time %s: %v", timeStr, err)
			return true
		}

		if !targetGame.DateTime.IsZero() && !t.Equal(targetGame.DateTime) {
			return true
		}

		// If we are here, we have found the game
		log.Printf("Found Game: %s %s %s %s 🏀", gym, gameType, gameLevel, timeStr)

		status := store.GameStatus{
			Target:     targetGame,
			Found:      true,
			IsFull:     false,
			ObservedAt: time.Now(),
		}
		// If the game time is not set, set it to the game time
		// so that we can use it as key in the cache.
		if targetGame.DateTime.IsZero() {
			status.Target.DateTime = t
		}

		status.Players = getPlayers(s)
		for _, players := range status.Players {
			status.TotalPlayers += len(players)
		}
		status.IsFull = isGameFull(s, &status)

		if status.IsFull {
			log.Println("Game is full 🔴")
		} else {
			log.Println("Game has free spots 🟢")
		}

		status.Target.GameKey = uniqueGameKey(&status)
		statuses = append(statuses, &status)
		return true
	})

	return statuses, nil
}

// isGameFull checks if the game is full or if the total number of players is >= 14
func isGameFull(s *goquery.Selection, gameStatus *store.GameStatus) bool {
	gameFull := false
	if s.Find("a:contains('FULL')").Length() > 0 {
		gameFull = true
	}

	if gameStatus.TotalPlayers >= 14 {
		gameFull = true
	}

	return gameFull
}

// getGameTime parses the game time from the time string and returns the time in the target game's timezone
func getGameTime(timeStr string) (time.Time, error) {
	t, err := time.Parse("Mon, January 2 3:04pm", timeStr)
	if err != nil {
		log.Printf("Error parsing time %s: %v", timeStr, err)
		return time.Now(), err
	}
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t = time.Date(2025, t.Month(), t.Day(),
		t.Hour(), t.Minute(), 0, 0, loc)

	return t, nil
}

// getPlayers parses the players from the roster and returns a map of team name to player names
func getPlayers(s *goquery.Selection) map[string][]string {
	players := make(map[string][]string)

	roster := s.Find("#roster")

	// Find all team divs
	roster.Find("div[style*='float:left;width:50%']").Each(func(_ int, teamDiv *goquery.Selection) {
		teamName := strings.TrimSpace(teamDiv.Find("span.team_name").Text())
		var playerNames []string

		// Find all player names for this team
		teamDiv.Find("span.player_name").Each(func(_ int, player *goquery.Selection) {
			playerName := strings.TrimSpace(player.Text())
			// Remove the number prefix (e.g., "1. ")
			if idx := strings.Index(playerName, ". "); idx != -1 {
				playerName = playerName[idx+2:]
			}
			playerNames = append(playerNames, playerName)
		})

		players[teamName] = playerNames
	})

	return players
}

// uniqueGameKey returns a unique, deterministic key for the given GameStatus so
// that we can cache a GameStatus for every different game we are tracking.
// The key is a combination of the gym, type, level, date/time, and team names.
func uniqueGameKey(status *store.GameStatus) string {
	return status.Target.Gym + 
	"|" + status.Target.Type + 
	"|" + status.Target.Level + 
	"|" + status.Target.DateTime.In(time.UTC).Format(time.RFC3339) +
	"|" + strings.Join(getMapKeys(status.Players), ",") // Team names
}

// getMapKeys returns the keys of a map as a comma-separated string
func getMapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}