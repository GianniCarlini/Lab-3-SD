# Lab-3-SD
Orden Maquinas:
Máquina 1 - broker
ip/hostname: 10.10.28.67
Comando: make bro

Máquina 2 - Dns1
ip/hostname: 10.10.28.68
Comando: make d1

Máquina 3 - Dns2
ip/hostname: 10.10.28.69
Comando: make d2

Máquina 4 - Dns3
ip/hostname: 10.10.28.7
Comando: make d3

Comando para admin y cliente: make adm/cliente.

El sistema funciona con el dns1 como cordinador por lo cual se sugiere dejar este como el ultimo en iniciar (lleva el contador de 5mins), se crean archivos log.txt en blanco para que no se caiga el sistema. Los comando son Get/Create/Delete/Update. Si necesita modificar el tiempo de merge en la linea 54 del dn1 se cambia el tiempo (esta en segundos). El cliente aveces se cae si el broker lo conecta a un dns que no tiene el dominio (si pasa esto intentar denuevo o mejor esperar el merge)

en caso de necesitar recompilar el proto
 export GO111MODULE=on
 go get github.com/golang/protobuf/protoc-gen-go
 go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0
 export PATH="$PATH:$(go env GOPATH)/bin"
 protoc --go_out=plugins=grpc:proto helloworld.proto
