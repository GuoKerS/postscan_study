echo "Building ....."

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags="-w -s" -o ../../../bin/linux_x86_postscan main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ../../../../bin/linux_x64_postscan main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ../../../../bin/windows_x64_postscan.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags="-w -s" -o ../../../../bin/windows_x86_postscan.exe main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ../../../../bin/macos_postscan main.go