package constants

var Websites = struct {
	MyDisabilityJobs string
	SimplyHired      string
	LinkedIn         map[string]string
	Indeed           string
	FoundIt          map[string]string
	Shine            string
	FreshersWorld    string
	Naukri           string
}{
	MyDisabilityJobs: "https://mydisabilityjobs.com/search-page/",
	SimplyHired:      "www.simplyhired.co.in",

	LinkedIn: map[string]string{
		"MainPage": "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?keywords=pwd&location=India&start=",
		"SubPage":  "https://www.linkedin.com/jobs-guest/jobs/api/jobPosting/",
	},

	Indeed:        "https://in.indeed.com/",
	Shine:         "https://www.shine.com/job-search",
	FreshersWorld: "https://www.freshersworld.com/jobs/jobsearch",
	FoundIt: map[string]string{
		"MainPage": "https://www.foundit.in/middleware/jobsearch?sort=1&limit=",
		"SubPage":  "https://www.foundit.in/middleware/jobdetail/",
	},
	Naukri: "https://www.naukri.com/jobapi/v3/search?noOfResults=",
}
