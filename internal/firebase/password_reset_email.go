// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package firebase is common logic and handling around firebase.
package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type sendPasswordResetEmailRequest struct {
	RequestType string `json:"requestType"`
	Email       string `json:"email"`
}

// SendPasswordResetEmail sends a password reset email to the user.
//
// See: https://firebase.google.com/docs/reference/rest/auth#section-send-password-reset-email
func (c *Client) SendPasswordResetEmail(ctx context.Context, email string) error {
	r := &sendPasswordResetEmailRequest{
		RequestType: "PASSWORD_RESET",
		Email:       email,
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(r); err != nil {
		return fmt.Errorf("failed to create json body: %w", err)
	}

	u := c.buildURL("/v1/accounts:sendOobCode")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &body)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response was %d, but failed to read body: %w", status, err)
		}
		return fmt.Errorf("failure %d: %s", status, string(b))
	}

	return nil
}
