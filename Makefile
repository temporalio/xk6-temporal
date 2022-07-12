GRAFANA_URL ?= http://localhost:3000

k6: go.mod go.sum *.go client/*.go worker/*.go logger/*.go metrics/*.go
	xk6 build --with github.com/grafana/xk6-output-prometheus-remote --with github.com/temporalio/xk6-temporal=.

dashboards:
	grabana apply -g ${GRAFANA_URL} -f xk6-temporal -i grafana/dashboards/k6-temporal.yml
	grabana apply -g ${GRAFANA_URL} -f Temporal -i grafana/dashboards/workers.yml
