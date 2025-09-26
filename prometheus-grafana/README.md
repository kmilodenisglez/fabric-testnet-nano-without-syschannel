# Description

This sample provides an environment to display and capture metrics from the test-network in real time. It consists of a docker-compose file that starts a Prometheus and Grafana server setup configured to collect and display metrics for the test network.

## Requirements

This sample has been tested and is recommended to be used on **linux** in order to fully benefit from its capabilities, however it can be deployed and works on MacOS-intel machine as well (some modification to the cadvisor docker image and related queries are required to show docker containers metrics).
You will need to have installed **docker-compose with version 1.29 or above** (note that this is higher than the v1.14 requirement requested for the test-network).

## How to use

1. Go to the [fabric-testnet-nano-without-syschannel](../) directory and follow instructions in README.md to run bring up the network
2. Bring up the Prometheus/Grafana network in the fabric-testnet-nano-without-syschannel/prometheus-grafana directory and run **docker-compose up -d**
3. Log in: type "localhost:3000" on your web browser -> username="admin", password="admin" -> set a new password
4. Browse dashboard and analyse results
   - The default dashboard "HLF Performances Review" can be found and displayed by hovering over the dashboard menu and clicking on the browse button.
   Once opened the dashboard, to display the collected metrics and data, adjust the timeframe on the top right to focus on the latest timespan when the network was up.
5. Deploy a chaincode, start using the network and use the Grafana dashboard to analyse and assess your network performances.
Extras: add new queries, modify dashboard & add relevant changes to main repo --> extract json and add it to "Grafana/dashboards/hlf-performances.json".
Metrics can also be displayed directly from Prometheus by going to "localhost:9090".

## Docker Compose

Brings up

- a Prometheus server (port 9090) -> pulls metrics from peers, orderer, system(node exporter) and containers(cadvisor)
- Grafana server (port 3000) -> collects and display data from Prometheus
- node exporter (port 9100) -> exposes systems metrics
- cadvisor (port 8080) -> exposes docker containers metrics

## Prometheus "configuration file"

### Prometheus.yml

Fabric metrics targets:

- ordere1: `localhost:8443`
- peer1:   `localhost:8446`
- peer2:   `localhost:8448`

System and docker metrics targets:

- `cadvisor:8080`
- `node-exporter:9100`

Check the state of the connections with targets on http://localhost:9090/targets.

## Sources

[Prometheus docs](https://prometheus.io/docs/introduction/overview/)
[Grafana docs](https://grafana.com/docs/)
