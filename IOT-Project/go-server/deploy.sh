env GOOS=linux GOARCH=arm GOARM=5 go build
scp go-server mndkk.dk:
ssh -t mndkk.dk './go-server'
