// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package fogbugz implements the structs and methods used
// to combine data for reports to output data for an application.
// this package contains the hour struct for the ClarkKent report
package fogbugz

import (
	"time"
)

// Hour holds metadata about the hours for each person.
// In FogBugz this report is similar to the ClarkKent report with more fields.
type Hour struct {
	ID            int
	StartDate     time.Time // time period begins in local time
	EndDate       time.Time // time period ends in local time
	Title         string    // BugID & Case Title {17137 : 2016 - PTO - Holiday - Me}
	Duration      float64   // duration of the time in hours
	Expense       string    // leave blank
	Employee      string    // person name
	Project       string    // project name
	MileStone     string    // milestone name -> sFixFor
	Customer      string    // customer name
	CaseNumber    int       // case number or bug number
	BillingPeriod string    // billing period is 1stCheck or 2ndCheck any period after 15th is 2ndCheck
	Area          string    // Area is accounting defined area
	Category      string    // category refers to accounting defined category... i have no clue what they are doing here anymore
	StartNote     string    // startnote was the task order # at one time, but now the project is the task order
	Year          int       // year as integer for pivot tables
	Month         string    // month as integer for pivot tables
	Day           int       // day as integer for pivot tables
	DOW           string    // Day of the week, to help review the hours per day
	Tags          string    // Tags on each case, used by EV and IP
}

// HourByEmployee will hold the total hours for a date range.
type HourByEmployee struct {
	PersonID   int
	Employee   string    // person name
	TotalHours float64   // total time for the duration
	StartDate  time.Time // time period begins in local time
	EndDate    time.Time // time period ends in local time
}
