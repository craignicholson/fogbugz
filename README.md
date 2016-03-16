# gofogbugz
Application to collect data from the FogBugz API - ClarkKent expansion
http://help.fogcreek.com/8202/xml-api

## TODOs:
Sometimes FogBugz fields have empty strings for dates and time.Time throughs
parse error, it would be good to see which interval, case etc... or maybe
just ignore these and make sure the data is still being captured.


## Overview
We use FogBugz to track projects for statements of work for clients instead
of the project being a software project.  This is atypical and the main reason
is for reporting hours used for each project.

Goals
* Provide web UI for accountant to pull hours for each month by date range.
* Rollup hours by employee
* Save the raw data to .csv file which is available by Download

## Using FogBugz to Report Hours

### How to quickly run this script
Edit the following values and click run, execute etc... in the sql tool (IDE)

@FirstDayofBillingMonth - SET THIS TO THE First Day of the month you are billing
@FirstDayofNextMonth    - SET THIS TO THE First Day of the next month

## Detail Overview

Users in FogBugz can use the working task to enter hours for each case.  FogBugz additionally has
a feature one can set for each user's Working Schedule under a user's own profile which allows
them to allocate their time and completion dates and set working days and a daily schedule which
can automatically start and stop work.  

Ocassionally people using this automatically start and stop work feature let the tasks run-on
without review of their hours and can over bill customers.

Also, when working outside of the workday schedule the user is prompted to answer 'yes/no' to
working outside of the workday, and if they say yes, and never stop the timer the timer
will run to midnight.  This will cause issues with the actual vs reported total and each user
will need to correct their own time.

They can sort this out using the Timesheet report under their own profile. (I think, a non admin
user needs to verify this)

## Troubleshooting Issues

## Using the FogBugz API
API Link:                       https://developers.fogbugz.com/default.asp?W194
Database Schema:                http://fogbugz.stackexchange.com/fogbugz-database-schema

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

The falling data pulls are necessary

### USERS
Get of LIST OF People in FogBugz so we can walk back to the sFullName

EXAMPLE
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

### TIME INTERVALS

GET All the time for a Month and the Case #s
We will need a way to manually edit or via a form request the dtStart and dtEnd
-----------------------------------------------------------------------------------------
https://company.fogbugz.com/api.asp?token=m8rr99dopu7hm2pmib5u5p4tpfd3ti&cmd=listIntervals&ixPerson=1&dtStart=2012-03-01&dtEnd=2012-03-28
https://company.fogbugz.com/api.asp?token=iierr9nrc2vg441t1vcn8ee5ftlrqq&cmd=listTags
cmd=listTags

ixPerson – optional
Specifies which user’s intervals should be returned. If omitted, list intervals for the logged on user.
If set to 1, list intervals for all users. Note that you must be an administrator to see time interval
information for users other than the logged on user.

All parameters starting with the letters “dt” only accept times expressed in UTC (Coordinated Universal Time). Similarly, all return values starting with those letters will be expressed in UTC.


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

### Case INFORMATION
Fetch the case Information We need for the time sheet for Ann
-----------------------------------------------------------------------------------------

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

### MILESTONES
------FixFor - sStartNote is not showing up yet.... sFixFor is the milestone.. to go fetch the sStartNote
We will need to pull a list fo the ixForFor Ids and request them all out.

cmd=viewFixFor&ixFixFor=224
https://company.fogbugz.com/api.asp?token=[]&cmd=viewFixFor&ixFixFor=224
https://compnay.fogbugz.com/api.asp?token=[]&cmd=listFixFors

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
