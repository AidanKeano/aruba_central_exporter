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

The credentials necessary to authenticate with Aruba Central should be located in two seperate yaml files located in the same directory as the exporter. The first should be named client.yaml and contain the client id and client secret created from aruba central.

<h4>client.yaml:</h4>

	clientId: client-id-goes-here
	clientSecret: client-secret-goes-here

For the other file named tokens.yaml, the initial access and refresh tokens obtained from the Aruba Central API Gateway should be entered here. This application will overwrite this data after refreshing tokens, as each token expires after two hours.

<h4>tokens.yaml:</h4>

	arubaAccessToken: aruba-access-token-goes-here
	arubaRefreshToken: aruba-refresh-token-goes-here
