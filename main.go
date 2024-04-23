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

type ApResponse struct {
	AccessPoints []AccessPoint `json:"aps"`
}

type AccessPoint struct {
	ApDeploymentNode   string   `json:"ap_deployment_mode"`
	ApGroup            string   `json:"ap_group"`
	ClientCount        int      `json:"client_count"`
	ClusterId          string   `json:"cluster_id"`
	ControllerName     string   `json:"controller_name"`
	CpuUtilization     int      `json:"cpu_utilization"`
	FirmwareVersion    string   `json:"firmware_version"`
	GatewayClusterId   string   `json:"gateway_cluster_id"`
	GatewayClusterName string   `json:"gateway_cluster_name"`
	GroupName          string   `json:"group_name"`
	IpAddress          string   `json:"ip_address"`
	Labels             []string `json:"labels"`
	LastModified       int      `json:"last_modified"`
	MacAddress         string   `json:"macaddr"`
	MemFree            int      `json:"mem_free"`
	MemTotal           int      `json:"mem_total"`
	MeshRole           string   `json:"mesh_role"`
	Model              string   `json:"model"`
	Name               string   `json:"name"`
	Notes              string   `json:"notes"`
	PublicIpAddress    string   `json:"public_ip_address"`
	Radios             []struct {
		Band          int    `json:"band"`
		Channel       string `json:"channel"`
		Index         int    `json:"index"`
		MacAddress    string `json:"macaddr"`
		Node          int    `json:"node"`
		RadioName     string `json:"radio_name"`
		RadioType     string `json:"radio_type"`
		SpatialStream string `json:"spatial_stream"`
		Status        string `json:"status"`
		TxPower       int    `json:"tx_power"`
		Utilization   int    `json:"utilization"`
	} `json:"radios"`
	Serial      string `json:"serial"`
	Site        string `json:"site"`
	SleepStatus bool   `json:"sleep_status"`
	Status      string `json:"status"`
	SubnetMask  string `json:"subnet_mask"`
	SwarmId     string `json:"swarm_id"`
	SwarmMaster bool   `json:"swarm_master"`
	SwarmName   string `json:"swarm_name"`
	Uptime      int    `json:"uptime"`
}

type McResponse struct {
	Count               int                  `json:"count"`
	MobilityControllers []MobilityController `json:"mcs"`
}

type MobilityController struct {
	CpuUtilization        int      `json:"cpu_utilization"`
	FirmwareBackupVersion string   `json:"firmware_backup_version"`
	FirmwareVersion       string   `json:"firmware_version"`
	GroupName             string   `json:"group_name"`
	IpAddress             string   `json:"ip_address"`
	Labels                []string `json:"labels"`
	MacRange              string   `json:"mac_range"`
	MacAddress            string   `json:"macaddr"`
	MemFree               int      `json:"mem_free"`
	MemTotal              int      `json:"mem_total"`
	Mode                  string   `json:"mode"`
	Model                 string   `json:"model"`
	Name                  string   `json:"name"`
	RebootReason          string   `json:"reboot_reason"`
	Role                  string   `json:"role"`
	Serial                string   `json:"serial"`
	Site                  string   `json:"site"`
	Status                string   `json:"status"`
	Uptime                int      `json:"uptime"`
}

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
	apClientCount    = prometheus.NewDesc("ap_client_count", "Number of clients connected to access point", []string{"name"}, nil)
	apCpuUtilization = prometheus.NewDesc("ap_cpu_utilization", "CPU Utilization of the access point in percentge", []string{"name"}, nil)
	apMemFree        = prometheus.NewDesc("ap_mem_free", "Amount of free memory of access point", []string{"name"}, nil)
	apMemTotal       = prometheus.NewDesc("ap_mem_total", "Total amount of  memory of access point", []string{"name"}, nil)
	apUptime         = prometheus.NewDesc("ap_uptime", "Uptime of the access point in seconds", []string{"name"}, nil)

	mcCpuUtilization = prometheus.NewDesc("mc_cpu_utilization", "CPU Utilization of the mobility controller in percentge", []string{"name"}, nil)
	mcMemFree        = prometheus.NewDesc("mc_mem_free", "Amount of free memory of mobility controller", []string{"name"}, nil)
	mcMemTotal       = prometheus.NewDesc("mc_mem_total", "Total amount of  memory of mobility controller", []string{"name"}, nil)
	mcUptime         = prometheus.NewDesc("mc_uptime", "Uptime of the mobility controller in seconds", []string{"name"}, nil)

	switchClientCount    = prometheus.NewDesc("switch_client_count", "Number of clients connected to switch", []string{"name"}, nil)
	switchCpuUtilization = prometheus.NewDesc("switch_cpu_utilization", "CPU Utilization of the switch in percentge", []string{"name"}, nil)
	switchUsage          = prometheus.NewDesc("switch_usage", "Bandwidth usage of the switch", []string{"name"}, nil)
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
	ch <- apClientCount
	ch <- apCpuUtilization
	ch <- apMemFree
	ch <- apMemTotal
	ch <- apUptime

	ch <- mcCpuUtilization
	ch <- mcMemFree
	ch <- mcMemTotal
	ch <- mcUptime

	ch <- switchClientCount
	ch <- switchCpuUtilization
	ch <- switchUsage
	ch <- switchUptime
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	refreshToken(e)
	listSwitches(e, ch)
	listAccessPoints(e, ch)
	listMobilityControllers(e, ch)

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

func listAccessPoints(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v2/aps?calculate_total=true&calculate_client_count=true&calculate_ssid_count=true&show_resource_details=true"

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
	var apResponse ApResponse
	if err := json.Unmarshal(body, &apResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, a := range apResponse.AccessPoints {

		ch <- prometheus.MustNewConstMetric(apClientCount, prometheus.GaugeValue, float64(a.ClientCount), a.Name)
		ch <- prometheus.MustNewConstMetric(apCpuUtilization, prometheus.GaugeValue, float64(a.CpuUtilization), a.Name)
		ch <- prometheus.MustNewConstMetric(apMemFree, prometheus.GaugeValue, float64(a.MemFree), a.Name)
		ch <- prometheus.MustNewConstMetric(apMemTotal, prometheus.GaugeValue, float64(a.MemTotal), a.Name)
		ch <- prometheus.MustNewConstMetric(apUptime, prometheus.GaugeValue, float64(a.Uptime), a.Name)
	}

}

func listMobilityControllers(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/mobility_controllers?calculate_total=false"

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
	var mcResponse McResponse
	if err := json.Unmarshal(body, &mcResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, m := range mcResponse.MobilityControllers {

		ch <- prometheus.MustNewConstMetric(mcCpuUtilization, prometheus.GaugeValue, float64(m.CpuUtilization), m.Name)
		ch <- prometheus.MustNewConstMetric(mcMemFree, prometheus.GaugeValue, float64(m.MemFree), m.Name)
		ch <- prometheus.MustNewConstMetric(mcMemTotal, prometheus.GaugeValue, float64(m.MemTotal), m.Name)
		ch <- prometheus.MustNewConstMetric(mcUptime, prometheus.GaugeValue, float64(m.Uptime), m.Name)
	}

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

		ch <- prometheus.MustNewConstMetric(switchClientCount, prometheus.GaugeValue, float64(s.ClientCount), s.Name+"_"+s.MacAddress)
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
