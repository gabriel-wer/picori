package picori

import "testing"

func TestShorten(t *testing.T) {
    url := []URL{
        {"", "google.com"},
        {"ao2d9jg", "https://youtube.com"},
        {" ", "bing.com"},
    }

    for _, u := range url {
        result := u.Shorten()
        if result == nil {
            t.Logf("Shorten() PASSED. Input %s Got value %s", u.LongURL, u.ShortURL)
        }else{
            t.Errorf("Shorten() FAILED. Expected a shortened URL but got the error: %s", result.Error())
        }
    }
}
