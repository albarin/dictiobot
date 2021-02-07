package words

import (
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestName(t *testing.T) {
	api := New("https://wordsapiv1.p.rapidapi.com/words", "b496d18214msh600d7ae382d49f1p16bb60jsnd2b26fae7025", http.DefaultClient)
	r, _ := api.Word("big")
	spew.Dump(r)
}
