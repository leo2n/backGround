package main

import "log"

func main() {
	for i :=0; i<100000; i++ {
		log.Printf("Order: %d", i)
		_, err := downloadImg()
		if err!=nil {
			log.Printf("%d task is failed!", i)
		}
	}
}
