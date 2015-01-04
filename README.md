gossip
======

Fetch &amp; store messages by room, author, and time.

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

