package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type SwitchResponse struct {
	Count    int      `json:"count"`
	Switches []Switch `json:"switches"`
}

type Switch struct {
	ClientCount      int      `json:"client_count"`
	CPUUtilization   int      `json:"cpu_utilization"`
	FanSpeed         string   `json:"fan_speed"`
	FirmwareVersion  string   `json:"firmware_version"`
	GroupID          int      `json:"group_id"`
	GroupName        string   `json:"group_name"`
	IPAddress        string   `json:"ip_address"`
	LabelIDs         []int    `json:"label_ids"`
	Labels           []string `json:"labels"`
	MacAddress       string   `json:"macaddr"`
	MaxPower         int      `json:"max_power"`
	MemFree          int      `json:"mem_free"`
	MemTotal         int      `json:"mem_total"`
	Model            string   `json:"model"`
	Name             string   `json:"name"`
	PoeConsumption   string   `json:"poe_consumption"`
	PowerConsumption int      `json:"power_consumption"`
	PublicIPAddress  string   `json:"public_ip_address"`
	Serial           string   `json:"serial"`
	Site             string   `json:"site"`
	SiteID           int      `json:"site_id"`
	StackID          string   `json:"stack_id"`
	StackMemberID    int      `json:"stack_member_id"`
	Status           string   `json:"status"`
	SwitchRole       int      `json:"switch_role"`
	SwitchType       string   `json:"switch_type"`
	Temperature      string   `json:"temperature"`
	UplinkPorts      []struct {
		Port string `json:"port"`
	} `json:"uplink_ports"`
	Uptime int `json:"uptime"`
	Usage  int `json:"usage"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

var (
	switchCpuUtilization = prometheus.NewDesc("switch_cpu_utilization", "CPU Utilization of the switch in percentge", []string{"name"}, nil)
	switchUsage          = prometheus.NewDesc("switch_usage", "Bandwidth usage of the switch in mb/s", []string{"name"}, nil)
	switchUptime         = prometheus.NewDesc("switch_uptime", "Uptime of the switch in seconds", []string{"name"}, nil)
)

type Exporter struct {
	arubaEndpoint, arubaAccessToken, arubaRefreshToken string
}

func NewExporter(arubaEndpoint string, arubaAccessToken string, arubaRefreshToken string) *Exporter {
	return &Exporter{
		arubaEndpoint:     arubaEndpoint,
		arubaAccessToken:  arubaAccessToken,
		arubaRefreshToken: arubaRefreshToken,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- switchCpuUtilization
	ch <- switchUsage
	ch <- switchUptime
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	refreshToken(e)
	listSwitches(e, ch)

}

func main() {

	arubaEndpoint := "https://apigw-eucentral3.central.arubanetworks.com/"

	// Read tokens from YAML file
	tokensFile, err := ioutil.ReadFile("tokens.yaml")

	if err != nil {
		fmt.Println("Error reading tokens from YAML file:", err)
		return
	}

	tokens := make(map[string]string)
	err = yaml.Unmarshal(tokensFile, &tokens)
	if err != nil {
		fmt.Println("Error unmarshaling tokens from YAML:", err)
		return
	}

	arubaAccessToken := tokens["arubaAccessToken"]
	arubaRefreshToken := tokens["arubaRefreshToken"]

	exporter := NewExporter(arubaEndpoint, arubaAccessToken, arubaRefreshToken)
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)

}

func listSwitches(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/switches?show_resource_details=true&calculate_client_count=true"

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	// Parse JSON
	var switchResponse SwitchResponse
	if err := json.Unmarshal(body, &switchResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, s := range switchResponse.Switches {

		ch <- prometheus.MustNewConstMetric(switchCpuUtilization, prometheus.GaugeValue, float64(s.CPUUtilization), s.Name+"_"+s.MacAddress)
		ch <- prometheus.MustNewConstMetric(switchUsage, prometheus.GaugeValue, float64(s.Usage), s.Name+"_"+s.MacAddress)
		ch <- prometheus.MustNewConstMetric(switchUptime, prometheus.GaugeValue, float64(s.Uptime), s.Name+"_"+s.MacAddress)
	}
}

func refreshToken(e *Exporter) {

	// Read tokens from YAML file
	tokensFile, err := ioutil.ReadFile("client.yaml")

	if err != nil {
		fmt.Println("Error reading client data from YAML file:", err)
		return
	}

	tokens := make(map[string]string)
	err = yaml.Unmarshal(tokensFile, &tokens)
	if err != nil {
		fmt.Println("Error unmarshaling client data from YAML:", err)
		return
	}

	clientId := tokens["clientId"]
	clientSecret := tokens["clientSecret"]

	url := e.arubaEndpoint + "oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=refresh_token&refresh_token=" + e.arubaRefreshToken

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

	}

	req.Header.Set("Authorization", "Bearer "+e.arubaAccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

	}

	// Parse JSON

	var tokenResponse TokenResponse

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)

	}

	e.arubaAccessToken = tokenResponse.AccessToken
	e.arubaRefreshToken = tokenResponse.RefreshToken

	// Write tokens to YAML file
	if e.arubaAccessToken != "" && e.arubaRefreshToken != "" {
		tokens := map[string]string{
			"arubaAccessToken":  e.arubaAccessToken,
			"arubaRefreshToken": e.arubaRefreshToken,
		}

		data, err := yaml.Marshal(tokens)
		if err != nil {
			fmt.Println("Error marshaling tokens to YAML:", err)
			return
		}

		err = ioutil.WriteFile("tokens.yaml", data, 0644)
		if err != nil {
			fmt.Println("Error writing tokens to YAML file:", err)
			return
		}
	}

}
