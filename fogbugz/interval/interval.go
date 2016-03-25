// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package interval implements the fetching of time intervals from the fogbugz api.
package interval

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Response contains the array of intervals.
type Response struct {
	XMLName   xml.Name `xml:"response"`
	Intervals Intervals
}

// Intervals contains the list of interval data and list of bugs and the
// distinct list of bugs passed to ListCases for all of the intervals collected.
type Intervals struct {
	XMLName  xml.Name   `xml:"intervals"`
	Interval []Interval `xml:"interval"`

	// Bugs: Distinct iXBug(s) in csv string passed to ListCases.
	Bugs string
}

// Interval holds the meta data about increment of time spent on a bug.
type Interval struct {
	IntervalID int       `xml:"ixInterval"`
	PersonID   int       `xml:"ixPerson"`
	Bug        int       `xml:"ixBug"`
	Start      time.Time `xml:"dtStart"` // UTC Date
	End        time.Time `xml:"dtEnd"`   // UTC Date
	Deleted    bool      `xml:"fDeleted"`
	Title      string    `xml:"sTitle"`

	// Empty values populated after the xml is unmarshalled
	StartLocal time.Time // Start Date converted to local time
	EndLocal   time.Time // End Date converted to local time
	Duration   float64   // Duration of the hours worked
}

// ListIntervals accepts start date and end date in UTC and returns a list of
// cases and converts the UTC dates to Local dates using the timezone, and
// calculates the duration of each time interval in hours.
func ListIntervals(token string, root *url.URL, dtStartUTC string, dtEndUTC string, timezone string) Intervals {
	// https://company.fogbugz.com/api.asp?token=[token]&cmd=listIntervals&ixPerson=1&dtStart=2012-03-01&dtEnd=2012-03-28
	// ixPerson = 1, will collect data for all persons during the timeframe.

	location, err4 := time.LoadLocation(timezone)
	if err4 != nil {
		fmt.Printf("LoadLocation : %s", err4)
	}

	v := url.Values{"token": {token}, "cmd": {"listIntervals"}, "ixPerson": {"1"}, "dtStart": {dtStartUTC}, "dtEnd": {dtEndUTC}}
	resp, errresp := http.PostForm(root.String(), v)
	if errresp != nil {
		fmt.Printf("PostForm : %s", errresp)
	}

	// Read the reponse data.
	defer resp.Body.Close()
	data, errbody := ioutil.ReadAll(resp.Body)
	if errbody != nil {
		fmt.Println(errbody)
		fmt.Printf("ioutil.ReadAll : %s", errbody)
	}

	// Load the data we collected into the PeopleResponse.
	r := &Response{}
	err := xml.Unmarshal(data, r)
	if err != nil {
		fmt.Printf("ListIntervals - Unmarshal : %s", err)
	}

	// Collect the unique set bugs for the set of intervals collected.
	// The struct will remain empty.  An empty struct uses 0 bytes.
	bugs := make(map[int]struct{})

	// Calculate duration and format the time for local dates.
	// https://golang.org/pkg/time/#Location
	for i := 0; i < len(r.Intervals.Interval); i++ {
		// Using the location we can avoid DST Issues letting time.In()
		// convert the dates from UTC to local time.
		// TODO: should I even do this here... and just do it in the api method instead
		// Because only the exported data and the UI care about the time in Local time
		r.Intervals.Interval[i].StartLocal = r.Intervals.Interval[i].Start.In(location)
		r.Intervals.Interval[i].EndLocal = r.Intervals.Interval[i].End.In(location)

		// TODO: DURATION BELONGS HERE... So DOES THE TOTAL Distinct BUGS
		// Duration Calulcation using the time.Sub(), and we report in hours.
		// https://golang.org/pkg/time/#Time.Sub
		duration := r.Intervals.Interval[i].End.Sub(r.Intervals.Interval[i].Start)
		r.Intervals.Interval[i].Duration = duration.Hours()
		bugs[r.Intervals.Interval[i].Bug] = struct{}{}
	}

	// Create a csv string of bugs to pull header data for in ListCases.
	var bug bytes.Buffer
	for key := range bugs {
		bug.WriteString(strconv.Itoa(key) + ",")
	}
	//TODO: Remove the last trailing comma, even though the api works with comma.

	r.Intervals.Bugs = bug.String()
	return r.Intervals
}
