package scrapper

import (
	"capAble-scrapper/pkg/constants"
	"capAble-scrapper/pkg/objects"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var opportunity objects.Opportunity
var sub_item = false

type FoundItJob struct {
	JobSearchResponse struct {
		Data []struct {
			JobId           int32
			Title           string
			Locations       string
			CreatedAt       int64
			Industries      []string
			EmploymentTypes []string
			CompanyName     string
			Exp             string
			Salary          string
			SeoJdUrl        string
			SeoCompanyUrl   string
		}
	}
}

type FoundItJobDetail struct {
	JobDetailResponse struct {
		Description string
	}
}

func (scrapperSvc *Scrapper) FoundItScrape() ([]objects.Opportunity, error) {

	limit := 100
	var opportunityList []objects.Opportunity

	url := fmt.Sprintf("%s%d%s", constants.Websites.FoundIt["MainPage"], limit, "&query=disabled")
	if collector == nil {
		collector = colly.NewCollector()

		collector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("authority", "wwwF.foundit.in")
			r.Headers.Set("accept", "application/json, text/plain, */*")
			r.Headers.Set("accept-language", "en-US,en;q=0.6")
			r.Headers.Set("cache-control", "no-cache")
			r.Headers.Set("pragma", "no-cache")
			r.Headers.Set("referer", "https://www.foundit.in")
			r.Headers.Set("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Brave\";v=\"120\"")
			r.Headers.Set("sec-ch-ua-mobile", "?0")
			r.Headers.Set("sec-ch-ua-platform", "\"Linux\"")
			r.Headers.Set("sec-fetch-dest", "empty")
			r.Headers.Set("sec-fetch-site", "same-origin")
			r.Headers.Set("sec-gpc", "1")
			r.Headers.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		})
	}

	collector.OnResponse(func(r *colly.Response) {
		if !sub_item {
			err = parseFoundItJob(r.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			opportunityList, err = parseFoundItJobDescription(r.Body, opportunityList)
			if err != nil {
				fmt.Println(err)
			}
		}

	})

	if err != nil {
		return nil, err
	}

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	err = collector.Visit(url)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Length Final: ", len(opportunityList))
	return opportunityList, nil
}

func parseFoundItJob(jobDatas []byte) error {

	var parsedData = &FoundItJob{}

	err := json.Unmarshal(jobDatas, &parsedData)
	if err != nil {
		return err
	}

	for _, job := range parsedData.JobSearchResponse.Data {
		opportunity = objects.Opportunity{}

		opportunity.Company = job.CompanyName
		opportunity.CompanyURL = "https://www.foundit.in" + job.SeoCompanyUrl
		opportunity.JobPostingUrl = "https://www.foundit.in" + job.SeoJdUrl
		opportunity.Location = job.Locations

		opportunity.Role = job.Title
		unixTimestampMillis := job.CreatedAt
		unixTimestampSeconds := unixTimestampMillis / 1000

		opportunity.Posted = (time.Unix(unixTimestampSeconds, 0)).UTC()
		opportunity.Experience = job.Exp

		if !strings.Contains(job.Salary, "0-0") {
			opportunity.Pay = job.Salary
		}

		opportunity.Industry = strings.Join(job.Industries, ",")
		opportunity.EmploymentType = strings.Join(job.EmploymentTypes, ",")

		sub_item = true
		collector.Visit(constants.Websites.FoundIt["SubPage"] + fmt.Sprint(job.JobId))
	}

	return err
}

func parseFoundItJobDescription(jobDatas []byte, opportunityList []objects.Opportunity) ([]objects.Opportunity, error) {
	var parsedData = &FoundItJobDetail{}
	err := json.Unmarshal(jobDatas, &parsedData)
	if err != nil {
		return opportunityList, nil
	}

	opportunity.Description = parsedData.JobDetailResponse.Description
	sub_item = false
	return append(opportunityList, opportunity), nil
}
