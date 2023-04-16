.PHONY: build clean deploy gomodgen deps

deps:
	go mod tidy
	go mod download

build: clean deps
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/getFoods food/getFoods/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/createMenu food/createMenu/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/updateMenu food/updateMenu/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/getMenu food/getMenu/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/getMenus food/getMenus/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/table/createTable table/createTable/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/table/getTable table/getTable/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/table/getTables table/getTables/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/table/updateTable table/updateTable/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/createFood food/createFood/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/order/createOrder order/createOrder/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/food/getFood food/getFood/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/order/getOrder order/getOrder/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/order/getOrders order/getOrders/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/orderItems/createOrderItem orderItems/createOrderItem/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/orderItems/getOrderItem orderItems/getOrderItem/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/orderItems/getOrderItems orderItems/getOrderItems/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/invoice/createInvoice invoice/createInvoice/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/invoice/getInvoice invoice/getInvoice/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/invoice/getInvoices invoice/getInvoices/main.go

clean:
	rm -rf ./bin ./vendor go.sum

start:
	sls offline --useDocker start --host 0.0.0.0

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
