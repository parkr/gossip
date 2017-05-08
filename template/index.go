package template

import (
	htmltemplate "html/template"

	"github.com/parkr/gossip/database"
)

type ListTemplateData struct {
	Messages      []database.Message
	Rooms         []string
	CurrentRoom   string
	CurrentAuthor string
}

func init() {
	ListTemplate = htmltemplate.Must(htmltemplate.New("listTemplate").Funcs(addtlFuncs).Parse(listTemplate))
}

var ListTemplate *htmltemplate.Template
var listTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Chat Log Server / {{if .CurrentAuthor}}{{.CurrentAuthor}}{{else}}{{.CurrentRoom}}{{end}}</title>
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
    <h1 id="latest_messages">latest {{len .Messages}} messages {{if .CurrentAuthor}}by {{.CurrentAuthor}}{{else}}in {{.CurrentRoom}}{{end}}</h1>
<code class="messages-list">

  {{range .Messages}}
  <article class="message-item">
    <a class="link" href="/messages/{{.ID}}/context">+</a>
    <span class="time-container">[<span class="time">{{.CreatedAtRFC3339}}</span>]</span>
    <span class="author"><a  href="/messages/by/{{.Author}}">{{.Author}}</a></span>
    <span class="message">{{.Message}}</span>
    <br />
  </article>
  {{end}}

</code>

  </div>
</body>
</html>
`
