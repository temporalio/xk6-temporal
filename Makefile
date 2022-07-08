k6: go.mod go.sum *.go client/*.go worker/*.go logger/*.go metrics/*.go
	xk6 build --with github.com/grafana/xk6-output-prometheus-remote --with github.com/temporalio/xk6-temporal=.

grafana/dashboards/%.json: grafana/dashboards/%.yml
	grabana render -i $< > $@

dashboards: grafana/dashboards/*.json
