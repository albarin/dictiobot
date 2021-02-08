package words

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type API struct {
	token   string
	baseURL string
	client  *http.Client
}

func New(baseURL, token string, client *http.Client) *API {
	return &API{
		token:   token,
		baseURL: baseURL,
		client:  client,
	}
}

type responsePayload struct {
	Word    string   `json:"word"`
	Results []Result `json:"results"`
}

type Result struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
	Examples   []string `json:"examples"`
}

func (r Result) Format(word string) string {
	result := "*" + strings.Title(word) + "*\n"
	result += "" + r.Definition + "\n\n"

	if len(r.Synonyms) != 0 {
		result += "- *Synonyms*: " + strings.Join(r.Synonyms, ", ") + "\n"
	}

	if len(r.Antonyms) != 0 {
		result += "- *Antonyms*: " + strings.Join(r.Antonyms, ", ") + "\n"
	}

	if len(r.Examples) != 0 {
		result += "- *Examples*: " + strings.Join(r.Examples, ", ") + "\n"
	}

	return escape(result)
}

func escape(s string) string {
	replacer := strings.NewReplacer(
		"+", "\\+",
		"(", "\\(",
		")", "\\)",
		"!", "\\!",
		".", "\\.",
		"-", "\\-",
	)

	return replacer.Replace(s)
}

func (api *API) Word(word string) ([]Result, error) {
	url := fmt.Sprintf("%s/%s", api.baseURL, word)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", api.token)

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	var payload responsePayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return payload.Results, nil
}
