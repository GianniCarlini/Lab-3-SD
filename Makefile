d1:
	cd dns1 && \
	go run dns1.go

d2:
	cd dns2 && \
	go run dns2.go

d3:
	cd dns3 && \
	go run dns3.go

bro:
	cd broker && \
	go run broker.go

adm:
	go run Admin.go

cliente:
	go run Client.go