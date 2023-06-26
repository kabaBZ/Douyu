package main

import (
	"github.com/kabaBZ/Barrage_Go/src/barrageCrawler"
)

func main() {
	roomID := "11144156"
	dyBarrageCrawler := barrageCrawler.NewDyDanmuCrawler(roomID)
	dyBarrageCrawler.Start()
}
