// Copyright 2016 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"strconv"

	"github.com/gizak/termui"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func mem_UsedPercent(UsedPercent_mem int) *termui.Gauge {
	g := termui.NewGauge()
	g.Percent = UsedPercent_mem
	g.Width = 50
	g.Height = 3
	g.BorderLabel = "Memory UsedPercent"
	g.BarColor = termui.ColorRed
	g.BorderFg = termui.ColorWhite
	g.BorderLabelFg = termui.ColorCyan
	return g
}

func cpu_UsedPercent(UsedPercent_cpu int) *termui.Gauge {
	g := termui.NewGauge()
	g.Percent = UsedPercent_cpu
	g.Width = 50
	g.Height = 3
	g.PercentColor = termui.ColorBlue
	g.X = 51
	g.BorderLabel = "CPU UsedPercent"
	g.BarColor = termui.ColorYellow
	g.BorderFg = termui.ColorWhite
	return g
}

func mem_info(v *mem.VirtualMemoryStat) *termui.Par {
	total := strconv.Itoa(int(v.Total >> 20))
	free := strconv.Itoa(int(v.Free >> 20))
	buf := strconv.Itoa(int(v.Buffers >> 20))
	cached := strconv.Itoa(int(v.Cached >> 20))
	mem_info := "Total:" + total + "MB  " + "Free:" + free + "MB\n" + "Buffers:" + buf + "MB  " + "Cached: " + cached + "MB"
	g := termui.NewPar(mem_info)
	g.Height = 4
	g.Width = 50
	g.Y = 3
	g.BorderLabel = "Memory Info"
	g.BorderFg = termui.ColorYellow
	return g
}

func cpu_info() *termui.Par {
	v, _ := load.Avg()
	load1 := strconv.Itoa(int(v.Load1))
	load5 := strconv.Itoa(int(v.Load5))
	load15 := strconv.Itoa(int(v.Load15))
	c, _ := cpu.Times(false)
	//{"cpu":"","user":45740.2,"system":20639.8,"idle":591461.7,"nice":0.0,"iowait":0.0,"irq":0.0,"softirq":0.0,"steal":0.0,"guest":0.0,"guestNice":0.0,"stolen":0.0}
	cpu_info := "Load1:" + load1 + "  load5:" + load5 + "  load15:" + load15 + "\n" +
		"cpu:" + c[0].CPU + " user:" + strconv.FormatFloat(c[0].User, 'f', 1, 64) + " system:" + strconv.FormatFloat(c[0].System, 'f', 1, 64) +
		"\nidle:" + strconv.FormatFloat(c[0].Idle, 'f', 1, 64) + " nice:" + strconv.FormatFloat(c[0].Nice, 'f', 1, 64) +
		"\niowait:" + strconv.FormatFloat(c[0].Iowait, 'f', 1, 64) + " irq:" + strconv.FormatFloat(c[0].Irq, 'f', 1, 64) +
		"\nsoftirq:" + strconv.FormatFloat(c[0].Softirq, 'f', 1, 64) + " steal:" + strconv.FormatFloat(c[0].Steal, 'f', 1, 64) +
		"\nguest:" + strconv.FormatFloat(c[0].Guest, 'f', 1, 64) + " guestNice:" + strconv.FormatFloat(c[0].GuestNice, 'f', 1, 64) +
		" stolen:" + strconv.FormatFloat(c[0].Stolen, 'f', 1, 64)
	g := termui.NewPar(cpu_info)
	g.Height = 8
	g.Width = 50
	g.Y = 3
	g.X = 51
	g.BorderLabel = "Cpu Info"
	g.BorderFg = termui.ColorYellow
	return g
}

func disk_info(disks []disk.PartitionStat, count int) *termui.Par {

	var UsedPercent string
	var d string
	var disk_info string
	var total string
	for i := 0; i < count; i++ {
		d = disks[i].Device
		dd, _ := disk.Usage(d)
		UsedPercent = strconv.FormatFloat(dd.UsedPercent, 'f', 1, 64)
		total = strconv.Itoa(int(dd.Total) >> 30)
		disk_info = disk_info + d + " - " + "Total:" + total + "G" + "  UsedPercent:" + UsedPercent + "%\n"
	}
	g := termui.NewPar(disk_info)
	g.Height = count * 2
	g.Width = 50
	g.Y = 7
	g.BorderLabel = "Disk Info"
	g.BorderFg = termui.ColorYellow
	return g
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(1000000000, true)
	disks, _ := disk.Partitions(false)
	count := len(disks)
	ga_mem := mem_UsedPercent(int(v.UsedPercent))
	ga_cpu := cpu_UsedPercent(int(c[0]))
	gc_mem := mem_info(v)
	gc_cpu := cpu_info()
	gc_disk := disk_info(disks, count)

	termui.Render(ga_mem, gc_mem, ga_cpu, gc_cpu, gc_disk)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)
		termui.SendCustomEvt("/usr/t", t.Count)

		if t.Count%2 == 0 {
			m, _ := mem.VirtualMemory()
			ga_mem.Percent = int(m.UsedPercent)
			cu, _ := cpu.Percent(1000000000, true)
			ga_cpu.Percent = int(cu[0])
			cc, _ := load.Avg()
			load1 := strconv.Itoa(int(cc.Load1))
			load5 := strconv.Itoa(int(cc.Load5))
			load15 := strconv.Itoa(int(cc.Load15))
			c, _ := cpu.Times(false)
			cpu_info := "Load1:" + load1 + "  load5:" + load5 + "  load15:" + load15 + "\n" +
				"cpu:" + c[0].CPU + " user:" + strconv.FormatFloat(c[0].User, 'f', 1, 64) + " system:" + strconv.FormatFloat(c[0].System, 'f', 1, 64) +
				"\nidle:" + strconv.FormatFloat(c[0].Idle, 'f', 1, 64) + " nice:" + strconv.FormatFloat(c[0].Nice, 'f', 1, 64) +
				"\niowait:" + strconv.FormatFloat(c[0].Iowait, 'f', 1, 64) + " irq:" + strconv.FormatFloat(c[0].Irq, 'f', 1, 64) +
				"\nsoftirq:" + strconv.FormatFloat(c[0].Softirq, 'f', 1, 64) + " steal:" + strconv.FormatFloat(c[0].Steal, 'f', 1, 64) +
				"\nguest:" + strconv.FormatFloat(c[0].Guest, 'f', 1, 64) + " guestNice:" + strconv.FormatFloat(c[0].GuestNice, 'f', 1, 64) +
				" stolen:" + strconv.FormatFloat(c[0].Stolen, 'f', 1, 64)
			gc_cpu.Text = cpu_info
		} else {
			m, _ := mem.VirtualMemory()
			ga_mem.Percent = int(m.UsedPercent)
			cu, _ := cpu.Percent(1000000000, true)
			ga_cpu.Percent = int(cu[0])
			cc, _ := load.Avg()
			load1 := strconv.Itoa(int(cc.Load1))
			load5 := strconv.Itoa(int(cc.Load5))
			load15 := strconv.Itoa(int(cc.Load15))
			c, _ := cpu.Times(false)
			cpu_info := "Load1:" + load1 + "  load5:" + load5 + "  load15:" + load15 + "\n" +
				"cpu:" + c[0].CPU + " user:" + strconv.FormatFloat(c[0].User, 'f', 1, 64) + " system:" + strconv.FormatFloat(c[0].System, 'f', 1, 64) +
				"\nidle:" + strconv.FormatFloat(c[0].Idle, 'f', 1, 64) + " nice:" + strconv.FormatFloat(c[0].Nice, 'f', 1, 64) +
				"\niowait:" + strconv.FormatFloat(c[0].Iowait, 'f', 1, 64) + " irq:" + strconv.FormatFloat(c[0].Irq, 'f', 1, 64) +
				"\nsoftirq:" + strconv.FormatFloat(c[0].Softirq, 'f', 1, 64) + " steal:" + strconv.FormatFloat(c[0].Steal, 'f', 1, 64) +
				"\nguest:" + strconv.FormatFloat(c[0].Guest, 'f', 1, 64) + " guestNice:" + strconv.FormatFloat(c[0].GuestNice, 'f', 1, 64) +
				" stolen:" + strconv.FormatFloat(c[0].Stolen, 'f', 1, 64)
			gc_cpu.Text = cpu_info
		}

		termui.Render(ga_mem, ga_cpu)
	})

	termui.Loop()
}
