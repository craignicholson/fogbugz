// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package milestone implements the transformation of FixFor data into a usable form.
// The data is very unpredicaable with empty strings and nulls for various
// elements and omitempty does not seem to help solve this issue when the values
// in the element needs to be parsed (int, bool, time.Time).  Because of this
// unresolved issue all the FixFor are strings until I can see if the issue is
// my code or is a bug.
package milestone

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Response contains the array of intervals.
type Response struct {
	XMLName xml.Name `xml:"response"`
	FixFors FixFors
}

// FixFors holds the count of cases and the cases.
type FixFors struct {
	XMLName xml.Name `xml:"fixfors"`
	FixFor  []FixFor `xml:"fixfor"`
}

// FixFor holds the metadata about a Milestone.
type FixFor struct {
	FixForID         int    `xml:"ixFixFor,omitempty"`
	FixFor           string `xml:"sFixFor"`
	Deleted          string `xml:"fDeleted"`
	Created          string `xml:"dt"`         // UTC Date using string these dates are empty
	StartDate        string `xml:"dtStart"`    // UTC Date using string
	StartNote        string `xml:"sStartNote"` // this is the one field we need
	ProjectID        string `xml:"ixProject,omitempty"`
	Project          string `xml:"sProject"`
	FixForDependency string `xml:"setixFixForDependency"`
	ReallyDeleted    string `xml:"fReallyDeleted,omitempty"`
}

// ListMilestone collects the milestones (FixFors)
func ListMilestone(token string, root *url.URL) map[int]FixFor {
	// https://company.fogbugz.com/api.asp?token=h5ed5c72inoa46cpa5jj503ccu676j&cmd=listFixFors
	v := url.Values{"token": {token}, "cmd": {"listFixFors"}}
	resp, err := http.PostForm(root.String(), v)
	if err != nil {
		fmt.Println(err)
	}

	// Read the reponse data.
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	// Load the data we collected.
	r := &Response{}
	err3 := xml.Unmarshal(data, r)
	if err3 != nil {
		fmt.Println(err3)
	}

	// Create a map of FixFors to perform lookups
	// using the ixFixFor (FixForID) in interval.go package.
	// Bug (ixBug) is the key in this map.
	f := make(map[int]FixFor)
	for i := 0; i < len(r.FixFors.FixFor); i++ {
		f[r.FixFors.FixFor[i].FixForID] = r.FixFors.FixFor[i]
	}

	return f
}
