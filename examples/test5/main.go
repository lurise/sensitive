package main

import (
	"github.com/lurise/sensitive"
	"log"
)

func main() {

	text := `你好傻逼呀，真的号傻逼要大姐夫`
	//textTemp := text
	filter := sensitive.New()
	if err := filter.LoadWordDict("./dict/dict.txt"); err != nil {
		log.Println(err.Error())
	}

	text = filter.Hightlight(text)
	text = text + "傻逼"
	text = filter.Hightlight(text)

	println("text=" + text)
}
