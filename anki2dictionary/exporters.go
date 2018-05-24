package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/tkrajina/go-reflector/reflector"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var sorter = collate.New(language.Croatian)

func init() {
	str := ` * **boška**: šuma `
	for _, r := range str {
		if r == 'š' {
			fmt.Println(r, "!!!")
		} else {
			fmt.Println(r, string(r), "jok")
		}
	}
	str = strings.Replace(str, "š", "___", -1)
	str = strings.Replace(str, "š", "___", -1)
	//fmt.Println(croatianLatexChars(str))
	//os.Exit(1)
}

var (
	markdownTmpl = template.Must(template.ParseFiles("templates/markdown.tmpl"))
	htmlTmpl     = template.Must(template.ParseFiles("templates/html.tmpl"))
	latexTmpl    = template.Must(template.ParseFiles("templates/latex.tmpl"))
)

func exportToJson(dict Dictionary) (map[string]string, error) {
	return map[string]string{
		"dictionary.json": jsonizePrettified(dict),
	}, nil
}

func exportWithTemplate(filename string, dict Dictionary, tmpl *template.Template) (map[string]string, error) {
	var res bytes.Buffer
	if err := tmpl.Execute(&res, dict.ToTemplateParams()); err != nil {
		return nil, err
	}

	return map[string]string{filename: res.String()}, nil
}

func exportToMarkdown(dict Dictionary) (map[string]string, error) {
	return exportWithTemplate("dictionary.md", dict, markdownTmpl)
}

func exportToHtml(dict Dictionary) (map[string]string, error) {
	params := dict.ToTemplateParams()
	files := map[string]string{
		"dictionary.html": TMPLmain(params),
	}

	for _, letter := range params.Letters {
		files[fmt.Sprintf("dictionary_%s.html", letter.Letter)] = TMPLletter(letter)
		words := make([]TemplateWord, 0, len(letter.Words))
		for _, word := range letter.Words {
			words = append(words, word)
			sort.Sort(TemplateWords(words))
		}
		for _, word := range words {
			files[fmt.Sprintf("dictionary_word_%s.html", word.Back)] = TMPLword(word)
		}
	}

	return files, nil
}

func exportToLatex(dict Dictionary) (map[string]string, error) {
	tmpl := dict.ToTemplateParams()
	for n := range tmpl.Letters {
		tmpl.Letters[n].Letter = croatianLatexChars(tmpl.Letters[n].Letter)
		for m := range tmpl.Letters[n].Words {
			word := tmpl.Letters[n].Words[m]
			obj := reflector.New(&word)
			for _, field := range obj.FieldsAll() {
				val, err := field.Get()
				if err != nil {
					return nil, err
				}
				if s, is := val.(string); is {
					if err := field.Set(croatianLatexChars(s)); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	var res bytes.Buffer
	if err := latexTmpl.Execute(&res, tmpl); err != nil {
		return nil, err
	}

	return map[string]string{"dictionary.tex": res.String()}, nil
}

func croatianLatexChars(s string) string {
	s = strings.Replace(s, "š", `\v{s}`, -1)
	s = strings.Replace(s, "Š", `\v{S}`, -1)
	s = strings.Replace(s, "Đ", `\DJ{}`, -1)
	s = strings.Replace(s, "đ", `\dj{}`, -1)
	s = strings.Replace(s, "č", `\v{c}`, -1)
	s = strings.Replace(s, "Č", `\v{C}`, -1)
	s = strings.Replace(s, "ć", `\'{c}`, -1)
	s = strings.Replace(s, "Ć", `\'{C}`, -1)
	s = strings.Replace(s, "ž", `\v{z}`, -1)
	s = strings.Replace(s, "Ž", `\v{Z}`, -1)

	// It's not s but s +  ̌
	s = strings.Replace(s, "š", `\v{s}`, -1)
	s = strings.Replace(s, "Š", `\v{S}`, -1)
	s = strings.Replace(s, "ž", `\v{z}`, -1)
	s = strings.Replace(s, "Ž", `\v{Z}`, -1)
	s = strings.Replace(s, "č", `\v{c}`, -1)
	s = strings.Replace(s, "Č", `\v{C}`, -1)
	s = strings.Replace(s, "ć", `\'{c}`, -1)
	s = strings.Replace(s, "Ć", `\'{C}`, -1)

	//fmt.Println(s)
	return s
}
