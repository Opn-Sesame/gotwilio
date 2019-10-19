package gotwilio

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Carrier is the carrier information returned via a lookup.
type Carrier struct {
	ErrorCode         string `json:"error_code"`
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"`
}

// LookupResponse is returned after a lookup request is made to Twilio
type LookupResponse struct {
	CallerName     string `json:"caller_name"`
	Carrier        Carrier
	CountryCode    string `json:"country_code"`
	NationalFormat string `json:"national_format"`
	PhoneNumber    string `json:"phone_number"`
	Url            string `json:"url"`
}

// Lookup uses Twilio to lookup a phone number.
// See https://www.twilio.com/docs/lookup/api for more information.
func (twilio *Twilio) Lookup(ctx context.Context, number, lookupType string) (*LookupResponse, *Exception, error) {
	u, err := url.Parse(twilio.LookupUrl)
	if err != nil {
		return nil, nil, err
	}

	u.Path += number
	params := url.Values{}
	params.Add("Type", lookupType)
	u.RawQuery = params.Encode()
	res, err := twilio.getWithContext(ctx, u.String())
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		exception := new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return nil, exception, err
	}

	lookupResponse := new(LookupResponse)
	err = json.Unmarshal(responseBody, lookupResponse)
	return lookupResponse, nil, err
}
