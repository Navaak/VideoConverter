package main

import "navaak/convertor/app"

func main() {
	a, _ := app.New(app.DefaultConfig)
	a.Run()
}
