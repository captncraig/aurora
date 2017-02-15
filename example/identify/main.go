package main

import (
	"fmt"

	"github.com/captncraig/aurora"
	"log"
	"time"
)

func main() {
	c := aurora.NewWithToken("http://10.0.1.50:16021", "f5TTjV33uPFn9B6tMAmgR9q3N5aOTmT7")
	dat, err := c.GetInfo()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range dat.Panels {
		fmt.Printf("Panel #%d @ %d,%d (%d)\n", p.ID, p.X, p.Y, p.Rotation)
	}
	ch, err := c.ActivateExtranalControl()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(dat.Panels); i++ {
		for j, p := range dat.Panels {
			pt := &aurora.PanelColorCommand{
				ID: byte(p.ID),
				R:  byte(0),
				G:  byte(0),
				B:  byte(0),
			}
			if i == j {
				pt.R = 255
				fmt.Println("!!!!!!! ", pt.ID)
			}
			ch <- pt
		}
		time.Sleep(2 * time.Second)
	}
}
