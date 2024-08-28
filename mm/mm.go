package mm

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func MoveMouse(seconds time.Duration) {
  fmt.Println("WIP...")
  i := 1

	for {
    boxes := strings.Repeat("■", i)
    fmt.Printf("\rProgress: %v%% %s", i, boxes)
		x, y := robotgo.Location()
		robotgo.MoveSmooth(x+5, y)
		x, y = robotgo.Location()
		robotgo.MoveSmooth(x-5, y)
		time.Sleep(seconds * time.Second)
    i++
	}
}
