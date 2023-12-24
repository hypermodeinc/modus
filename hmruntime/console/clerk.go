package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClerkAPI struct {
	UserID              string
	BearerToken         string
	ClerkFrontendAPIURL string
	JWT                 string
}

type SignInResponse struct {
	Token string `json:"token"`
}

type ClientSession struct {
	LastActiveToken LastActiveToken `json:"last_active_token"`
}

type LastActiveToken struct {
	JWT string `json:"jwt"`
}

func (c *ClerkAPI) Login() error {
	// First POST request to get the token
	tokenURL := "https://api.clerk.com/v1/sign_in_tokens"
	tokenPayload := map[string]interface{}{
		"user_id":            c.UserID,
		"expires_in_seconds": 2592000,
	}
	jsonData, _ := json.Marshal(tokenPayload)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var signInResponse SignInResponse
	json.Unmarshal(body, &signInResponse)
	token := signInResponse.Token

	// Second POST request to get the JWT
	signInURL := fmt.Sprintf("%s/v1/client/sign_ins?_is_native=true", c.ClerkFrontendAPIURL)
	payload := fmt.Sprintf("strategy=ticket&ticket=%s", token)

	req, err = http.NewRequest("POST", signInURL, bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var session ClientSession
	json.Unmarshal(body, &session)
	c.JWT = session.LastActiveToken.JWT

	return nil
}
