package main

import (
	"fmt"
	"github.com/GuoKerS/portscan/cmd"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "port_scan v0.1"
	app.Author = "Guoker"
	app.Version = "2021/12/31"
	app.Usage = "tcp connect port scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)

	//调试用
	test := os.Args
	//test = append(test, "scan")

	err := app.Run(test)

	if err != nil {
		fmt.Println(err)
	}
}
