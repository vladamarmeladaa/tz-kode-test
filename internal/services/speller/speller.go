package speller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Speller struct {
	BaseURL string
}

type SpellerResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func New(source string) *Speller {
	return &Speller{
		BaseURL: source,
	}
}

func (s *Speller) Validate(texts []string) error {
	// const op = "internal/services/speller"

	params := url.Values{}
	for i := 0; i < len(texts); i++ {
		params.Add("text", texts[i])
	}

	requestURL := fmt.Sprintf("%s?%s", s.BaseURL, params.Encode())

	resp, err := http.Get(requestURL)
	if err != nil {
		return fmt.Errorf("error sending the request: %v", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading the request: %v", err)
	}

	sliceSpellerResult := [][]SpellerResult{}
	err = json.Unmarshal(bytes, &sliceSpellerResult)
	if err != nil {
		return fmt.Errorf("error decoding speller result: %v", err)
	}

	spellerError := ""
	if len(sliceSpellerResult[0]) != 0 {
		for i := 0; i < len(sliceSpellerResult[0]); i++ {
			spellerError += fmt.Sprintf("Word: %s, Position: %d, Length: %d, Corrects: %v ", sliceSpellerResult[0][i].Word, sliceSpellerResult[0][i].Pos, sliceSpellerResult[0][i].Len, sliceSpellerResult[0][i].S)
		}
	}
	if len(sliceSpellerResult[1]) != 0 {
		for i := 0; i < len(sliceSpellerResult[1]); i++ {
			spellerError += fmt.Sprintf("Word: %s, Position: %d, Length: %d, Corrects: %v ", sliceSpellerResult[1][i].Word, sliceSpellerResult[1][i].Pos, sliceSpellerResult[1][i].Len, sliceSpellerResult[1][i].S)
		}
	}

	if spellerError != "" {
		return fmt.Errorf("spelling error: %v", spellerError)
	}

	return nil
}
