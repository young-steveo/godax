package gdax

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"

	"github.com/young-steveo/godax/message"
)

var baseURL string
var httpClient = &http.Client{Timeout: 15 * time.Second}

// Request to GDAX api.
func Request(method string, uri string, msg interface{}) (res *http.Response, err error) {
	var data []byte
	var body *bytes.Reader

	if msg == nil {
		body = bytes.NewReader(make([]byte, 0))
	} else {
		data, err = json.Marshal(msg)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	if baseURL == `` {
		baseURL = fmt.Sprintf("https://%s", os.Getenv(`GDAX_REST_HOST`))
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, uri), body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req.Header.Add(`Accept`, `application/json`)
	req.Header.Add(`Content-Type`, `application/json`)
	req.Header.Add(`User-Agent`, `GODAX Bot 1.0`)
	req.Header.Add(`CB-ACCESS-KEY`, os.Getenv(`GDAX_KEY`))               // @todo: cache this
	req.Header.Add(`CB-ACCESS-PASSPHRASE`, os.Getenv(`GDAX_PASSPHRASE`)) // @todo: cache this
	req.Header.Add(`CB-ACCESS-TIMESTAMP`, timestamp)

	sigMsg := fmt.Sprintf("%s%s%s%s", timestamp, method, uri, data)
	sig, err := message.Signature([]byte(sigMsg), os.Getenv(`GDAX_SECRET`))
	if err == nil {
		req.Header.Add(`CB-ACCESS-SIGN`, sig)
	}

	res, err = httpClient.Do(req)
	if err != nil {
		if res != nil {
			res.Body.Close()
		}
		return res, err
	}

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		errorMessage, _ := jsonparser.GetString(body, `message`)
		res.Body.Close()
		return res, errors.New(errorMessage)
	}

	return res, nil
}
