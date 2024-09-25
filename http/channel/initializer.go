package channel

import (
	"fmt"
	"hotNews/http/controllers"
	"time"
)

func Init() {
	ch := make(chan string)

	go sendData(ch)
	go getData(ch)

	time.Sleep(1e9 * 10)
}

func sendData(ch chan string) {
	for {
		ch <- "ZhTop"
	}
}

func getData(ch chan string) {
	var input string
	for {
		input = <-ch
		switch input {
		case "ZhTop":
			controllers.ZhTop()
		default:
			fmt.Println("休整中。。。")
		}
		time.Sleep(1e9 * 10)
	}
}
