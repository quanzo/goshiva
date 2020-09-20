SET GOOS=windows
SET GOARCH=386
go build -o shiva-win-386.exe
SET GOARCH=amd64
go build -o shiva-win-64.exe
SET GOOS=linux
SET GOARCH=386
go build -o shiva-linux-386
SET GOARCH=amd64
go build -o shiva-linux-64
