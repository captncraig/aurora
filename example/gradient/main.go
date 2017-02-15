package main

import (
	"flag"
	"fmt"
	"github.com/captncraig/aurora"
	"github.com/captncraig/aurora/example/detectAndAuth"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	flag.Parse()
	fixedMap := map[int]bool{}
	fixed := []*aurora.PanelColorCommand{}
	mixed := []*aurora.PanelColorCommand{}

	for _, a := range flag.Args() {
		i, err := strconv.Atoi(a)
		if err != nil {
			log.Fatalf("BAD ARG: %s. Must be number.", a)
		}
		fixedMap[i] = true
		fixed = append(fixed, &aurora.PanelColorCommand{
			ID: byte(i),
		})
	}

	c, err := detect.FindClient()
	if err != nil {
		log.Fatal(err)
	}
	info, err := c.GetInfo()
	if err != nil {
		log.Fatal(err)
	}

	panels := map[byte]*aurora.Panel{}
	for _, p := range info.Panels {
		panels[byte(p.ID)] = p
		if _, ok := fixedMap[p.ID]; !ok {
			mixed = append(mixed, &aurora.PanelColorCommand{
				ID: byte(p.ID),
			})
		}
	}

	//select initial colors
	targets := []*aurora.PanelColorCommand{}
	for i, f := range fixed {
		switch i % 3 {
		case 0:
			f.R = 255
		case 1:
			f.B = 255
		case 2:
			f.B = 12
		}
		targets = append(targets, &aurora.PanelColorCommand{
			R: byte(rand.Int()),
			G: byte(rand.Int()),
			B: byte(rand.Int()),
		})
		log.Println("FIXED:", f.ID)
	}

	for _, m := range mixed {
		log.Println("MIXED:", m.ID)
	}

	ch, err := c.ActivateExtranalControl()
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i, f := range fixed {
			adjustFixed(f, targets[i])
			ch <- f
		}
		for _, m := range mixed {
			adjustColors(m, panels, fixed)
			ch <- m
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func adjustFixed(f *aurora.PanelColorCommand, target *aurora.PanelColorCommand) {
	if f.R == target.R && f.G == target.G && f.B == target.B {
		target.B = byte(rand.Int())
		target.G = byte(rand.Int())
		target.R = byte(rand.Int())
		return
	}
	if f.R < target.R {
		f.R++
	}
	if f.R > target.R {
		f.R--
	}
	if f.G < target.G {
		f.G++
	}
	if f.G > target.G {
		f.G--
	}
	if f.B < target.B {
		f.B++
	}
	if f.B > target.B {
		f.B--
	}
}

func adjustColors(p *aurora.PanelColorCommand, panels map[byte]*aurora.Panel, fixed []*aurora.PanelColorCommand) {
	pan := panels[p.ID]
	fmt.Println(p.ID, "------------------")
	fmt.Println("MY: ", pan.X, pan.Y)
	dists := make([]float64, len(fixed))
	totalDist := float64(0)

	for i, f := range fixed {
		pan2 := panels[f.ID]
		dx := float64(pan2.X - pan.X)
		dy := float64(pan2.Y - pan.Y)
		dist := math.Sqrt(dx*dx + dy*dy)
		dists[i] = dist
		totalDist += dist
		fmt.Println(f.ID, dist)
	}

	var r, g, b float64

	for i, f := range fixed {
		weight := (totalDist - dists[i]) / totalDist
		fmt.Println(f.ID, "W:", weight)
		r += weight * float64(f.R)
		g += weight * float64(f.G)
		b += weight * float64(f.B)
	}

	fmt.Println("FINAL", r, g, b)
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}
	p.R = byte(r)
	p.G = byte(g)
	p.B = byte(b)
	return
}
