package main

import(
	"flag"
	"fmt"
	"strings"
	"time"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus-community/pro-bing"
)

type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *StringSlice) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

type ServerConfig struct{
	port string
	interval time.Duration
	hosts StringSlice
}

type HostPing struct{
	ID int `json: id`
	Host string `json: host`
	Date time.Time `json: date`
	Stats probing.Statistics  `json: Stats`
}

func main () {
	hostPings := []HostPing{}

	var serverConfig ServerConfig
	flag.StringVar(&serverConfig.port, "port", "3000", "Server Port")
	flag.DurationVar(&serverConfig.interval, "interval", 5*time.Second, "Time interval")
	flag.Var(&serverConfig.hosts,"hosts", "Comma-separated list of hostnames")
	flag.Parse()
	
	serverConfig.hosts = append(serverConfig.hosts, "127.0.0.1")

	fmt.Printf("%+v\n", serverConfig)
	ticker := time.NewTicker(serverConfig.interval)
	go func(){
		for{
			select {
			case t := <-ticker.C:
				fmt.Println("Ticker at:", t)
				for index, value := range serverConfig.hosts {
				fmt.Printf("Index: %d, Value: %s\n", index, value)
				pinger, err := probing.NewPinger(value)
				if err != nil {
					panic(err)
				}
				pinger.Count = 1
				err = pinger.Run() // Blocks until finished.
				if err != nil {
					panic(err)
				}
				stats := pinger.Statistics()
				hostPing := &HostPing{
					ID: len(hostPings) + 1,
					Host: value, 
					Date: time.Now(), 
					Stats:  *stats,
				}
				hostPings = append(hostPings, *hostPing)
			}
			}
		}
	}()

	app := fiber.New()

	app.Get("/api/interval", func (c *fiber.Ctx) error {
        return c.Status(200).JSON(fiber.Map{"interval": serverConfig.interval})
    })

    app.Get("/api/hostPings", func (c *fiber.Ctx) error {
        return c.Status(200).JSON(fiber.Map{"body": hostPings})
    })

    log.Fatal(app.Listen(":" + serverConfig.port))
}