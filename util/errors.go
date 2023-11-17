package errors

import (
	"net/http"
    "encoding/json"

	"github.com/gabriel-wer/picori/types"
)

func JsonData(w http.ResponseWriter, url types.URL) []byte {
    jsonData, err := json.Marshal(url)
    if err != nil{ 
        w.Write([]byte("Cannot marshall JSON Data"))
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    return jsonData
}
