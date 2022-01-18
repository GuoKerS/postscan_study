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
	// todo 根据ip列表先做一遍存活判断
	ipsSurvival, _ := scaner.GetSurviving_IPs(ips)

	ports, err := scaner.GetPorts(vars.Port)

	// todo 生成存活列表在做端口列表生成(对列表中的元素进行删除效率很低，可能需要设计一个链表出来，用来给接下来的扫描做准备)
	tasks, _ := scaner.Gen_PortScanTask(ipsSurvival, ports)

	scaner.RunTask(tasks)
	scaner.PrintResult()
	return err
}
