while true
do
    watchman-wait -p "**/*.go" -- .
    clear
    go test -cover -count=1 -timeout 10s -race ./...
done
