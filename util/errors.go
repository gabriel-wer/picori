package errors

import (
	"net/http"
    "encoding/json"

	"github.com/gabriel-wer/picori/picori"
)

func JsonData(w http.ResponseWriter, url picori.URL) []byte {
    jsonData, err := json.Marshal(url)
    if err != nil{ 
        w.Write([]byte("Cannot marshall JSON Data"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return jsonData
}
