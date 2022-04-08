package vars

import "sync"

var (
	Host string
	Port = "22"
	//Mode string
	Timeout   = 2
	ThreadNum = 800 // 过大可能会影响端口扫描准确里，后续考虑将存活探测和端口扫描的最大线程分开
	TaskDone  = false

	// 第一种存活map
	//Scan_Task = map[string]map[string]int{
	//	"ip": {
	//		"isPingOK": 0,
	//	},
	//}

	// 第二种 简单存活map （并发安全map）
	IsPingsOK *sync.Map
	Result    *sync.Map

	Seed  int64
	Index int64
)

func init() {
	Result = &sync.Map{}
	IsPingsOK = &sync.Map{}
}
