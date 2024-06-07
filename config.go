package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

type Response struct {
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Config struct {
	ArubaEndpoint string `yaml:"arubaEndpoint"`
	ArubaUser     []struct {
		ArubaUsername string `yaml:"user"`
		ArubaPassword string `yaml:"password"`
	} `yaml:"arubaUser"`
	ArubaApplicationCredentials []struct {
		ClientID     string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
		CustomerID   string `yaml:"customerId"`
	} `yaml:"arubaApplicationCredentials"`
	ExporterConfig []struct {
		ExporterEndpoint string `yaml:"exporterEndpoint"`
		ExporterPort     string `yaml:"exporterPort"`
	} `yaml:"exporterConfig"`
}

func readConfig(c *Config, configPath string) {
	// Read the YAML file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the YAML data into a Config struct
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatal(err)
	}

}

func authenticate(c *Config, configPath string, response *Response) {

	readConfig(c, configPath)

	clientId := c.ArubaApplicationCredentials[0].ClientID
	clientSecret := c.ArubaApplicationCredentials[1].ClientSecret
	customerId := c.ArubaApplicationCredentials[2].CustomerID

	username := c.ArubaUser[0].ArubaUsername
	password := c.ArubaUser[1].ArubaPassword

	// First request to get the session and CSRF token
	url := c.ArubaEndpoint + "/oauth2/authorize/central/api/login?client_id=" + clientId

	reqBodyValues := map[string]string{
		"username": username,
		"password": password,
	}
	requestBody, _ := json.Marshal(reqBodyValues)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Login failed with status code: %d", resp.StatusCode)
	}

	cookies := resp.Header["Set-Cookie"]
	var csrftoken, session string
	for _, cookie := range cookies {
		if strings.HasPrefix(cookie, "csrftoken") {
			parts := strings.SplitN(cookie, ";", 2)
			if len(parts) > 0 {
				csrftoken = strings.TrimPrefix(parts[0], "csrftoken=")
			}
		}
		if strings.HasPrefix(cookie, "session") {
			parts := strings.SplitN(cookie, ";", 2)
			if len(parts) > 0 {
				session = strings.TrimPrefix(parts[0], "session=")
			}
		}
	}

	// Second request to get the authorization code
	url = c.ArubaEndpoint + "/oauth2/authorize/central/api?client_id=" + clientId + "&response_type=code&scope=all"
	reqBodyValues = map[string]string{"customer_id": customerId}
	requestBody, _ = json.Marshal(reqBodyValues)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "session="+session)
	req.Header.Set("X-CSRF-Token", csrftoken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Authorization request failed with status code: %d", resp.StatusCode)
	}

	// Extract the authorization code from the response
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatalf("Error unmarshaling response body: %v", err)
	}
	code, ok := responseData["auth_code"].(string)
	if !ok {
		log.Fatalf("Authorization code not found in response")
	}

	// Third request to get the access token
	url = c.ArubaEndpoint + "/oauth2/token"
	reqBodyValues = map[string]string{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
	}
	requestBody, _ = json.Marshal(reqBodyValues)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalf("Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Token request failed with status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

}
