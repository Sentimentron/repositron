package main

import (
	"github.com/Sentimentron/repositron/client/go/repoclient"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os"
)

var qs = []*survey.Question{
	{
		Name: "location",
		Prompt: &survey.Input{
			Message: "Where should the config file live?",
			Default: repoclient.BuildDefaultClientConfigurationPath(),
		},
	},
	{
		Name: "URL",
		Prompt: &survey.Input{
			Message: "What's the base URL?",
			Help:    "e.g. http://localhost:8000/repositron",
			Default: "http://localhost:8000",
		},
	},
}

func main() {

	app := cli.NewApp()
	app.Name = "repositron-cli"
	app.Usage = "Control a Repositron server"
	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from FILE",
			Value: repoclient.BuildDefaultClientConfigurationPath(),
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "configure",
			Usage: "Complete repositron setup",
			Action: func(c *cli.Context) error {
				answers := struct {
					ConfigPath string `survey:"location"`
					BaseURL    string `survey:"URL"`
				}{}
				err := survey.Ask(qs, &answers)
				if err != nil {
					return err
				}

				config := repoclient.ClientConfiguration{answers.BaseURL, "1"}
				repoclient.WriteClientConfiguration(&config, answers.ConfigPath)
				log.Printf("Successfully wrote configuration to %s", answers.ConfigPath)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
