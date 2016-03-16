// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package cases implements the interface(?) to the fogbugz api
// to pull cases using the search command.  To extend the data
// pulled add additional columns and remove the customer fields.
// cols:{"sTitle,ixProject,sProject,sArea,sCategory,sFixFor,ixFixFor,
// cols: sStartNote,tags,plugin_customfields_at_fogcreek_com_customere1c"}}
package cases

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
	Cases   Cases
}

// Cases holds the count of cases and the cases.
type Cases struct {
	XMLName xml.Name `xml:"cases"`
	Count   string   `xml:"count,attr"`
	Case    []Case   `xml:"case"`
}

// Case holds metadata about a case. Customer is a customer field and
// might need to be removed for other users of this software.
type Case struct {
	Bug        int    `xml:"ixBug,attr"`
	Operations string `xml:"operations,attr"`
	Title      string `xml:"sTitle"`
	ProjectID  int    `xml:"ixProject"`
	Project    string `xml:"sProject"`
	Area       string `xml:"sArea"`
	FixFor     string `xml:"sFixFor"`
	FixForID   int    `xml:"ixFixFor"`
	StartNote  string `xml:"sStartNote"`
	Category   string `xml:"sCategory"`
	Customer   string `xml:"plugin_customfields_at_fogcreek_com_customere1c"`
	Tags       Tags
}

// Tags holds all the tags assigned to a case.
type Tags struct {
	XMLName xml.Name `xml:"tags"`
	Tag     []string `xml:"tag"`
}

// ListCases takes in a comma seperated list of strings and returns a list
// of cases, cases can be comma separated list of case numbers without spaces.
// Using a map allows us to searches to find the Bug/Case by the ixBug
// For example, interval.go has a reference only to ixBug.
// https://company.fogbugz.com/api.asp?token=csjm9ljou4q4fkc3tgepgm2nm0vmfk&cmd=search&q=7880&cols=tags,sTitle,ixProject,sProject,sArea,sCategory,sFixFor,ixFixFor,sStartNote,plugin_customfields_at_fogcreek_com_customere1c,plugin_customfields_at_fogcreek_com_taskxordert0e
// The query term you are searching for.  Can be a string, a case number, a
// comma separated list of case numbers without spaces, e.g. 12,25,556.
// The last cols in the list of values is a custom field, this will need
// to be removed or maybe it will run if you chose not to edit...
func ListCases(token string, root *url.URL, casescsv string) map[int]Case {

	v := url.Values{"token": {token}, "cmd": {"search"}, "q": {casescsv}, "cols": {"sTitle,ixProject,sProject,sArea,sCategory,sFixFor,ixFixFor,sStartNote,tags,plugin_customfields_at_fogcreek_com_customere1c"}}
	resp, err := http.PostForm(root.String(), v)
	if err != nil {
		fmt.Println(err)
	}

	// Read the response into the data { []byte }.
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	// Load the xml data we collected into the object.
	r := &Response{}
	err3 := xml.Unmarshal(data, r)
	if err3 != nil {
		fmt.Println(err3)
	}

	// Create a map of cases to perform lookups
	// using the BugIDs in interval.go package.
	// Bug (ixBug) is the key in this map.
	c := make(map[int]Case)
	for i := 0; i < len(r.Cases.Case); i++ {
		c[r.Cases.Case[i].Bug] = r.Cases.Case[i]
	}

	return c
}
