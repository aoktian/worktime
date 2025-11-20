go build -o worktime.exe

SET GOPROXY=https://goproxy.cn
set GOARCH=amd64
set GOOS=linux

go build -o worktime

pause