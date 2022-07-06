k6: go.mod go.sum *.go
	xk6 build --with github.com/grafana/xk6-output-prometheus-remote --with xk6-temporal=.

grafana/dashboards/%.json: grafana/dashboards/%.yml
	grabana render -i $< > $@

dashboards: grafana/dashboards/*.json
