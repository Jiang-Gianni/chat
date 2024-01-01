dg:
	d2 --watch img/gRPC.d2 img/gRPC.svg

dn:
	d2 --watch img/NATS.d2 img/NATS.svg

ds:
	d2 --watch img/ssr.d2 img/ssr.svg

dgc:
	d2 --watch img/gRPCchat.d2 img/gRPCchat.svg

dnc:
	d2 --watch img/NATSchat.d2 img/NATSchat.svg

dgl:
	d2 --watch img/gRPClogin.d2 img/gRPClogin.svg

dnl:
	d2 --watch img/NATSlogin.d2 img/NATSlogin.svg

cover:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o=cover.html