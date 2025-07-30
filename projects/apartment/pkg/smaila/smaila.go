package smaila

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type Sender struct {
	endpoint string
}

func NewSender(endpoint string) *Sender {
	return &Sender{endpoint: endpoint}
}

func (s *Sender) Send(to []string, subject string, body []byte, html bool) error {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add form fields
	_ = writer.WriteField("to", strings.Join(to, ","))
	_ = writer.WriteField("subject", subject)
	_ = writer.WriteField("body", string(body))
	_ = writer.WriteField("is_html", fmt.Sprintf("%v", html))

	// Close the writer to finalize the body
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the POST request
	req, err := http.NewRequest("POST", s.endpoint, &b)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received error status from mail server: %s", resp.Status)
	}

	return nil
}
