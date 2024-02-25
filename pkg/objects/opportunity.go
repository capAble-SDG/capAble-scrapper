package objects

import "time"

type Opportunity struct {
	Company        string
	CompanyLogo    string
	CompanyURL     string
	JobPostingUrl  string
	Location       string
	Industry       string
	Role           string
	Description    string
	Experience     string
	SeniorityLevel string
	EmploymentType string
	Job            string
	Pay            string
	Posted         time.Time
}
