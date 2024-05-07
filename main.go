package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

type TopNClientResponse struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	MacAddress  string `json:"macaddr"`
	Name        string `json:"name"`
	RxDataBytes int    `json:"rx_data_bytes"`
	TxDataBytes int    `json:"tx_data_bytes"`
}

var (
	apClientCount    = prometheus.NewDesc("ap_client_count", "Number of clients connected to access point", []string{"name", "groupName", "site", "status", "firmwareVersion"}, nil)
	apCpuUtilization = prometheus.NewDesc("ap_cpu_utilization", "CPU Utilization of the access point in percentge", []string{"name", "groupName", "site", "status", "firmwareVersion"}, nil)
	apMemFree        = prometheus.NewDesc("ap_mem_free", "Amount of free memory of access point", []string{"name", "groupName", "site", "status", "firmwareVersion"}, nil)
	apMemTotal       = prometheus.NewDesc("ap_mem_total", "Total amount of  memory of access point", []string{"name", "groupName", "site", "status", "firmwareVersion"}, nil)
	apUptime         = prometheus.NewDesc("ap_uptime", "Uptime of the access point in seconds", []string{"name", "groupName", "site", "status", "firmwareVersion"}, nil)

	apRadioTxPower     = prometheus.NewDesc("ap_radio_tx_power", "Radio tx power", []string{"band", "channel", "radioName", "apName"}, nil)
	apRadioUtilization = prometheus.NewDesc("ap_radio_utilization", "Radip cpu utilization", []string{"band", "channel", "radioName", "apName"}, nil)

	clientRxDataBytes = prometheus.NewDesc("client_rx_data_bytes", "Volume of data received", []string{"name", "mac"}, nil)
	clientTxDataBytes = prometheus.NewDesc("client_tx_data_bytes", "Volume of data transmitted", []string{"name", "mac"}, nil)

	mcCpuUtilization = prometheus.NewDesc("mc_cpu_utilization", "CPU Utilization of the mobility controller in percentge", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcMemFree        = prometheus.NewDesc("mc_mem_free", "Amount of free memory of mobility controller", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcMemTotal       = prometheus.NewDesc("mc_mem_total", "Total amount of  memory of mobility controller", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)
	mcUptime         = prometheus.NewDesc("mc_uptime", "Uptime of the mobility controller in seconds", []string{"name", "groupName", "mode", "model", "site", "status", "firmwareVersion"}, nil)

	switchClientCount    = prometheus.NewDesc("switch_client_count", "Number of clients connected to switch", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)
	switchCpuUtilization = prometheus.NewDesc("switch_cpu_utilization", "Current Switch CPU utilization percentage", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)
	switchMemFree        = prometheus.NewDesc("switch_mem_free", "Switch free memory", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)
	switchMemTotal       = prometheus.NewDesc("switch_mem_total", "Switch total memory", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)
	switchUsage          = prometheus.NewDesc("switch_usage", "Switch uptime", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)
	switchUptime         = prometheus.NewDesc("switch_uptime", "Switch usage", []string{"name", "stackMemberId", "groupId", "groupName", "site", "siteId", "switchRole", "switchType", "status", "firmwareVersion"}, nil)

	expiresIn = 0
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

	ch <- apRadioTxPower
	ch <- apRadioUtilization

	ch <- clientRxDataBytes
	ch <- clientTxDataBytes

	ch <- mcCpuUtilization
	ch <- mcMemFree
	ch <- mcMemTotal
	ch <- mcUptime

	ch <- switchClientCount
	ch <- switchCpuUtilization
	ch <- switchMemFree
	ch <- switchMemTotal
	ch <- switchUsage
	ch <- switchUptime
}

func decrementExpiresIn() {
	for {
		time.Sleep(time.Second)
		expiresIn--
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	refreshToken(e)
	listSwitches(e, ch)
	listAccessPoints(e, ch)
	listMobilityControllers(e, ch)
	listTopClients(e, ch)

}

func main() {

	go decrementExpiresIn()
	config := Config{}
	readConfig(&config)

	arubaEndpoint := config.ArubaEndpoint
	arubaAccessToken := config.ArubaTokens[0].ArubaAccessToken
	arubaRefreshToken := config.ArubaTokens[1].ArubaRefreshToken
	exporterEndpoint := config.ExporterConfig[0].ExporterEndpoint
	exporterPort := config.ExporterConfig[1].ExporterPort

	exporter := NewExporter(arubaEndpoint, arubaAccessToken, arubaRefreshToken)
	prometheus.MustRegister(exporter)

	http.Handle(exporterEndpoint, promhttp.Handler())

	err := http.ListenAndServe(exporterPort, nil)

	if err != nil {
		if err.Error() == "listen tcp :8080: bind: address already in use" {
			fmt.Println("Error: Port", exporterPort, "is already in use.")
		} else {
			fmt.Println("Error starting server:", err)
		}
	}

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

		ch <- prometheus.MustNewConstMetric(apClientCount, prometheus.GaugeValue, float64(a.ClientCount), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(apCpuUtilization, prometheus.GaugeValue, float64(a.CpuUtilization), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(apMemFree, prometheus.GaugeValue, float64(a.MemFree), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(apMemTotal, prometheus.GaugeValue, float64(a.MemTotal), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(apUptime, prometheus.GaugeValue, float64(a.Uptime), a.Name, a.GroupName, a.Site, a.Status, a.FirmwareVersion)

		for _, r := range a.Radios {

			ch <- prometheus.MustNewConstMetric(apRadioTxPower, prometheus.GaugeValue, float64(r.TxPower), strconv.Itoa(r.Band), r.Channel, r.RadioName, a.Name)
			ch <- prometheus.MustNewConstMetric(apRadioUtilization, prometheus.GaugeValue, float64(r.Utilization), strconv.Itoa(r.Band), r.Channel, r.RadioName, a.Name)
		}
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

		ch <- prometheus.MustNewConstMetric(mcCpuUtilization, prometheus.GaugeValue, float64(m.CpuUtilization), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcMemFree, prometheus.GaugeValue, float64(m.MemFree), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcMemTotal, prometheus.GaugeValue, float64(m.MemTotal), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(mcUptime, prometheus.GaugeValue, float64(m.Uptime), m.Name, m.GroupName, m.Mode, m.Model, m.Site, m.Status, m.FirmwareVersion)
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

		ch <- prometheus.MustNewConstMetric(switchClientCount, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(switchCpuUtilization, prometheus.GaugeValue, float64(s.CPUUtilization), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(switchMemFree, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(switchMemTotal, prometheus.GaugeValue, float64(s.ClientCount), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(switchUsage, prometheus.GaugeValue, float64(s.Usage), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
		ch <- prometheus.MustNewConstMetric(switchUptime, prometheus.GaugeValue, float64(s.Uptime), s.Name, strconv.Itoa(s.StackMemberID), strconv.Itoa(s.GroupID), s.GroupName, s.Site, strconv.Itoa(s.SiteID), strconv.Itoa(s.SwitchRole), s.SwitchType, s.Status, s.FirmwareVersion)
	}
}

func listTopClients(e *Exporter, ch chan<- prometheus.Metric) {

	url := e.arubaEndpoint + "monitoring/v1/clients/bandwidth_usage/topn?count=100"

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
	var topNClientResponse TopNClientResponse
	if err := json.Unmarshal(body, &topNClientResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, t := range topNClientResponse.Clients {

		ch <- prometheus.MustNewConstMetric(clientRxDataBytes, prometheus.GaugeValue, float64(t.RxDataBytes), t.Name, t.MacAddress)
		ch <- prometheus.MustNewConstMetric(clientTxDataBytes, prometheus.GaugeValue, float64(t.TxDataBytes), t.Name, t.MacAddress)

	}

}

func refreshToken(e *Exporter) {

	if expiresIn < 60 {

		config := Config{}
		readConfig(&config)

		clientId := config.ArubaApplicationCredentials[0].ClientID
		clientSecret := config.ArubaApplicationCredentials[1].ClientSecret

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

		expiresIn = tokenResponse.ExpiresIn

		e.arubaAccessToken = tokenResponse.AccessToken
		e.arubaRefreshToken = tokenResponse.RefreshToken

		configData := []byte(`arubaEndpoint: "` + config.ArubaEndpoint + `"
arubaTokens:
  - arubaAccessToken: "` + e.arubaAccessToken + `"
  - arubaRefreshToken: "` + e.arubaRefreshToken + `"
arubaApplicationCredentials:
  - clientId: "` + config.ArubaApplicationCredentials[0].ClientID + `"
  - clientSecret: "` + config.ArubaApplicationCredentials[1].ClientSecret + `"
exporterConfig:
  - exporterEndpoint: "` + config.ExporterConfig[0].ExporterEndpoint + `"
  - exporterPort: "` + config.ExporterConfig[1].ExporterPort + `" `)

		err = ioutil.WriteFile("exporter_config.yaml", configData, 0644)

		if err != nil {
			fmt.Println("Error writing config file:", err)
		}
	}

}
