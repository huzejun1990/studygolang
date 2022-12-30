package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"time"
)

//获取cpu信息

// cpu info
func getCpuInfo() {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Println("get cpu info failed,err:%v", err)
	}
	for _, c1 := range cpuInfos {
		fmt.Println(c1)
	}
	// CPU使用率
	for {
		percent, _ := cpu.Percent(time.Second, false)
		fmt.Printf("cpu percent:%v\n", percent)
	}

}

func main() {
	//getCpuInfo()
}
