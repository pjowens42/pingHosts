useing: github.com/gofiber/fiber/v2 github.com/prometheus-community/pro-bing (to run on Linux sudo sysctl -w net.ipv4.ping_group_range="0 2147483647")

go run cmd/main.go --interval 5s --hosts www.google.com --port 4000