/*
 * Copyright 2023 Hypermode, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the License);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ClerkUsername    string `json:"clerkUsername"`
	ClerkFrontendURL string `json:"clerkFrontendURL"`
	ConsoleURL       string `json:"consoleURL"`
}

func GetAppConfiguration(env string) (*Config, error) {
	var config Config

	const configPath = "config/%s.json"

	path := fmt.Sprintf(configPath, env)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &config, nil
}
