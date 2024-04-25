<h3>Metrics:</h3>

<h4>/monitoring/v1/switches:</h4>

- switch_client_count
- switch_cpu_utilization
- switch_usage
- switch_uptime

<h4>/monitoring/v1/mobility_controllers:</h4>

- mc_cpu_utilization
- mc_mem_free
- mc_mem_total
- mc_uptime

<h4>/monitoring/v2/aps:</h4>

- ap_client_count
- ap_cpu_utilization
- ap_mem_free
- ap_mem_total
- ap_uptime



***
<h3>Usage:</h3>

./aruba_central_exporter

The configuration options should be specified in the exporter_config.yaml file in the same directory as the exporter. The application reads the necessary credentials and configuration options from this file, and also writes the new token values to the file (as each access token expires after 2 hours)

<h4>exporter_config.yaml</h4>

	arubaEndpoint: "https://apigw-eucentral3.central.arubanetworks.com/" //Replace with alternative endpoint if necessary
	arubaTokens:
	  - arubaAccessToken: "aruba-access-token-goes-here"
	  - arubaRefreshToken: "aruba-refresh-token-goes-here"
	arubaApplicationCredentials:
	  - clientId: "aruba-application-client-id-goes-here"
	  - clientSecret: "aruba-application-client-secret-goes-here"
	exporterConfig:
	  - exporterEndpoint: "/metrics" //Replace with alternative endpoint directory if necessary
	  - exporterPort: ":8080" //Choose whatever port you wish here
