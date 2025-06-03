go build -o pkg/worker.exe cmd/worker/main.go
go build -o pkg/api.exe cmd/api/main.go

Start-Process -FilePath ".\pkg\api.exe"
Start-Process -FilePath ".\pkg\worker.exe"