# asterisk-dialer [![Build Status](https://travis-ci.org/incu6us/asterisk-dialer.svg?branch=master)](https://travis-ci.org/incu6us/asterisk-dialer) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/b90f76e81907427e9f2d0a4124a34028)](https://www.codacy.com/app/incu6us/asterisk-dialer?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=incu6us/asterisk-dialer&amp;utm_campaign=Badge_Grade)

### api

 * `GET: /api/v1/dialer/start` - start dialer 
 * `GET: /api/v1/dialer/stop` - stop dialer
 * `GET: /api/v1/dialer/status` - get status for dialer (is it started or stopped)
 * `GET: /api/v1/dialer/registeredUsers` - get registered users, which could accept a call
 * `GET: /api/v1/dialer/msisdn/all` - get all msisdns
 * `GET: /api/v1/dialer/msisdn/inProgress?limit=1&page=1` - get msisdns which are int progress(limit & page are optionals)
 * `PUT: /api/v1/dialer/msisdn/inProgress/{id}` - update msisdn priority
   
   Body example:
   ```
   {"priority": 1}
   ```
 * `DELETE: /api/v1/dialer/msisdn/inProgress/{id}` - delete msisdn record
 * `POST: /api/v1/dialer/upload-msisdn` - api for list uploading 
 
    Example for uploading:
    ```
    0631234567;
    0931234567;
    ```
 * `GET: /api/v1/ready` - app check


MSISDNs older then 7 days will be removed