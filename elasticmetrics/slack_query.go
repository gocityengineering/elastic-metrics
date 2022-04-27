package elasticmetrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func slackQuery(team string, markdown string) error {

	// blocks:
	//   - type: section
	//     text:
	// 	     type: mrkdwn
	// 	     text: ...

	var payload = SlackPayload{
		Blocks: []Block{
			{
				Type: "section",
				Text: TextMap{
					Type: "mrkdwn",
					Text: markdown,
				},
			},
		},
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload object")
	}
	body := strings.NewReader(string(bytes))

	// identify Slack API URL
	envVar := team + "SlackApiUrl"
	slackApiUrl := os.Getenv(envVar)
	if slackApiUrl == "" {
		return fmt.Errorf("missing Slack API URL for team %s", team)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(slackApiUrl), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
