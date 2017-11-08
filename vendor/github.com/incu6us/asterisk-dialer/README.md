# asterisk-dialer [![Build Status](https://travis-ci.org/incu6us/asterisk-dialer.svg?branch=master)](https://travis-ci.org/incu6us/asterisk-dialer) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/d1e83421924045a6a6bdaba87588a881)](https://www.codacy.com/app/incu6us/asterisk-dialer?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=incu6us/asterisk-dialer&amp;utm_campaign=Badge_Grade)

### api

 * `/api/v1/dialer/start` - start dialer 
 * `/api/v1/dialer/stop` - stop dialer
 * `/api/v1/dialer/status` - get status for dialer (is it started or stopped)
 * `/api/v1/dialer/registeredUsers` - get registered users, which could accept a call
 * `/api/v1/ready` - app check
 * `/api/v1/dialer/upload-msisdn` - api for list uploading (POST-method) 

Example for uploading:
```
0631234567;
0931234567;
```

MSISDNs older then 7 days will be removed