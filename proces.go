package main

import (
	"fmt"

	"os"

	"strconv"

	"time"

	"github.com/shirou/gopsutil/process"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please Input PID")
		os.Exit(0)
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Please Input pid")
		os.Exit(0)
	}
	ps, _ := process.Pids()
	count := len(ps)
	ok := false
	for i := 0; i < count; i++ {
		if ps[i] == int32(pid) {
			ok = true
		}
	}
	if !ok {
		fmt.Println("Please Input pid")
		os.Exit(0)
	}
	p := process.Process{Pid: int32(pid)}
	a, _ := p.MemoryPercent()
	mem, _ := p.MemoryInfo()
	fmt.Printf("Memory Status:\n------------------\nRS:%d | VMS:%d | swap:%d | MemUsed:%f%%\n", mem.RSS, mem.VMS, mem.Swap, a)
	//fmt.Println(mem)
	fmt.Println("================")
	cpu, _ := p.Percent(1000000000)
	fmt.Printf("CPUUsed: %f%%\n", cpu)
	fmt.Println("================")
	c, _ := p.Cmdline()
	fmt.Printf("Cmd: %s\n", c)
	fmt.Println("================")
	u, _ := p.Username()
	fmt.Printf("UserName: %s\n", u)
	fmt.Println("================")
	t, _ := p.CreateTime()
	ti := time.Unix(t/1000, 0)
	fmt.Println("CreateTime:", ti)
	fmt.Println("================")

}
