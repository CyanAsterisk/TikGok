package main

import (
	sociality "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"log"
)

func main() {
	svr := sociality.NewServer(new(SocialityServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
