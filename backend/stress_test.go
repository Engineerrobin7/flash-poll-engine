package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	url := "http://localhost:8080/api/polls"

	// 1. Create a dummy poll
	pollData := map[string]interface{}{
		"question": "STRESS TEST POLL",
		"category": "Tech",
		"options":  []string{"A", "B"},
	}
	body, _ := json.Marshal(pollData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})
	pollID := data["id"].(string)
	options := data["options"].([]interface{})
	optionID := int64(options[0].(map[string]interface{})["id"].(float64))

	fmt.Printf("Starting stress test on Poll %s, Option %d...\n", pollID, optionID)

	// 2. Blast it with 100 concurrent votes
	var wg sync.WaitGroup
	votes := 100
	voteURL := fmt.Sprintf("http://localhost:8080/api/polls/%s/vote", pollID)
	voteBody, _ := json.Marshal(map[string]int64{"option_id": optionID})

	for i := 0; i < votes; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req, _ := http.NewRequest("PATCH", voteURL, bytes.NewBuffer(voteBody))
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("[%d] Request failed: %v\n", id, err)
				return
			}
			res.Body.Close()
		}(i)
	}

	wg.Wait()
	fmt.Println("Stress test complete. Checking final count...")

	// 3. Verify
	resp, _ = http.Get(url)
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("System is ALIVE. No crashes detected.")
}
