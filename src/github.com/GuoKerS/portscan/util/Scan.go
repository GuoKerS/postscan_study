package util

import (
	"github.com/GuoKerS/portscan/scaner"
	"github.com/GuoKerS/portscan/vars"
	"github.com/urfave/cli"
)

func Scan(ctx *cli.Context) error {
	//fmt.Println("util.Scan被运行")
	if ctx.IsSet("iplist") {
		vars.Host = ctx.String("iplist")
		//fmt.Println(vars.Host)
	}

	if ctx.IsSet("ports") {
		vars.Port = ctx.String("ports")
		//fmt.Println(vars.Port)
	}

	if ctx.IsSet("timeout") {
		vars.Timeout = ctx.Int("timeout")
		//fmt.Println(vars.Timeout)
	}

	if ctx.IsSet("concurrency") {
		vars.ThreadNum = ctx.Int("concurrency")
		//fmt.Println(vars.ThreadNum)
	}

	ips, err := scaner.GetIps(vars.Host)

	ports, err := scaner.GetPorts(vars.Port)
	tasks, _ := scaner.Gen_PortScanTask(ips, ports)
	scaner.RunTask(tasks)
	scaner.PrintResult()
	return err
}
