go mod init example
go get github.com/namihq/walrus-go
go run ingest-file-walrus.go --file ./sample.webp --epochs 1
