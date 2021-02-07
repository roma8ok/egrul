package egrul

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const egrulHost = "egrul.nalog.ru"

type captcha struct {
	T               string `json:"t"`
	CaptchaRequired bool   `json:"captchaRequired"`
}

type subject struct {
	FullName        string `json:"n"`
	AbbreviatedName string `json:"c"`
	INN             string `json:"i"`
	KPP             string `json:"p"`
	OGRN            string `json:"o"`
	RegisterDate    string `json:"r"`
	Head            string `json:"g"`
	Address         string `json:"a"`
}

type subjects struct {
	All []subject `json:"rows"`
}

type Scrapper struct {
	Client *http.Client
}

func (s *Scrapper) getCaptcha(query string) (c string, err error) {
	u := &url.URL{
		Scheme: "https",
		Host:   egrulHost,
	}

	resp, err := s.Client.PostForm(
		u.String(),
		map[string][]string{
			"query":                     {query},
			"vyp3CaptchaToken":          {""},
			"page":                      {""},
			"nameEq":                    {"on"},
			"region":                    {""},
			"PreventChromeAutocomplete": {""},
		})
	if err != nil {
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	var captchaResponse captcha
	err = json.NewDecoder(resp.Body).Decode(&captchaResponse)
	if err != nil {
		return
	}

	c = captchaResponse.T
	return
}

func (s *Scrapper) Find(query string) (subs []subject, err error) {
	captcha, err := s.getCaptcha(query)
	if err != nil {
		return
	}

	u := &url.URL{
		Scheme: "https",
		Host:   egrulHost,
		Path:   fmt.Sprintf("search-result/%s", captcha),
	}

	resp, err := s.Client.Get(u.String())
	if err != nil {
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	var subjectsRes subjects
	err = json.NewDecoder(resp.Body).Decode(&subjectsRes)
	if err != nil {
		return
	}

	subs = subjectsRes.All
	return
}
