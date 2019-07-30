package Story

type Story struct {
	Title   string
	Story   []string
	Options []Option
}

type Option struct {
	Text string
	Arc  string
}

const StoryTemplate = 
`<h1>{{ .Title }}</h1>
<p>{{ range $sentence := .Story }} {{ $sentence }} {{ end }}</p>
<h2>Choices:</h2>
{{ range .Options }}
<a href="/arcs/{{ .Arc }}"> {{ .Text }} </a><br>
{{ end }}`
