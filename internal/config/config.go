package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PhpVMSBaseURL string
	PhpVMSAPIKey  string

	UDPBindHost string
	UDPBindPort int

	TUIEnabled bool

	SelectedAirlineID  int
	SelectedAircraftID int

	SimbriefUserID string

	LogLevel string
}

func DefaultConfig() *Config {
	return &Config{
		PhpVMSBaseURL:      "",
		PhpVMSAPIKey:       "",
		UDPBindHost:        "0.0.0.0",
		UDPBindPort:        47777,
		TUIEnabled:         true,
		SelectedAirlineID:  0,
		SelectedAircraftID: 0,
		SimbriefUserID:     "",
		LogLevel:           "info",
	}
}

func (c *Config) LoadFromDotEnv(filePath string) error {
	if filePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
		filePath = filepath.Join(cwd, ".env")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := godotenv.Load(filePath); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return c.LoadFromEnv()
}

func (c *Config) LoadFromEnv() error {
	if val := os.Getenv("PHPVMS_BASE_URL"); val != "" {
		c.PhpVMSBaseURL = val
	}

	if val := os.Getenv("PHPVMS_API_KEY"); val != "" {
		c.PhpVMSAPIKey = val
	}

	if val := os.Getenv("UDP_BIND_HOST"); val != "" {
		c.UDPBindHost = val
	}

	if val := os.Getenv("UDP_BIND_PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid UDP_BIND_PORT: %w", err)
		}
		c.UDPBindPort = port
	}

	if val := os.Getenv("TUI_ENABLED"); val != "" {
		enabled, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid TUI_ENABLED: %w", err)
		}
		c.TUIEnabled = enabled
	}

	if val := os.Getenv("SIMBRIEF_USER_ID"); val != "" {
		c.SimbriefUserID = val
	}

	if val := os.Getenv("SELECTED_AIRLINE_ID"); val != "" {
		id, err := strconv.Atoi(val)
		if err == nil {
			c.SelectedAirlineID = id
		}
	}

	if val := os.Getenv("SELECTED_AIRCRAFT_ID"); val != "" {
		id, err := strconv.Atoi(val)
		if err == nil {
			c.SelectedAircraftID = id
		}
	}

	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.LogLevel = strings.ToLower(val)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.PhpVMSBaseURL == "" {
		return fmt.Errorf("PHPVMS_BASE_URL is required")
	}

	if c.PhpVMSAPIKey == "" {
		return fmt.Errorf("PHPVMS_API_KEY is required")
	}

	if c.UDPBindPort <= 0 || c.UDPBindPort > 65535 {
		return fmt.Errorf("UDP_BIND_PORT must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	return nil
}

func (c *Config) SavePreferences(filePath string) error {
	if filePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		filePath = filepath.Join(homeDir, ".phpvms-xplane-prefs")
	}

	prefs := map[string]string{
		"SELECTED_AIRLINE_ID":  strconv.Itoa(c.SelectedAirlineID),
		"SELECTED_AIRCRAFT_ID": strconv.Itoa(c.SelectedAircraftID),
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create preferences file: %w", err)
	}
	defer file.Close()

	for key, value := range prefs {
		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return fmt.Errorf("failed to write preference to file: %w", err)
		}
	}

	return nil
}

func (c *Config) LoadPreferences(filePath string) error {
	if filePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		filePath = filepath.Join(homeDir, ".phpvms-xplane-prefs")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read preferences file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "SELECTED_AIRLINE_ID":
			if id, err := strconv.Atoi(value); err == nil {
				c.SelectedAirlineID = id
			}
		case "SELECTED_AIRCRAFT_ID":
			if id, err := strconv.Atoi(value); err == nil {
				c.SelectedAircraftID = id
			}
		}
	}

	return nil
}
