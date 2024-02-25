package scrapper

import "capAble-scrapper/pkg/objects"

type Scrapper struct {
}

type ScrapperInterface interface {
	LinkedInScrape() ([]objects.Opportunity, error)
	NaukriScrape() ([]objects.Opportunity, error)
	FoundItScrape() ([]objects.Opportunity, error)
}

func NewService() ScrapperInterface {
	return &Scrapper{}
}
