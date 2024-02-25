package scrapper

import (
	"capAble-scrapper/pkg/constants"
	"capAble-scrapper/pkg/objects"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type NaukriJob struct {
	JobDetails []struct {
		Title                  string
		LogoPathV3             string
		FooterPlaceholderLabel string
		CompanyName            string
		TagsAndSkills          string
		JdURL                  string
		Placeholders           []struct {
			Type  string
			Label string
		}
		JobDescription string
		CreatedDate    int64
	}
}

var opportunityList []objects.Opportunity
var err error

func (scrapperSvc *Scrapper) NaukriScrape() ([]objects.Opportunity, error) {
	var noOfResults = 100
	var collector *colly.Collector

	url := fmt.Sprintf("%s%d%s", constants.Websites.Naukri, noOfResults, "&urlType=search_by_keyword&searchType=adv&keyword=disabled&pageNo=1&k=disabled&seoKey=disabled-jobs&src=jobsearchDesk&latLong=")
	if collector == nil {
		collector = colly.NewCollector(
			colly.Async(true),
		)

		collector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("authority", "www.naukri.com")
			r.Headers.Set("accept", "application/json")
			r.Headers.Set("accept-language", "en-US,en;q=0.8")
			r.Headers.Set("appid", "109")
			r.Headers.Set("content-type", "application/json")
			r.Headers.Set("referer", "https://www.naukri.com/disabled-jobs?k=disabled")
			r.Headers.Set("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Brave\";v=\"120\"")
			r.Headers.Set("sec-ch-ua-mobile", "?0")
			r.Headers.Set("sec-ch-ua-platform", "\"Linux\"")
			r.Headers.Set("sec-fetch-dest", "empty")
			r.Headers.Set("sec-fetch-site", "same-origin")
			r.Headers.Set("sec-gpc", "1")
			r.Headers.Set("systemid", "Naukri")
			r.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		})
	}

	collector.OnResponse(func(r *colly.Response) {
		err = parseNaukriJob(r.Body)

	})
	if err != nil {
		return nil, err
	}

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	err = collector.Visit(url)
	collector.Wait()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return opportunityList, nil
}

func parseNaukriJob(jobDatas []byte) error {
	opportunity := objects.Opportunity{}

	var parsedData = &NaukriJob{}

	err := json.Unmarshal(jobDatas, &parsedData)
	if err != nil {
		return err
	}

	for _, job := range parsedData.JobDetails {
		if err != nil {
			return err
		}

		opportunity.Company = job.CompanyName
		opportunity.Role = job.Title
		opportunity.JobPostingUrl = "https://www.naukri.com" + job.JdURL
		opportunity.CompanyLogo = job.LogoPathV3
		opportunity.Description = job.JobDescription

		unixTimestampMillis := job.CreatedDate
		unixTimestampSeconds := unixTimestampMillis / 1000

		opportunity.Posted = (time.Unix(unixTimestampSeconds, 0)).UTC()

		for i, value := range job.Placeholders {
			switch i {
			case 0:
				opportunity.Experience = value.Label
			case 1:
				opportunity.Pay = value.Label
			case 2:
				opportunity.Location = value.Label

			}
		}
		opportunityList = append(opportunityList, opportunity)

	}

	return err
}
