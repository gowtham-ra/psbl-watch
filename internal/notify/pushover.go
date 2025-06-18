package notify

import (
	"log"
	"os"
	"time"

	"github.com/gregdel/pushover"
)

func SendPushover(message string) {
	app := pushover.New(os.Getenv("PUSHOVER_APP_TOKEN"))
	user := pushover.NewRecipient(os.Getenv("PUSHOVER_USER_TOKEN"))

	msg := pushover.NewMessage(message)
	msg.Title = "PSBL Game Available"
	msg.Priority = pushover.PriorityEmergency
	msg.Sound = pushover.SoundBugle
	msg.Retry = 60 * time.Second
	msg.Expire = 300 * time.Second
	msg.URLTitle = "View Game"
	msg.URL = "https://mobile.pugetsoundbasketball.com/"

	_, err := app.SendMessage(msg, user)
	if err != nil {
		log.Printf("Error sending Pushover message: %v", err)
	}
	log.Println("Sent push notification!")
}
