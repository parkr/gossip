package template

import (
	htmltemplate "html/template"

	"github.com/parkr/gossip/database"
)

type ShowTemplateData struct {
	PriorMessages      []database.Message
	Message            database.Message
	SubsequentMessages []database.Message
	Rooms              []string
	CurrentRoom        string
}

func init() {
	ShowTemplate = htmltemplate.Must(htmltemplate.New("showTemplate").Funcs(addtlFuncs).Parse(showTemplate))
}

var ShowTemplate *htmltemplate.Template
var showTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Chat Log Server / Message {{.Message.ID}} in {{.Message.Room}}</title>
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
        <p><input type="search" name="q" value=""></p>
        <p><input type="submit" value="search"></p>
      </form>
    </div>
    <article class="context pre">
      {{range .PriorMessages}}
      <article class="message-item">
        <a class="link" href="/messages/{{.ID}}/context">+</a>
        <span class="time-container">[<span class="time">{{.CreatedAtRFC3339}}</span>]</span>
        <span class="author"><a  href="/messages/by/{{.Author}}">{{.Author}}</a></span>
        <span class="message">{{.Message}}</span>
        <br />
      </article>
      {{end}}
    </article>
    <article class="main-message">
      {{with .Message}}
      <h1 id="message-id">message no. {{.ID}}</h1>
      <p>
        Posted by <a  href="/messages/by/{{.Author}}">{{.Author}}</a> in <a  href="/room/{{encodeRoom .Room}}">{{.Room}}</a> at {{.CreatedAtRFC3339}}
      </p>
      <div class="monospaced">
          {{.Message}}
      </div>
      {{end}}
    </article>
    <article class="context post">
      {{range .SubsequentMessages}}
      <article class="message-item">
        <a class="link" href="/messages/{{.ID}}/context">+</a>
        <span class="time-container">[<span class="time">{{.CreatedAtRFC3339}}</span>]</span>
        <span class="author"><a  href="/messages/by/{{.Author}}">{{.Author}}</a></span>
        <span class="message">{{.Message}}</span>
        <br />
      </article>
      {{end}}
    </article>
  </div>
</body>
</html>
`
