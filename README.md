gossip
======

Fetch &amp; store messages by room, author, and time.

[![Build Status](https://travis-ci.org/parkr/gossip.svg?branch=master)](https://travis-ci.org/parkr/gossip)

## Running

### Docker

To run this with Docker, simply:

```console
~$ hostip=$(ip route show 0.0.0.0/0 | grep -Eo 'via \S+' | awk '{ print $2 }')
~$ docker run -it --rm \
    --publish 8080:8080 \
    --add-host=mysql:$hostip \
    --env GOSSIP_DB_HOSTNAME=mysql \
    --env GOSSIP_DB_USERNAME=mysqlusername \
    --env GOSSIP_DB_PASSWORD=mysqlpassword \
    --env GOSSIP_DB_DBNAME=mysqldbname \
    --env GOSSIP_AUTH_TOKEN=authtokentovalidateclients \
    parkr/gossip \
    gossip -bind=:8080
```

## API

### Storing messages

Send a `POST` request to `/api/messages/log` with the following data as a URL-encoded string:

```json
{
    "room": "#jekyll",
    "author": "parkr",
    "message": "hey y'all",
    "time": "Mon, 02 Jan 2006 15:04:05 MST"
}
```

### Fetching messages

Send a `GET` request to `/api/messages/log`. You can optionally add a `limit=N` to the querystring to limit results. `limit` defaults to `10`.

```json
{
    "ok": "true",
    "values": [ ... ]
}
```

## Server Configuration

Some environment variables are required to connect for proper functionality:

- `GOSSIP_DB_HOSTNAME` (optional)
- `GOSSIP_DB_USERNAME`
- `GOSSIP_DB_PASSWORD`
- `GOSSIP_DB_DBNAME`
- `GOSSIP_AUTH_TOKEN` (used to authenticate api requests from the client)

Optionally, set the `GOSSIP_LOGFILE` when running `script/deploy`.

## Credits / License

Copyright (c) 2014 Parker Moore (@parkr)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

