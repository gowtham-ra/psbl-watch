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

func GameStatusData(data *fetch.Result, targetGame store.TargetGame) (*store.GameStatus, error) {
	status := store.GameStatus{
		Target:     targetGame,
		Found:      false,
		IsFull:     false,
		ObservedAt: time.Now(),
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data.Body))
	if err != nil {
		return nil, err
	}

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
		t, err := getGameTime(timeStr, targetGame)
		if err != nil {
			return true
		}

		if !t.Equal(targetGame.DateTime) {
			return true
		}

		// If we are here, we have found the game
		log.Println("Found the game 🏀")
		status.Found = true

		status.Players = getPlayers(s)
		status.IsFull = isGameFull(s, status.Players)
		
		if status.IsFull {
			log.Println("Game is full, no spots available 😢")
		}

		return false // break
	})

	return &status, nil
}

// isGameFull checks if the game is full by checking if the game is full or if both teams are full
func isGameFull(s *goquery.Selection, players map[string][]string) bool {
	gameFull := false
	if s.Find("a:contains('FULL')").Length() > 0 {
		gameFull = true
	}

	var fullTeams int
	for _, team := range players {
		if len(team) >= 7 {
			fullTeams++
		}
	}

	if fullTeams == 2 {
		gameFull = true
	}

	return gameFull
}

// getGameTime parses the game time from the time string and returns the time in the target game's timezone
func getGameTime(timeStr string, targetGame store.TargetGame) (time.Time, error) {
	t, err := time.Parse("Mon, January 2 3:04pm", timeStr)
	if err != nil {
		log.Printf("Error parsing time %s: %v", timeStr, err)
		return time.Now(), err
	}
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t = time.Date(targetGame.DateTime.Year(), t.Month(), t.Day(),
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