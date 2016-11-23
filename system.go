// Copyright 2016 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import "github.com/gizak/termui"
import "github.com/shirou/gopsutil/mem"
import "github.com/shirou/gopsutil/cpu"

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	//termui.UseTheme("helloworld")

	ga_mem := termui.NewGauge()
	v, _ := mem.VirtualMemory()
	ga_mem.Percent = int(v.UsedPercent)
	ga_mem.Width = 50
	ga_mem.Height = 3
	ga_mem.BorderLabel = "Memory UsedPercent"
	ga_mem.BarColor = termui.ColorRed
	ga_mem.BorderFg = termui.ColorWhite
	ga_mem.BorderLabelFg = termui.ColorCyan

	ga_cpu := termui.NewGauge()
	c, _ := cpu.Percent(1000000000, true)
	ga_cpu.Percent = int(c[0])
	ga_cpu.Width = 50
	ga_cpu.Height = 3
	ga_cpu.PercentColor = termui.ColorBlue
	ga_cpu.X = 51
	ga_cpu.BorderLabel = "CPU UsedPercent"
	ga_cpu.BarColor = termui.ColorYellow
	ga_cpu.BorderFg = termui.ColorWhite

	gc_mem := termui.NewPar("Simple colored text\nwith label. It [can be](fg-red) multilined with \\n or [break automatically](fg-red,fg-bold)")
	gc_mem.Height = 5
	gc_mem.Width = 50
	gc_mem.Y = 3
	gc_mem.BorderLabel = "Memory Info"
	gc_mem.BorderFg = termui.ColorYellow

	g3 := termui.NewGauge()
	g3.Percent = 50
	g3.Width = 50
	g3.Height = 3
	g3.Y = 11
	g3.BorderLabel = "Gauge with custom label"
	g3.Label = "{{percent}}% (100MBs free)"
	g3.LabelAlign = termui.AlignRight

	g4 := termui.NewGauge()
	g4.Percent = 50
	g4.Width = 50
	g4.Height = 3
	g4.Y = 14
	g4.BorderLabel = "Gauge"
	g4.Label = "Gauge with custom highlighted label"
	g4.PercentColor = termui.ColorYellow
	g4.BarColor = termui.ColorGreen
	g4.PercentColorHighlighted = termui.ColorBlack

	termui.Render(ga_mem, gc_mem, ga_cpu, g3, g4)

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
