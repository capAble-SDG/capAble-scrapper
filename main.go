package main

import (
	"capAble-scrapper/api/pkg/firebase"
	"capAble-scrapper/api/pkg/scrapper"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {

	scrapperSvc := scrapper.NewService()
	firebaseSvc := firebase.NewService()

	opportunities, err := scrapperSvc.LinkedInScrape()

	if err != nil {
		if err != colly.ErrAlreadyVisited {
			log.Fatal(err)
		}
	}

	err = firebaseSvc.WriteToFirebase(opportunities)

	if err != nil {
		log.Fatal(err)
	}

	opportunities, err = scrapperSvc.NaukriScrape()

	if err != nil {
		fmt.Print(err)
	}

	err = firebaseSvc.WriteToFirebase(opportunities)

	if err != nil {
		log.Fatal(err)
	}

	opportunities, err = scrapperSvc.FoundItScrape()

	if err != nil {
		log.Fatal(err)
	}

	err = firebaseSvc.WriteToFirebase(opportunities)

	if err != nil {
		log.Fatal(err)
	}

}
