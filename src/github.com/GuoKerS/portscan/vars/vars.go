package vars

import "sync"

var (
	Host string
	Port = "22"
	//Mode string
	Timeout   = 2
	ThreadNum = 100
	Result    *sync.Map
)

func init() {
	Result = &sync.Map{}
}
