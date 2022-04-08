package util

import (
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"os"
	"os/signal"
)

func MonitorCtrlC() {
	sigChan := make(chan os.Signal)
	defer close(sigChan)
	signal.Notify(sigChan)
	select {
	case sig := <-sigChan:
		//fmt.Printf("[!] Process Exit.msg:%s\n", sig)
		_ = sig
		fmt.Printf("[Seed] %d | [Index] %d | [Host] %s\n", vars.Seed, vars.Index, vars.Host)
		panic("[!] Process Exit.msg")

	}
}
