package console

import (
	"hmruntime/utils"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

const (
	CLERK_SIGNIN_URL = "https://api.clerk.com/v1/sign_in_tokens"
)

var modelSpecCache = map[string]ModelSpec{}

type SignInToken struct {
	Token string `json:"token"`
}

type JWTToken struct {
	Client struct {
		Sessions []struct {
			LastActiveToken struct {
				JWT string `json:"jwt"`
			} `json:"last_active_token"`
		} `json:"sessions"`
	} `json:"client"`
}

func GetConsoleJWT() (string, error) {
	payload := map[string]interface{}{
		"user_id":            os.Getenv("CLERK_USER_ID"),
		"expires_in_seconds": 2592000,
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv("CLERK_API_KEY"),
	}

	log.Info().Msgf("Payload: %v", payload)
	log.Info().Msgf("Headers: %v", headers)

	tokenResponse, err := utils.PostHttp[SignInToken](CLERK_SIGNIN_URL, payload, headers)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("Token Response: %v", tokenResponse)

	log.Info().Msgf("Token: %s", tokenResponse.Token)

	token := tokenResponse.Token

	// Define form data
	data := url.Values{}
	data.Set("strategy", "ticket")
	data.Set("ticket", token)

	headers = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	jwtResponse, err := utils.PostHttp[JWTToken](os.Getenv("CLERK_JWT_URL")+"/v1/client/sign_ins?_is_native=true", data.Encode(), headers)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("JWT Response: %v", jwtResponse)

	log.Info().Msgf("JWT: %s", jwtResponse.Client.Sessions[0].LastActiveToken.JWT)

	return jwtResponse.Client.Sessions[0].LastActiveToken.JWT, nil
}

type ModelSpec struct {
	ID        string `json:"id"`
	ModelType string `json:"modelType"`
	Endpoint  string `json:"endpoint"`
}

type ModelSpecPayload struct {
	Data struct {
		ModelSpec ModelSpec `json:"modelSpec"`
	} `json:"data"`
}

func GetModelSpec(modelId string) (ModelSpec, error) {
	// Check if the model spec is already in the cache
	if spec, ok := modelSpecCache[modelId]; ok {
		return spec, nil
	}

	jwt, err := GetConsoleJWT()
	if err != nil {
		return ModelSpec{}, err
	}

	serviceURL := os.Getenv("CONSOLE_URL") + "/graphql"

	query := `
		query GetModelSpec($id: ID!) {
			modelSpec: getModelSpec(id: $id) {
				id
				modelType
				endpoint
			}
		}`

	request := map[string]any{
		"query":     query,
		"variables": map[string]any{"id": modelId},
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": jwt,
	}

	response, err := utils.PostHttp[ModelSpecPayload](serviceURL, request, headers)
	if err != nil {
		return ModelSpec{}, err
	}

	spec := response.Data.ModelSpec
	if spec.ID != modelId {
		return ModelSpec{}, err
	}

	// Cache the model spec
	modelSpecCache[modelId] = spec

	return response.Data.ModelSpec, nil
}
