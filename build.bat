rsrc -manifest exe.manifest -ico main.ico
go-bindata -o icon_files.go main.ico status.ico
go build -buildvcs=false -ldflags="-H windowsgui -w -s" -o simple-ip-monitor-windows.exe