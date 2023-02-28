package watch

import (
	"bytes"
	"github.com/alecthomas/kingpin"
	"text/template"
	"time"
)

func newRendering(
	command func() *command,
	interval func() time.Duration,
) *rendering {
	return &rendering{
		headerTemplate: &templateValue{name: "header"},
		footerTemplate: &templateValue{name: "footer"},
		renderingContext: &renderingContextHolder{
			command:    command,
			interval:   interval,
			timeFormat: "2006-01-02 15:04:05",
		},
	}
}

type rendering struct {
	headerTemplate   *templateValue
	footerTemplate   *templateValue
	renderingContext *renderingContextHolder
}

func (this *rendering) ConfigureCli(cli *kingpin.Application) {
	cli.Flag("header", "Will print a header what will be executed. If empty no header will be displayed. See: https://golang.org/pkg/text/template/").
		Short('h').
		Default("[{{.Now}}] Execute [{{.Command}}] every {{.Interval}}\n").
		SetValue(this.headerTemplate)
	cli.Flag("footer", "Will print a footer what was executed. If empty no footer will be displayed. See: https://golang.org/pkg/text/template/").
		Short('t').
		Default("[{{.Now}}] {{.Command.ResultSummary}}\n").
		SetValue(this.footerTemplate)

	this.renderingContext.ConfigureCli(cli)
}

type renderingContextHolder struct {
	command    func() *command
	interval   func() time.Duration
	timeFormat string
}

func (this *renderingContextHolder) ConfigureCli(cli *kingpin.Application) {
	cli.Flag("timeFormat", "How to format the time. See: https://golang.org/pkg/time/").
		Default("2006-01-02 15:04:05").
		StringVar(&this.timeFormat)
}

func (this *renderingContextHolder) Command() *command {
	return this.command()
}

func (this *renderingContextHolder) Interval() time.Duration {
	return this.interval()
}

func (this *renderingContextHolder) Now() string {
	return time.Now().Format(this.timeFormat)
}

func (this *rendering) render(tmpl *templateValue) string {
	buf := new(bytes.Buffer)
	err := tmpl.tmpl.Execute(buf, this.renderingContext)
	if err != nil {
		fatal("Cannot render %s: %v", tmpl.name, err)
	}
	return buf.String()
}

func (this *rendering) renderHeader() string {
	return this.render(this.headerTemplate)
}

func (this *rendering) renderFooter() string {
	return this.render(this.footerTemplate)
}

type templateValue struct {
	name   string
	tmpl   *template.Template
	source string
}

func (this *templateValue) String() string {
	return this.source
}

func (this *templateValue) Get() interface{} {
	return this.tmpl
}

func (this *templateValue) Set(source string) error {
	tmpl, err := template.New(this.name).Parse(source)
	if err != nil {
		return err
	}
	this.tmpl = tmpl
	this.source = source
	return nil
}
