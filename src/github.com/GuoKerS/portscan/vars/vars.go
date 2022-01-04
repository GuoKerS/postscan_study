package vars

import "sync"

var (
	Host string
	Port = "22"
	//Mode string
	Timeout   = 2
	ThreadNum = 100

	// 第一种存活map
	//Scan_Task = map[string]map[string]int{
	//	"ip": {
	//		"isPingOK": 0,
	//	},
	//}

	// 第二种 简单存活map （并发安全map）
	IsPingsOK *sync.Map
	Result    *sync.Map
)

func init() {
	Result = &sync.Map{}
	IsPingsOK = &sync.Map{}
}
