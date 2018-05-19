package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"os"

	"github.com/tkrajina/anki"
)

func usageAndExitf(msg string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", params...)
	os.Exit(1)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		usageAndExitf(err.Error())
	}

	fmt.Println(cfg.DatabaseFilename)

	db, err := anki.OpenOriginalDB(cfg.DatabaseFilename)
	if err != nil {
		usageAndExitf(err.Error())
	}
	defer func() {
		fmt.Println("Closing db")
		db.Close()
	}()

	fmt.Println(db)
	dict, err := loadDictionary(db, *cfg)
	if err != nil {
		usageAndExitf(err.Error())
	}

	dict.SortColumn = 1
	dict.Sort()

	fmt.Printf("found %d words", len(dict.Rows))

	exporters := []func(Dictionary) (map[string]string, error){
		exportToJson,
		exportToMarkdown,
		exportToHtml,
		exportToLatex,
	}

	for _, fnc := range exporters {
		out, err := fnc(*dict)
		if err != nil {
			usageAndExitf(err.Error())
		}

		for filename, contents := range out {
			if err := ioutil.WriteFile(fmt.Sprintf("output/%s", filename), []byte(contents), os.FileMode(0770)); err != nil {
				usageAndExitf(err.Error())
			}

			fmt.Printf("Saved to %s\n", filename)
		}
	}
}

func loadDictionary(db *anki.DB, cfg Config) (*Dictionary, error) {

	collection, err := db.Collection()
	if err != nil {
		return nil, err
	}

	for modelId, model := range collection.Models {
		fmt.Printf("model [%d] %s deck=%d\n", modelId, model.Name, model.DeckID)
	}

	for deckId, deck := range collection.Decks {
		fmt.Printf("deck [%d/%d] %s\n", deckId, deck.ID, deck.Name)
	}

	notesById := map[anki.ID]anki.Note{}
	notes, err := db.Notes()
	if err != nil {
		return nil, err
	}
	for notes.Next() {
		note, err := notes.Note()
		if err != nil {
			return nil, err
		}
		notesById[note.ID] = *note
	}
	notes.Close()

	allCards := []anki.Card{}

	cards, err := db.Cards()
	if err != nil {
		return nil, err
	}
	for cards.Next() {
		card, err := cards.Card()
		if err != nil {
			return nil, err
		}
		allCards = append(allCards, *card)
	}
	cards.Close()

	var dictionary Dictionary
	for _, card := range allCards {
		deck, found := collection.Decks[card.DeckID]
		//fmt.Println("found", deck.Name, deckName)
		if found && deck.Name == cfg.DeckName {
			note, found := notesById[card.NoteID]
			if !found {
				fmt.Printf("Note %d not found\n", card.NoteID)
				continue
			}

			model, found := collection.Models[note.ModelID]
			if !found {
				fmt.Println("Model not found", note.ModelID, note.FieldValues)
				continue
			}

			//fmt.Println("*", model.Name, cardType)
			if model.Name == cfg.Type {
				processNote(&dictionary, note, *model)
			} else {
				fmt.Printf("Note %#v in deck %s but not of type %s\n", note.FieldValues, cfg.DeckName, cfg.Type)
			}
		}
	}

	/*
		bytes, err := json.MarshalIndent(dictionary, "", "    ")
		if err != nil {
			return nil, err
		}

		fmt.Println(string(bytes))
	*/

	return &dictionary, nil
}

func processNote(dictionary *Dictionary, note anki.Note, model anki.Model) {
	fields := make([]string, len(model.Fields))
	values := make([]string, len(model.Fields))
	for n, f := range model.Fields {
		fields[n] = f.Name
		values[n] = html.UnescapeString(string(note.FieldValues[n]))
	}

	dictionary.Columns = fields
	dictionary.Rows = append(dictionary.Rows, values)
}
