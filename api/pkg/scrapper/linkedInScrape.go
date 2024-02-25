package scrapper

import (
	"capAble-scrapper/pkg/constants"
	"capAble-scrapper/pkg/objects"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var count int
var collector *colly.Collector

func (scrapperSvc *Scrapper) LinkedInScrape() ([]objects.Opportunity, error) {
	count = 0
	start := 0
	var opportunityList []objects.Opportunity
	var err error

	if collector == nil {
		collector = colly.NewCollector()

		err := collector.Limit(&colly.LimitRule{
			DomainGlob: "*",
			Delay:      2 * time.Second,
		})

		if err != nil {
			log.Fatal(err)
		}

		collector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("authority", "www.linkedin.com")
			r.Headers.Set("accept", "*/*")
			r.Headers.Set("accept-language", "en-US,en;q=0.5")
			r.Headers.Set("referer", "https://www.linkedin.com/jobs/search?trk=guest_homepage-basic_guest_nav_menu_jobs&position=1&pageNum=0")
			r.Headers.Set("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Brave\";v=\"120\"")
			r.Headers.Set("sec-ch-ua-mobile", "?0")
			r.Headers.Set("sec-ch-ua-platform", "\"Linux\"")
			r.Headers.Set("sec-fetch-dest", "empty")
			r.Headers.Set("sec-fetch-site", "same-origin")
			r.Headers.Set("sec-gpc", "1")
			r.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		})
	}

	collector.OnHTML(".base-card__full-link", func(e *colly.HTMLElement) {
		itemURL := e.Attr("href")
		opportunity := scrapeJobItem(itemURL)
		opportunityList = append(opportunityList, opportunity)
	})

	// Set up error handling
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println(err.Error())
	})

	for start < 50 {
		err = collector.Visit(constants.Websites.LinkedIn["MainPage"] + fmt.Sprint(start))
		for err != nil {
			err = collector.Visit(constants.Websites.LinkedIn["MainPage"] + fmt.Sprint(start))
			fmt.Println(err.Error())
		}
		fmt.Printf("___%d Done___", start+25)
		start += 25
	}

	return opportunityList, nil
}

func scrapeJobItem(itemUrl string) objects.Opportunity {
	opportunity := objects.Opportunity{}
	opportunity.JobPostingUrl = itemUrl
	count++
	fmt.Printf("Processing %d...\n", count)
	patternItemID := regexp.MustCompile(`\d+\?refId=\w+`)
	itemID := patternItemID.FindString(itemUrl)

	collector.OnHTML(".top-card-layout__title", func(e *colly.HTMLElement) {
		opportunity.Role = strings.TrimSpace(e.Text)
	})

	collector.OnHTML(".topcard__org-name-link", func(e *colly.HTMLElement) {
		opportunity.Company = strings.TrimSpace(e.Text)
	})

	collector.OnHTML(".topcard__flavor", func(e *colly.HTMLElement) {
		opportunity.Location = strings.TrimSpace(e.Text)
	})

	collector.OnHTML(".show-more-less-html__markup", func(e *colly.HTMLElement) {
		htmlDescription, err := e.DOM.Html()
		if err != nil {
			fmt.Print(err)
		}
		opportunity.Description = htmlDescription
	})

	collector.OnHTML(".description__job-criteria-list", func(e *colly.HTMLElement) {
		e.ForEach(".description__job-criteria-text", func(i int, h *colly.HTMLElement) {
			switch i {
			case 0:
				opportunity.SeniorityLevel = strings.TrimSpace(h.Text)
			case 1:
				opportunity.EmploymentType = strings.TrimSpace(h.Text)
			case 2:
				opportunity.Job = strings.TrimSpace(h.Text)
			case 3:
				opportunity.Industry = strings.TrimSpace(h.Text)

			}

		})
	})

	err = collector.Visit(constants.Websites.LinkedIn["SubPage"] + itemID)

	return opportunity
}
