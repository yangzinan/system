// Copyright 2016 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import "github.com/gizak/termui"
import "github.com/shirou/gopsutil/mem"
import "github.com/shirou/gopsutil/cpu"
import "strconv"

func mem_UsedPercent(UsedPercent int) *termui.Gauge {
	g := termui.NewGauge()
	g.Percent = UsedPercent
	g.Width = 50
	g.Height = 3
	g.BorderLabel = "Memory UsedPercent"
	g.BarColor = termui.ColorRed
	g.BorderFg = termui.ColorWhite
	g.BorderLabelFg = termui.ColorCyan
	return g
}

func cpu_UsedPercent(UsedPercent int) *termui.Gauge {
	g := termui.NewGauge()
	g.Percent = UsedPercent
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
	mem_info := "Memtory Total: " + total + "MB\n" + "Memory Free:" + free + "MB\n"
	g := termui.NewPar(mem_info)
	g.Height = 5
	g.Width = 50
	g.Y = 3
	g.BorderLabel = "Memory Info"
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

	ga_mem := mem_UsedPercent(int(v.UsedPercent))
	ga_cpu := cpu_UsedPercent(int(c[0]))
	gc_mem := mem_info(v)

	termui.Render(ga_mem, gc_mem, ga_cpu)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)
		termui.SendCustomEvt("/usr/t", t.Count)

		if t.Count%2 == 0 {
			v, _ := mem.VirtualMemory()
			ga_mem.Percent = int(v.UsedPercent)
			c, _ := cpu.Percent(1000000000, true)
			ga_cpu.Percent = int(c[0])
		} else {
			v, _ := mem.VirtualMemory()
			ga_mem.Percent = int(v.UsedPercent)
			c, _ := cpu.Percent(1000000000, true)
			ga_cpu.Percent = int(c[0])
		}

		termui.Render(ga_mem, ga_cpu)

	})

	termui.Loop()
}
