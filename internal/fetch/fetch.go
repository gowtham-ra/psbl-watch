package fetch

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

const (
	maxRetries     = 10
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
	URL            = "https://mobile.pugetsoundbasketball.com/registration-hod.php"
)

type Result struct {
	Body []byte
}

func HoopsOnDemandData() (*Result, error) {
	var lastErr error
	backoff := initialBackoff

	for range maxRetries {
		log.Printf("Fetching Hoops-on-Demand data from %s", URL)
		resp, err := http.Get(URL)
		if err != nil {
			log.Printf("error fetching: %v, retrying...", err)
			lastErr = err
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected status code: %d, retrying...", resp.StatusCode)
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
			continue
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading body: %v, retrying...", err)
			lastErr = err
			time.Sleep(backoff)
			backoff = time.Duration(math.Min(float64(backoff*2), float64(maxBackoff)))
			continue
		}

		return &Result{
			Body: body,
		}, nil
	}

	return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}
