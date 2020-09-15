package command

import (
	"fmt"
	"github.com/magazin/data"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "magazin",
	Usage: "why not do shopping from your command-line?",
	Before: func(c *cli.Context) error {
		return data.Init()
	},
	Action: func(c *cli.Context) error {
		return cli.ShowAppHelp(c)
	},
	After: func(c *cli.Context) error {
		return nil
	},
}

func Main(args []string) error {
	app.Commands = []*cli.Command{
		buyCommand,
	}

	if err :=  app.Run(args); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
