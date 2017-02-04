package main

import (
	"fmt"

	"github.com/captncraig/aurora"
	"log"
	"math/rand"
	"time"
)

func main() {
	c := aurora.NewWithToken("http://10.0.1.50:16021", "yVui4pFJYFISPxWtOMqKiBcopG9SGOHC")
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
	for {
		for _, p := range dat.Panels {
			pt := &aurora.PanelColorCommand{
				ID: byte(p.ID),
				R:  byte(rand.Int()),
				G:  byte(rand.Int()),
				B:  byte(rand.Int()),
			}
			ch <- pt
		}
		time.Sleep(time.Second)
	}
}
