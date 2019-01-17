package main

import (
	"bytes"
	"text/template"
	"time"
)

var (
	headerTemplate   = &templateValue{name: "header"}
	footerTemplate   = &templateValue{name: "footer"}
	renderingContext = &renderingContextHolder{
		Command:  command,
		Interval: interval,
	}
	timeFormat = app.Flag("timeFormat", "How to format the time. See: https://golang.org/pkg/time/").
			Default("2006-01-02 15:04:05").
			String()
)

func init() {
	app.Flag("header", "Will print a header what will be executed. If empty no header will be displayed. See: https://golang.org/pkg/text/template/").
		Short('h').
		Default("[{{.Now}}] Execute [{{.Command}}] every {{.Interval}}\n").
		SetValue(headerTemplate)
	app.Flag("footer", "Will print a footer what was executed. If empty no footer will be displayed. See: https://golang.org/pkg/text/template/").
		Short('t').
		Default("[{{.Now}}] {{.Command.ResultSummary}}\n").
		SetValue(footerTemplate)
}

type renderingContextHolder struct {
	Command  *commandHolder
	Interval *time.Duration
}

func (instance *renderingContextHolder) Now() string {
	return time.Now().Format(*timeFormat)
}

func printHeader() {
	printHighlightedToTerminal(render(headerTemplate))
}

func printFooter() {
	printHighlightedToTerminal(render(footerTemplate))
}

func render(tmpl *templateValue) string {
	buf := new(bytes.Buffer)
	err := tmpl.tmpl.Execute(buf, renderingContext)
	if err != nil {
		fatal("Cannot render %s: %v", tmpl.name, err)
	}
	return buf.String()
}

type templateValue struct {
	name   string
	tmpl   *template.Template
	source string
}

func (instance *templateValue) String() string {
	return instance.source
}

func (instance *templateValue) Get() interface{} {
	return instance.tmpl
}

func (instance *templateValue) Set(source string) error {
	tmpl, err := template.New(instance.name).Parse(source)
	if err != nil {
		return err
	}
	instance.tmpl = tmpl
	instance.source = source
	return nil
}
