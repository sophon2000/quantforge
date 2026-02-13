package itick

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetInfo() {
	// url := "https://api.itick.org/stock/info?type=stock&region=HK&code=700"
	u := &url.URL{
		Scheme: "https",
		Host:   "api.itick.org",
		Path:   "/stock/info",
	}

	q := u.Query()
	q.Add("type", "stock")
	q.Add("region", "US")
	q.Add("code", "SE")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("token", "")

	res, _ := http.DefaultClient.Do(req)

	defer func() {
		_ = res.Body.Close()
	}()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
