package elasticmetrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func emailServiceQuery(recipients []string, subject, markdown, emailService string, emailServicePort int) error {

	// recipients: [ "mail1@domain", "mail2@domain" ]
	// subject: "somesubject"
	// markdownBody: "somemarkdown"

	var payload = EmailServicePayload{
		Recipients:   recipients,
		Subject:      subject,
		MarkdownBody: markdown,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload object")
	}
	body := strings.NewReader(string(bytes))

	if emailService == "" {
		return fmt.Errorf("missing email service URL")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/", emailService, emailServicePort), body)
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
