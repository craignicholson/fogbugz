// Copyright 2015 Craig Nicholson. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package people implements some I/O utility functions.
package people

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Response holds the People list.
type Response struct {
	XMLName xml.Name `xml:"response"`
	People  People
}

// People holds the list of persons.
type People struct {
	XMLName xml.Name `xml:"people"`
	Person  []Person `xml:"person"`
}

// Person holds the users metadata in the FogBugz Application.
type Person struct {
	PersonID           int       `xml:"ixPerson"`
	FullName           string    `xml:"sFullName"`
	Email              string    `xml:"sEmail"`
	Phone              string    `xml:"sPhone"`
	Administrator      bool      `xml:"fAdministrator"`
	Community          bool      `xml:"fCommunity"`
	Virtual            bool      `xml:"fVirtual"`
	Deleted            bool      `xml:"fDeleted"`
	Notify             bool      `xml:"fNotify"`
	Homepage           string    `xml:"sHomepage"`
	Locale             string    `xml:"sLocale"`
	Language           string    `xml:"sLanguage"`
	TimeZoneKey        string    `xml:"sTimeZoneKey"`
	LDAPUid            string    `xml:"sLDAPUid"`
	LastActivity       time.Time `xml:"dtLastActivity"` // UTC Date
	RecurseBugChildren bool      `xml:"fRecurseBugChildren"`
	PaletteExpanded    bool      `xml:"fPaletteExpanded"`
	BugWorkingOnID     int       `xml:"ixBugWorkingOn"`
	From               string    `xml:"sFrom"`
}

// ListPeople lists all of the users in FogBugz and return a map of Person.
// Using a map allows us to searches to find the person by the ixPersonID
// For example, interval.go has a reference only to ixPersonID.
func ListPeople(token string, root *url.URL) map[int]Person {
	// http://company.fogbugz.com/api.asp?token=[token]&cmd=listPeople
	// fIncludeDeleted=0
	// fIncludeActive=1 – default 1
	v := url.Values{"token": {token}, "cmd": {"listPeople"}, "fIncludeDeleted": {"1"}, "fIncludeActive": {"1"}}
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

	// Aa map Person(s) is used to perform lookups
	// on the data using the ixPersonID in interval.go.
	// PersonID is the key in this map.
	p := make(map[int]Person)
	for i := 0; i < len(r.People.Person); i++ {
		p[r.People.Person[i].PersonID] = r.People.Person[i]
	}

	return p
}
