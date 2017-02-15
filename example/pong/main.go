package main

import (
	"github.com/captncraig/aurora"
	"github.com/captncraig/aurora/example/detectAndAuth"
	"log"
	"math/rand"
	"time"
)

var panels = []byte{248, 83, 222, 249, 186, 77, 224, 91, 140}

func main() {
	c, err := detect.FindClient()
	if err != nil {
		log.Fatal(err)
	}
	ch, err := c.ActivateExtranalControl()
	if err != nil {
		log.Fatal(err)
	}
	//current := 0
	blank := &aurora.PanelColorCommand{
		R: 0,
		G: 0,
		B: 0,
	}

	for _, p := range panels {
		blank.ID = p
		ch <- blank
		time.Sleep(time.Microsecond)
	}
	time.Sleep(500 * time.Millisecond)
	//current := 0
	//dx := 1
	//color.ID = panels[0]
	//ch <- color
	for {
		modifyColor()
		time.Sleep(101 * time.Millisecond)

		for _, p := range panels {
			color.ID = p
			ch <- color
			time.Sleep(time.Microsecond)
		}
		// blank.ID = panels[current]
		// ch <- blank
		// if (current == 0 && dx == -1) || (current == len(panels)-1 && dx == 1) {
		// 	dx = -dx
		// }
		// current += dx
		// color.ID = panels[current]
		// ch <- color
	}
}

var color = &aurora.PanelColorCommand{
	R: 255,
}

var curTarget = rand.Int()

func modifyColor() {
	tarR := byte((curTarget & 0xff0000) >> 16)
	tarG := byte((curTarget & 0xff00) >> 8)
	tarB := byte(curTarget & 0xff)
	//log.Println(tarR, tarG, tarB, color.R, color.G, color.B)
	if tarR == color.R && tarG == color.G && tarB == color.B {
		curTarget = rand.Int()
		return
	}
	if color.R > tarR {
		color.R--
	}
	if color.R < tarR {
		color.R++
	}
	if color.G > tarG {
		color.G--
	}
	if color.G < tarG {
		color.G++
	}
	if color.B > tarB {
		color.B--
	}
	if color.B < tarB {
		color.B++
	}
}
