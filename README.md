<h3>Usage:</h3>

	Usage: ./aruba_exporter [options]
	Options:
 		-f string
    		Specify config file (default "exporter_config.yaml")
  		-v
			Enable verbose mode - prints HTTP status code and response headers to the terminal

If no configuration file is specified then the default of exporter_config.yaml will be assumed. The application reads the necessary credentials and configuration options from this file, and also writes the new token values to the file (as each access token expires after 2 hours)

<h4>exporter_config.yaml</h4>

	arubaEndpoint: "https://apigw-eucentral3.central.arubanetworks.com/"
	arubaTokens:
	  - arubaAccessToken: "aruba-access-token-goes-here"
	  - arubaRefreshToken: "aruba-refresh-token-goes-here"
	arubaApplicationCredentials:
	  - clientId: "aruba-application-client-id-goes-here"
	  - clientSecret: "aruba-application-client-secret-goes-here"
	exporterConfig:
	  - exporterEndpoint: "/metrics"
	  - exporterPort: ":8080"

The arubaEndpoint, exporterEndpoint and exporterPort values should also be amended to fit the required configuration.

***

<h3>Metrics:</h3>

<h4>/monitoring/v1/switches:</h4>

- switch_client_count
- switch_cpu_utilization
- switch_mem_free
- switch_mem_total
- switch_usage
- switch_uptime

<h4>/monitoring/v2/clients:</h4>

- client_rx_data_bytes
- client_tx_data_bytes

<h4>/monitoring/v1/mobility_controllers:</h4>

- mc_cpu_utilization
- mc_mem_free
- mc_mem_total
- mc_uptime

<h4>/monitoring/v2/aps:</h4>

- ap_client_count
- ap_cpu_utilization
- ap_radio_tx_power
- ap_radio_utilization
- ap_mem_free
- ap_mem_total
- ap_uptime

<h4>/branchhealth/v1/site:</h4>

- aruba_site_connected_count
- aruba_site_device_down
- aruba_site_device_high_ch_2_4ghz
- aruba_site_device_high_ch_5ghz
- aruba_site_device_high_cpu
- aruba_site_device_high_mem
- aruba_site_device_high_noise_2_4ghz
- aruba_site_device_high_noise_5ghz
- aruba_site_device_up
- aruba_site_wired_cpu_high
- aruba_site_wired_device_status_down
- aruba_site_wired_device_status_up
- aruba_site_wired_mem_high
- aruba_site_wlan_cpu_high
- aruba_site_wlan_device_status_down
- aruba_site_wlan_device_status_up
- aruba_site_wlan_mem_high


***

<h3>Prometheus Configuration:</h3>

For Prometheus configuration, it should be noted that the scraping interval greatly depends on the daily API call limit which difers per organisation. Each time the data is scraped, 5 API calls are made, inlcuding an additional 12 API calls per day for refresh tokens. For example, setting the interval at 30 seconds should result in 14,412 calls per day.

