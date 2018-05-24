package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/tkrajina/go-reflector/reflector"
)

type Config struct {
	DatabaseFilename string `json:"database_filename"`
	DeckName         string `json:"deck_name"`
	Type             string `json:"type"`
}

var audioRegexp = regexp.MustCompile(`\[sound:.*?\]`)

func valueAndAudioFile(val string) (string, string) {
	var audioFile string
	val = audioRegexp.ReplaceAllStringFunc(val, func(s string) string {
		s = strings.Replace(strings.Trim(s, "[]"), "sound:", "", -1)
		audioFile = s
		return ""
	})
	val = strings.TrimSpace(val)
	return val, audioFile
}

// Dictionary is used to store all the data in a pretty printed JSON file
type Dictionary struct {
	Columns    []string   `json:"columns"`
	SortColumn int        `json:"sort_column"`
	Rows       [][]string `json:"rows"`
}

func (d Dictionary) ToTemplateParams() Template {
	var res Template

	var firstLtr rune
	for _, row := range d.Rows {
		word := row[d.SortColumn]
		//fmt.Println(word)
		l, _ := firstLetter(word)
		if strings.ToUpper(string(l)) != strings.ToUpper(string(firstLtr)) {
			res.Letters = append(res.Letters, TemplateLetter{
				Letter: strings.ToUpper(string(l)),
				Words:  []TemplateWord{},
			})
		}

		rowMap := d.RowToMap(row)
		last := &res.Letters[len(res.Letters)-1]
		last.Words = append(last.Words, rowMap)

		firstLtr = l
	}

	return res
}

func (d *Dictionary) RowToMap(row []string) TemplateWord {
	var rowMap TemplateWord
	obj := reflector.New(&rowMap)
	for _, field := range obj.FieldsAll() {
		fld, err := field.Tag("field")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read tag for field "+field.Name())
		}
		for columnNo, column := range d.Columns {
			if column == fld {
				val, audioFile := valueAndAudioFile(row[columnNo])
				if audioFile != "" {
					if rowMap.AudioFiles == nil {
						rowMap.AudioFiles = make([]string, 0, 1)
					}
					rowMap.AudioFiles = append(rowMap.AudioFiles, audioFile)
				}
				if err := field.Set(val); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Error setting field %s", err.Error()))
				}
			}
		}
	}
	return rowMap
}

func (d *Dictionary) Len() int      { return len(d.Rows) }
func (d *Dictionary) Swap(i, j int) { d.Rows[i], d.Rows[j] = d.Rows[j], d.Rows[i] }
func (d *Dictionary) Less(i, j int) bool {
	a := ignoreNonLetters(d.Rows[i][d.SortColumn])
	b := ignoreNonLetters(d.Rows[j][d.SortColumn])
	return sorter.CompareString(a, b) < 0
}
func (d *Dictionary) Sort() { sort.Sort(d) }

type Template struct {
	Letters []TemplateLetter
}

type TemplateLetter struct {
	Letter string
	Words  []TemplateWord
}

type TemplateWord struct {
	Front     string `field:"Front"`
	Back      string `field:"Back"`
	Primjer   string `field:"Primjer"`
	Varijante string `field:"Varijante"`
	Sinonimi  string `field:"Sinonimi"`
	Vezano    string `field:"Vezano`

	AudioFiles []string `field:"AudioFiles"`
}

func (tw TemplateWord) sortField() string { return tw.Back }

type TemplateWords []TemplateWord

func (tw TemplateWords) Len() int { return len(tw) }
func (tw TemplateWords) Less(i, j int) bool {
	a := tw[i].sortField()
	b := tw[j].sortField()
	return sorter.CompareString(a, b) < 0
}
func (tw TemplateWords) Swap(i, j int) { tw[i], tw[j] = tw[j], tw[i] }
