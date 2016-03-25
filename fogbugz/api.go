// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package fogbugz (api.go) implements the interface for calling the methods
// to fetch the data from the fogbugz api and handles the authentication.
package fogbugz

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	// Am I doing this correctly here, with code strucutre and formatting
	"github.com/craignicholson/fogbugz/fogbugz/case"
	"github.com/craignicholson/fogbugz/fogbugz/interval"
	"github.com/craignicholson/fogbugz/fogbugz/milestone"
	"github.com/craignicholson/fogbugz/fogbugz/people"
)

// API encapsultaes token, site url, and site version returned from FogBugz.
type API struct {
	token   string   // security token
	Root    *url.URL // fogbugz api URL
	version string   // API verions
}

// logonResult stores the result of the logon.
type logonResult struct {
	XMLName xml.Name `xml:"response"`
	Error   string   `xml:"error"`
	Token   string   `xml:"token"`
}

// Login to the FogBugz API and collect the token so we can send commands (cmd).
func (api *API) Login(email string, password string) error {
	v := url.Values{"cmd": {"logon"}, "email": {email}, "password": {password}}
	fmt.Println(v)
	resp, err := http.PostForm(api.Root.String(), v)
	if err != nil {
		fmt.Printf("Login PostForm failed: %s\nvalues%v\n", err, v)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Login ReadAll failed: %s\n", err)
	}

	// Emit error to site when logon fails.
	r := &logonResult{}
	if err := xml.Unmarshal(body, &r); err != nil {
		fmt.Printf("Login Unmarshal failed: %s\n", err)
	}

	api.token = r.Token
	return nil
}

// InvalidateToken Logoff from the FogBugz API and retire the token.
func (api *API) InvalidateToken() error {
	v := url.Values{"cmd": {"logoff"}, "token": {api.token}}
	resp, err := http.PostForm(api.Root.String(), v)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// GetCase pulls all the case header information which is seen at the
// top of fogbugz cases.
func (api *API) GetCase(casescsv string) map[int]cases.Case {
	list := cases.ListCases(api.token, api.Root, casescsv)
	return list
}

// GetInterval pulls time intervals for all people for a date range in UTC
// so we can roll up the durations and bill our customers.
func (api *API) GetInterval(startDateLocal time.Time, endDateLocal time.Time, timezone string) interval.Intervals {
	// LoadLocation uses http://golang.org/pkg/time/#LoadLocation.
	// icann names - get alist of these... we can use... for documentation
	// https://www.iana.org/time-zones
	location, err4 := time.LoadLocation("UTC")
	if err4 != nil {
		fmt.Printf("LoadLocation : %s", err4)
	}
	fmt.Println(startDateLocal)
	fmt.Println(endDateLocal)

	startDateUTC := startDateLocal.In(location)
	endDateUTC := endDateLocal.In(location)
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	//Print the dates out for debugging
	fmt.Println(startDateUTC.Format(longForm))
	fmt.Println(endDateUTC.Format(longForm))

	// Send the Dates in as RFC3339,
	list := interval.ListIntervals(api.token,
		api.Root,
		startDateUTC.Format(time.RFC3339),
		endDateUTC.Format(time.RFC3339),
		timezone)

	return list
}

// GetMilestone pulls all milestones available.
func (api *API) GetMilestone() map[int]milestone.FixFor {
	m := milestone.ListMilestone(api.token, api.Root)
	return m
}

// GetPeople returns a map of Persons and is to be used
// to lookup the Fullname and email address of each person.
func (api *API) GetPeople() map[int]people.Person {
	list := people.ListPeople(api.token, api.Root)
	return list
}

// GetHours in a accountant like format.
func (api *API) GetHours(startDateLocal string, endDateLocal string, timezone string) []Hour {
	// Convert the strings to time.Time
	loc, _ := time.LoadLocation(timezone)
	const shortFormlayout = "2006-01-02"
	startDate, _ := time.ParseInLocation(shortFormlayout, startDateLocal, loc)
	endDate, _ := time.ParseInLocation(shortFormlayout, endDateLocal, loc)

	//
	people := api.GetPeople()
	intervals := api.GetInterval(startDate, endDate, timezone)
	bugs := api.GetCase(intervals.Bugs)
	//misspelled, arg!
	milstones := api.GetMilestone()

	var hours []Hour
	// The indexes on PersonID (ixPerson) is larger than the actual
	// number of people so we have to make sure the max ixPerson exists
	// to size the map
	//hoursbyEmp := make([]HourByEmployee, 150)
	var hoursbyEmp map[int]HourByEmployee
	for i := 0; i < len(intervals.Interval); i++ {

		data := Hour{}
		data.ID = intervals.Interval[i].PersonID
		data.StartDate = intervals.Interval[i].StartLocal
		data.EndDate = intervals.Interval[i].EndLocal
		data.Title = bugs[intervals.Interval[i].Bug].Title
		// Truncate the float to 2 decimal places
		data.Duration = float64(int(intervals.Interval[i].Duration*100)) / 100
		data.Expense = ""
		data.Employee = people[intervals.Interval[i].PersonID].FullName
		data.Project = bugs[intervals.Interval[i].Bug].Project
		data.MileStone = bugs[intervals.Interval[i].Bug].FixFor
		data.Customer = bugs[intervals.Interval[i].Bug].Customer
		data.CaseNumber = intervals.Interval[i].Bug

		//Set the predefined BillingPeriod
		//TODO: Specific to my needs - need to rework
		data.BillingPeriod = "1stCheck"
		if intervals.Interval[i].EndLocal.Day() > 15 {
			data.BillingPeriod = "2ndCheck"
		}

		data.Area = bugs[intervals.Interval[i].Bug].Area
		data.Category = bugs[intervals.Interval[i].Bug].Category

		//Look up the start note - very inefficient ...
		//data.StartNote = api.GetMilestone(bugs[intervals.Interval[i].Bug].FixForID)
		data.StartNote = milstones[bugs[intervals.Interval[i].Bug].FixForID].StartNote

		data.Year = intervals.Interval[i].EndLocal.Year()
		data.Month = intervals.Interval[i].EndLocal.Month().String()
		data.Day = intervals.Interval[i].EndLocal.Day()
		data.DOW = intervals.Interval[i].EndLocal.Weekday().String()

		// For display and export to csv, write out all the tags as csv string
		var sTags string
		tags := bugs[intervals.Interval[i].Bug].Tags
		for i := 0; i < len(tags.Tag); i++ {
			sTags = sTags + tags.Tag[i] + ","
		}
		data.Tags = strings.TrimSuffix(sTags, ",")

		// roll up the hours for each employee using a map
		// OR maybe we have the count of people just make an
		t := hoursbyEmp[intervals.Interval[i].PersonID]
		t.PersonID = intervals.Interval[i].PersonID
		t.Employee = people[intervals.Interval[i].PersonID].FullName
		t.StartDate = intervals.Interval[i].StartLocal
		t.EndDate = intervals.Interval[i].EndLocal
		t.TotalHours = t.TotalHours + intervals.Interval[i].Duration

		hours = append(hours, data)
	}
	//Map the Employee total hours to the correct Employee
	//fmt.Println(hoursbyEmp)
	return hours
}

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
