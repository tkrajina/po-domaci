// Package main is generated with ftmpl {{{v0.3.1}}}, do not edit!!!! */
package main

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"os"
)

func init() {
	_ = fmt.Sprintf
	_ = errors.New
	_ = os.Stderr
	_ = html.EscapeString
}

// TMPLERRletter evaluates a template letter.tmpl
func TMPLERRletter(letter TemplateLetter) (string, error) {
	_template := "letter.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
`)
	_w(`<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
`)
	_w(`
<h1>`)
	_w(fmt.Sprintf(`%s`, _escape(letter.Letter)))
	_w(`</h1>

<ul>
`)
	for _, word := range letter.Words {
		_w(`    <li/> <a href="dictionary_word_`)
		_w(fmt.Sprintf(`%s`, _escape(word.Back)))
		_w(`.html">`)
		_w(fmt.Sprintf(`%s`, _escape(word.Back)))
		_w(`</a>
`)
	}
	_w(`</ul>
`)
	_w(`</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLletter evaluates a template letter.tmpl
func TMPLletter(letter TemplateLetter) string {
	html, err := TMPLERRletter(letter)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template letter.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRmain evaluates a template main.tmpl
func TMPLERRmain(template Template) (string, error) {
	_template := "main.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
`)
	_w(`<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
`)
	_w(`
bu

<ul>
`)
	for _, letter := range template.Letters {
		_w(`    <li/> <a href="dictionary_`)
		_w(fmt.Sprintf(`%s`, _escape(letter.Letter)))
		_w(`.html">`)
		_w(fmt.Sprintf(`%s`, _escape(letter.Letter)))
		_w(`</a> (`)
		_w(fmt.Sprintf(`%d`, len(letter.Words)))
		_w(` riječi)
`)
	}
	_w(`</ul>
`)
	_w(`</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLmain evaluates a template main.tmpl
func TMPLmain(template Template) string {
	html, err := TMPLERRmain(template)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template main.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRword evaluates a template word.tmpl
func TMPLERRword(word TemplateWord) (string, error) {
	_template := "word.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
`)
	_w(`<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
`)
	_w(`
Bu
`)
	_w(`</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLword evaluates a template word.tmpl
func TMPLword(word TemplateWord) string {
	html, err := TMPLERRword(word)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template word.tmpl:" + err.Error())
	}
	return html
}
