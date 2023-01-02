package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

var (
	cli                    client.Client
	lastNetIOStatTimeStamp int64    // 上一次获取网络IO数据的时间点
	lastNetIO              *NetInfo //上一次的网络IO数据
)

// connect
func initconnInflux() (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	return
}

// insert
/*func writesCpuPoints(data *CpuInfo) {
	//data := info.Data.(CpuInfo)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "monitor",
		Precision: "s",	// 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型插入数据
		tags := map[string]string{"cpu": "cpu0"}
		fields := map[string]interface{}{
			"cpu_percent":	data.CpuPercent,
		}
	pt, err := client.NewPoint("cpu_percent",tags,fields,time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert cpu info success")

}*/

/*func writesMemPoints(data *MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "monitor",
		Precision: "s",	// 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型插入数据
	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":	int64(data.Total),
		"available": int64(data.Available),
		"used": int64(data.Used),
		"used_percent": data.UsedPercent,
	}
	pt, err := client.NewPoint("memory",tags,fields,time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert mem info success")

}*/

func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	//CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v\n", percent)
	//写入到influxDB中
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo)
}

func getMemInfo() {
	var memInfo = new(MemInfo)
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("get mem info failed, %v:", err)
		return
	}
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo)
}

func getDiskInfo() {
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		//拿到每一个分区
		usageStat, err := disk.Usage(part.Mountpoint) //传挂载点
		if err != nil {
			fmt.Printf("get %s usage stat failed,err:%v", err)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	writesDiskPoints(diskInfo)
}

func getNetInfo() {
	var netInfo = &NetInfo{
		NetIoCountersStat: make(map[string]*IOStat, 8),
	}
	currentTimeStamp := time.Now().Unix()
	netIOs, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("get net io counters failed, err:%v", err)
		return
	}
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv
		//将具体网卡数据的ioStat变量添加到map中
		netInfo.NetIoCountersStat[netIO.Name] = ioStat // 不要放到continue下面
		// 开始计算网卡相关的速率
		if lastNetIOStatTimeStamp == 0 || lastNetIO == nil {
			continue
		}
		//计算时间间隔
		interval := currentTimeStamp - lastNetIOStatTimeStamp
		//计算速率
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetIO.NetIoCountersStat[netIO.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetIO.NetIoCountersStat[netIO.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetIO.NetIoCountersStat[netIO.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetIO.NetIoCountersStat[netIO.Name].PacketsRecv)) / float64(interval)

	}
	//更新全局记录的上一次采集网卡的时间点和网卡数据
	lastNetIOStatTimeStamp = currentTimeStamp //更新时间
	lastNetIO = netInfo
	// 发送
	writesNetPoints(netInfo)
}

func run(interval time.Duration) {
	ticker := time.Tick(interval)
	for _ = range ticker {
		getCpuInfo()
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}

func main() {
	err := initconnInflux()
	if err != nil {
		fmt.Printf("connect to influcDB failed, err:%v", err)
		return
	}
	run(time.Second)
}
