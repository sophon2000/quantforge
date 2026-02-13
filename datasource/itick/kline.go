package itick

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetKline() {

	u := &url.URL{
		Scheme: "https",
		Host:   "api.itick.org",
		Path:   "/stock/kline",
	}

	q := u.Query()
	q.Add("region", "US")
	q.Add("code", "SE")
	q.Add("kType", "3")
	q.Add("limit", "10")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("token", "")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
