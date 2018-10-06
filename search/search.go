package search

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	SearchEngine *TitleSearcher
)

type TitleSearcher struct {
	searcher  *riot.Engine
	titleAids map[string]int64
	reexp     *regexp.Regexp
}

func InitEngineWithFiles(filename string) *TitleSearcher {
	titleSearcher := &TitleSearcher{searcher: &riot.Engine{}, titleAids: make(map[string]int64), reexp: regexp.MustCompile("[0-9]+")}

	titleSearcher.searcher.Init(types.EngineOpts{
		Using:     3,
		NumShards: 2,
		// IDOnly:        true,
		GseDict: "dictionary.txt",
	})

	words := titleSearcher.LoadWords(filename)
	fmt.Println("Debug len:", len(words))
	for idx, word := range words {
		titleSearcher.searcher.Index(uint64(idx+1), types.DocData{Content: word})
	}

	// 等待索引刷新完毕
	titleSearcher.searcher.Flush()

	return titleSearcher
}

func (titleSearcher *TitleSearcher) LoadWords(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		handleErr(err)
		return nil
	}
	defer file.Close()

	wordsToTest := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		res := strings.Split(text, ",")
		if len(res) < 2 {
			continue
		}

		aid, err := strconv.ParseInt(strings.TrimSpace(res[len(res)-1]), 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		title := strings.TrimSpace(res[len(res)-2])
		titleSearcher.titleAids[title] = aid

		wordsToTest = append(wordsToTest, title)
	}

	return wordsToTest
}

func (titleSearcher *TitleSearcher) GetAidByTitle(title string) int64 {
	aid, ok := titleSearcher.titleAids[strings.TrimSpace(title)]
	if !ok {
		return -1
	}
	return aid
}

func (titleSearcher *TitleSearcher) GetExactEpisode(queryText string) string {
	queryText = strings.Replace(queryText, "\"", "", -1)
	queryText = strings.TrimSpace(queryText)

	res := titleSearcher.reexp.FindAllString(queryText, -1)
	if len(res) > 0 {
		var buffer bytes.Buffer
		for _, v := range queryText {
			if unicode.IsDigit(v) {
				break
			}
			buffer.WriteRune(v)
		}
		for _, v := range res {
			buffer.WriteString(" " + v)
		}
		queryText = buffer.String()
	}

	output := titleSearcher.searcher.Search(types.SearchReq{Text: queryText})
	log.Println("Debug engine result:", queryText, output)
	if output.NumDocs == 0 {
		return ""
	}
	switch v := output.Docs.(type) {
	case types.ScoredDocs:
		return v[0].Content
	default:
		log.Printf("Default: %T\n", v)
	}
	return ""
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
