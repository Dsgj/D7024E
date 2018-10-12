package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "DHT client CLI"
	app.Usage = "Store and fetch files to/from the network"

	app.Commands = []cli.Command{
		{
			Name:  "store",
			Usage: "Stores the given file on the network, returns the address",
			Action: func(c *cli.Context) error {
				path := c.Args().Get(0)
				// should probably validate path
				fileBytes, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				resp, err := resty.R().
					SetBody(fileBytes).
					SetContentLength(true).
					Post("http://localhost:3000/files")
				if err != nil {
					fmt.Printf("\nError: %v", err)
					return err
				}
				if resp.StatusCode() != 200 {
					fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
					fmt.Printf("\nResponse Status: %v", resp.Status())
					fmt.Printf("\nResponse Time: %v", resp.Time())
					fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
					fmt.Printf("\nResponse Body:%v\n", resp)
				} else {
					fmt.Printf("\n%v\n", resp)
				}
				return nil
			},
		},
		{
			Name:  "cat",
			Usage: "fetches a file from the network and prints content",
			Action: func(c *cli.Context) error {
				// should probably validate key
				key := c.Args().Get(0)
				// GET request
				resp, err := resty.R().SetPathParams(map[string]string{
					"key": key,
				}).Get("http://localhost:3000/files/{key}")

				if err != nil {
					fmt.Printf("\nError: %v", err)
					return err
				}
				if resp.StatusCode() != 200 {
					fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
					fmt.Printf("\nResponse Status: %v", resp.Status())
					fmt.Printf("\nResponse Time: %v", resp.Time())
					fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
					fmt.Printf("\nResponse Body:%v\n", resp)
				} else {
					fmt.Printf("\n%v\n", resp)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
