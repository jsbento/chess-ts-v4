package openingbook

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
	"sync"

	e "github.com/jsbento/chess-server-v4/internal/engine"
)

type TSVEntry struct {
	ECO   string
	Name  string
	PGN   string
	Moves []Move
}

type Move struct {
	WhiteMove string
	BlackMove string
}

type fileEntries struct {
	path    string
	entries []TSVEntry
}

func NewOpeningBook(engine *e.Engine) (*OpeningBook, error) {
	// parallel load all tsv files
	// non-docker data dir: ../../internal/engine/opening-book/data/
	tsvFiles := []string{
		"./internal/engine/opening-book/data/a.tsv",
		"./internal/engine/opening-book/data/b.tsv",
		"./internal/engine/opening-book/data/c.tsv",
		"./internal/engine/opening-book/data/d.tsv",
		"./internal/engine/opening-book/data/e.tsv",
	}

	entriesChan := make(chan fileEntries, len(tsvFiles))
	var wg sync.WaitGroup
	var loadErr error
	var loadErrMu sync.Mutex

	for _, filePath := range tsvFiles {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			entries, err := loadTSV(path)
			if err != nil {
				loadErrMu.Lock()
				if loadErr == nil {
					loadErr = err
				}
				loadErrMu.Unlock()
				return
			}
			entriesChan <- fileEntries{path: path, entries: entries}
		}(filePath)
	}

	go func() {
		wg.Wait()
		close(entriesChan)
	}()

	allEntries := []TSVEntry{}
	for fileData := range entriesChan {
		allEntries = append(allEntries, fileData.entries...)
	}

	if loadErr != nil {
		return nil, loadErr
	}

	// Build the opening book graph from all entries (sequential)
	book, err := BuildOpeningBookFromTSV(allEntries, engine)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func loadTSV(filePath string) ([]TSVEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := []TSVEntry{}
	for _, record := range records[1:] {
		moves, err := movesFromPGN(record[2])
		if err != nil {
			return nil, err
		}

		entries = append(entries, TSVEntry{
			ECO:   record[0],
			Name:  record[1],
			PGN:   record[2],
			Moves: moves,
		})
	}

	return entries, nil
}

func movesFromPGN(pgn string) ([]Move, error) {
	if pgn == "" {
		return []Move{}, errors.New("pgn is empty")
	}

	moves := []Move{}

	parts := strings.Split(pgn, " ")
	i := 0
	for i < len(parts) {
		blackMove := ""
		if i+2 < len(parts) {
			blackMove = parts[i+2]
		}
		moves = append(moves, Move{
			WhiteMove: parts[i+1],
			BlackMove: blackMove,
		})
		i += 3
	}

	return moves, nil
}
