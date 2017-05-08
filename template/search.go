package template

import (
	htmltemplate "html/template"

	"github.com/parkr/gossip/database"
)

type SearchTemplateData struct {
	Results map[string][]database.Message
	Total   int
	Rooms   []string
	Query   string
}

func init() {
	SearchTemplate = htmltemplate.Must(htmltemplate.New("searchTemplate").Funcs(addtlFuncs).Parse(searchTemplate))
}

var SearchTemplate *htmltemplate.Template
var searchTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Chat Log Server / Search for {{.Query}}</title>
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
  <link rel="stylesheet" href="/assets/normalize.css">
  <link rel="stylesheet" href="/assets/main.css">
  <script src="/assets/moment.min.js" type="text/javascript" charset="utf-8"></script>
  <script src="/assets/messages.js" type="text/javascript" charset="utf-8"></script>
</head>
<body>
  <div id="container">
    <div id="header">
      <h1 id="irc_logs">irc logs</h1>

      {{range .Rooms}}
      <h2 id="{{.}}"><a  href="/room/{{encodeRoom .}}">{{.}}</a></h2>
      {{end}}

      <form action="/search" method="get" accept-charset="utf-8">
        <p><input type="search" name="q" value="{{.Query}}"></p>
        <p><input type="submit" value="search"></p>
      </form>
    </div>
    <h1 id="latest_messages">{{.Total}} messages for query '{{.Query}}'</h1>
    <code class="messages-list">
        {{range $room, $messages := .Results}}
        <h2 id="messages_in_{{$room}}">in <a  href="/room/{{encodeRoom $room}}">{{$room}}</a></h2>
            {{range $messages}}
            <article class="message-item">
              <a class="link" href="/messages/{{.ID}}/context">+</a>
              <span class="time-container">[<span class="time">{{.CreatedAtRFC3339}}</span>]</span>
              <span class="author"><a  href="/messages/by/{{.Author}}">{{.Author}}</a></span>
              <span class="message">{{.Message}}</span>
              <br />
            </article>
            {{end}}
        {{end}}
    </code>
  </div>
</body>
</html>
`
