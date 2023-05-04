# api_manager_cncf_tool

This is a simple example of how to use the CNCF tools to build a simple API manager.

CNCF graduated project used:
- [Envoy](https://www.envoyproxy.io/)
- [OPA](https://www.openpolicyagent.org/)
- [Prometheus](https://prometheus.io/)

plus [Grafana](https://grafana.com/) for seeing the metrics.

Find out what graduated project means [here](https://www.cncf.io/projects/#:~:text=Project%20maturity%20levels).

![Grafana Screen](https://github.com/paranoiasystem/api_manager_cncf_tool/blob/main/grafana_screen.png?raw=true)


## How to run

```bash
docker compose up -d
curl --location 'http://localhost:10000/api/rnk/character/1' --header 'x-auth-token: peppe'
```