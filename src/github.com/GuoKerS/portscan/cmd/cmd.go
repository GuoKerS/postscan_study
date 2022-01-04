package cmd

import (
	"github.com/GuoKerS/portscan/util"
	"github.com/urfave/cli"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to scan port",
	Description: "start to scan port",
	Action:      util.Scan,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "iplist, i", Value: "10.0.0.1/24", Usage: "ip list"},
		cli.StringFlag{Name: "ports, p", Value: "22", Usage: "port list"},
		//cli.StringFlag{Name: "mode, m", Value: "", Usage: "scan mode"},
		cli.IntFlag{Name: "timeout, t", Value: 3, Usage: "timeout"},
		cli.IntFlag{Name: "threads, c", Value: 100, Usage: "threads max"},
	},
}
