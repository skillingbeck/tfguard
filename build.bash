cd cli/
env GOOS=linux GOARCH=amd64 go build -o ../bin/compare_linux_amd64
env GOOS=windows GOARCH=amd64 go build -o ../bin/compare_windows_amd64.exe