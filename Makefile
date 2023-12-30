dg:
	d2 --watch img/gRPC.d2 img/gRPC.svg

dn:
	d2 --watch img/NATS.d2 img/NATS.svg

cover:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o=cover.html