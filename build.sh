mkdir -p backend
export GO111MODULE=on
go get
go build -o backend/goapi ./...
chmod u+x backend/*