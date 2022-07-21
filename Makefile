k6: go.mod go.sum *.go client/*.go logger/*.go metrics/*.go
	xk6 build --with github.com/temporalio/xk6-temporal=.
