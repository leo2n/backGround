package main

import "time"

func main() {
	for i :=0; i<500; i++ {
		downloadImg()
		time.Sleep(30*time.Second)
	}
}
