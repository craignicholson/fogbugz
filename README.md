# FogBugz Package for FogBugz API

Introduction
------------
API interface to the FogBugz xml api.  Used to collect data from the
FogBugz API and produce a similar report as ClarkKent.
http://help.fogcreek.com/8202/xml-api  

Installation and usage
----------------------

The import path for the packages are *github.com/craignicholson/fogbugz/fogbugz*.

To install it, run:

    go get github.com/craignicholson/fogbugz/fogbugz


Package Structure
----------------------
The fogbugz package contains a folder *fogbugz* which has one .go file in
sub folder.  The reason for this is case.go, interval.go, milestone.go, and
people.go each contain a **Response** struct and inside the **Response** struct
is the the struct for the data (case, interval, milestone, people).  

Placing each package in separate folders helps keep the naming schema the same
for each package for the **Response** struct.

Grouping all the packages into one folder allows the user to import
all packages with on *go get*.

Learning how to structure go code and packages is part of this applications
goal as well.  Learning by doing.  
* https://golang.org/doc/code.html
* https://peter.bourgon.org/go-in-production/

Using this package
----------------------
A simple CLI application was created to show how to use this package.

    github.com/craignicholson/fogbugzexporter


Additional Goals
----------------------
* Provide web UI for accountant to pull hours for each month by date range.
* Rollup hours by employee
* Save the raw data to .csv file which is available by Download

FogBugz API Documentation
----------------------

## Using the FogBugz API
API Link:         https://developers.fogbugz.com/default.asp?W194  Database Schema:  http://fogbugz.stackexchange.com/fogbugz-database-schema

### Get A Token

    https://company.fogbugz.com/api.asp?cmd=logon&email=scullyandmulder&password=trustnoone

```xml
<response>
  <token>
    <![CDATA[  [token]  ]]>
  </token>
</response>

<response>
  <error code="1">
  <![CDATA[ Incorrect password or username ]]>
  </error>
</response>

```

### Get Users
Gets list of People in FogBugz so we can walk back to the sFullName since the other API's return only the ixPerson.

    https://company.fogbugz.com/api.asp?token=d39eioefo0nkqf20adkan35qiokdi5&cmd=listPeople

```xml
<?xml version="1.0" encoding="UTF-8"?>
<response>
  <people>
    <person>
    <ixPerson>83</ixPerson>
      <sFullName>
      <![CDATA[ Person Abc ]]>
      </sFullName>
      <sEmail>
      <![CDATA[ hello@company.com ]]>
      </sEmail>
      <sPhone/>
      <fAdministrator>false</fAdministrator>
      <fCommunity>false</fCommunity>
      <fVirtual>false</fVirtual>
      <fDeleted>false</fDeleted>
      <fNotify>true</fNotify>
      <sHomepage/>
      <sLocale>
      <![CDATA[ * ]]>
      </sLocale>
      <sLanguage>
      <![CDATA[ * ]]>
      </sLanguage>
      <sTimeZoneKey>
      <![CDATA[ * ]]>
      </sTimeZoneKey>
      <sLDAPUid/>
      <dtLastActivity>2015-06-30T18:22:31Z</dtLastActivity>
      <fRecurseBugChildren>true</fRecurseBugChildren>
      <fPaletteExpanded>false</fPaletteExpanded>
      <ixBugWorkingOn>0</ixBugWorkingOn>
      <sFrom/>
    </person>
  <person>
    <ixPerson>66</ixPerson>
    <sFullName>
    <![CDATA[ A Person ]]>
    </sFullName>
    <sEmail>
    <![CDATA[ hello@company.com ]]>
    </sEmail>
    <sPhone>
    <![CDATA[ (111) 111-1111 ]]>
    </sPhone>
    <fAdministrator>true</fAdministrator>
    <fCommunity>false</fCommunity>
    <fVirtual>false</fVirtual>
    <fDeleted>false</fDeleted>
    <fNotify>true</fNotify>
    <sHomepage/>
    <sLocale>
    <![CDATA[ * ]]>
    </sLocale>
    <sLanguage>
    <![CDATA[ * ]]>
    </sLanguage>
    <sTimeZoneKey>
    <![CDATA[ * ]]>
    </sTimeZoneKey>
    <sLDAPUid/>
    <dtLastActivity>2016-02-10T23:00:16Z</dtLastActivity>
    <fRecurseBugChildren>false</fRecurseBugChildren>
    <fPaletteExpanded>true</fPaletteExpanded>
    <ixBugWorkingOn>0</ixBugWorkingOn>
    <sFrom/>
    </person>
  </people>
</response>
```

### Get Time Intervals
Get the time for a date range or the Case #s.

    https://company.fogbugz.com/api.asp?token=m8rr99dopu7hm2pmib5u5p4tpfd3ti&cmd=listIntervals&ixPerson=1&dtStart=2012-03-01&dtEnd=2012-03-28

ixPerson – optional
> Specifies which user’s intervals should be returned. If omitted, list intervals
  for the logged on user.  If set to 1, list intervals for all users. Note that
  you must be an administrator to see time interval information for users other
  than the logged on user.

> All parameters starting with the letters “dt” only accept times expressed in
  UTC (Coordinated Universal Time). Similarly, all return values starting with
  those letters will be expressed in UTC.

```xml
<response>
  <intervals>
    <interval>
      <ixInterval>4777</ixInterval>
      <ixPerson>28</ixPerson>
      <ixBug>1940</ixBug>
      <dtStart>2012-03-01T05:29:00Z</dtStart>
      <dtEnd>2012-03-01T05:59:00Z</dtEnd>
      <fDeleted>false</fDeleted>
      <sTitle>
      <![CDATA[Remove ssn and phone # from hacked tables]]>
      </sTitle>
    </interval>
    <interval>
      <ixInterval>4813</ixInterval>
      <ixPerson>28</ixPerson>
      <ixBug>2114</ixBug>
      <dtStart>2012-03-01T09:00:00Z</dtStart>
      <dtEnd>2012-03-01T09:30:00Z</dtEnd>
      <fDeleted>false</fDeleted>
      <sTitle>
      <![CDATA[ remove all data from table ABC ]]>
      </sTitle>
    </interval>
  </intervals>
</response>
```

### Get Case
Get the case Information We need for the time sheet.

> Additional data can be pulled by adding the element name to
the query string in case.go

    https://company.fogbugz.com/api.asp?token=iierr9nrc2vg441t1vcn8ee5ftlrqq&cmd=search&q=17308&cols=sTitle,ixProject,sProject,sArea,sCategory,sFixFor,ixFixFor,tags,plugin_customfields_at_fogcreek_com_customere1c

```xml
<response>
  <cases count="1">
      <case ixBug="3470" operations="edit,reopen,email,remind">
        <sTitle><![CDATA[ Permission error on Remittance computer ]]></sTitle>
        <ixProject>77</ixProject>
        <sProject><![CDATA[ IT - Company ABC ]]></sProject>
        <sArea><![CDATA[ Misc ]]></sArea>
        <sFixFor><![CDATA[ IT Services Support ]]></sFixFor>
        <ixFixFor>224</ixFixFor>
        <sCategory><![CDATA[ IT-Support ]]></sCategory>
        <tags>
          <tag>
          <![CDATA[ R-0041 ]]>
          </tag>
        </tags>
        <plugin_customfields_at_fogcreek_com_customere1c>
           <![CDATA[ Company ABC ]]>
        </plugin_customfields_at_fogcreek_com_customere1c>
      </case>
  </cases>
</response>
```

### Get Milestones
sFixFor is the milestone.
sStartNote is used by use internally so we need to pull these values.


    https://company.fogbugz.com/api.asp?token=[]&cmd=viewFixFor&ixFixFor=224

    https://company.fogbugz.com/api.asp?token=[]&cmd=listFixFors

```xml
<response>
  <fixfor>
    <ixFixFor>224</ixFixFor>
    <sFixFor><![CDATA[ IT Services Support ]]></sFixFor>
    <fInactive>false</fInactive>
    <dt>2013-12-31T06:00:00Z</dt>
    <ixProject>77</ixProject>
    <dtStart>2012-01-01T06:00:00Z</dtStart>
    <sStartNote><![CDATA[ Company ABC ]]></sStartNote>
    <setixFixForDependency/>
  </fixfor>
</response>

<response>
  <fixfors>
    <fixfor>
      <ixFixFor>1</ixFixFor>
      <sFixFor><![CDATA[ Undecided ]]></sFixFor>
      <fDeleted>false</fDeleted>
      <dt/>
      <dtStart/>
      <sStartNote/>
      <ixProject/>
      <sProject/>
      <setixFixForDependency/>
      <fReallyDeleted>false</fReallyDeleted>
    </fixfor>
    <fixfor>
      <ixFixFor>650</ixFixFor>
      <sFixFor><![CDATA[ Company EC ]]></sFixFor>
      <fDeleted>false</fDeleted>
      <dt/>
      <dtStart>2015-12-22T06:00:00Z</dtStart>
      <sStartNote><![CDATA[ Receipt of Customer Profile ]]></sStartNote>
      <ixProject>246</ixProject>
      <sProject><![CDATA[ BASE - software ]]></sProject>
      <setixFixForDependency/>
      <fReallyDeleted>false</fReallyDeleted>
    </fixfor>
  </fixfors>
<response>

```
