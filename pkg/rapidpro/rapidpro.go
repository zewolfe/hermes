package rapidpro

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	contentType string = "application/json"
)

type RapidProResponse struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	To         string `json:"to"`
	ToNoPlus   string `json:"to_no_plus"`
	From       string `json:"from"`
	FromNoPlus string `json:"from_no_plus"`
	Channel    string `json:"channel"`
}

type Service struct {
	o          options
	baseUrl    string
	channelId  string
	channelURL string
}

func New(url, channelId string, opts ...Options) *Service {
	options := newWithOptions(opts...)

	return &Service{
		o:          options,
		baseUrl:    url,
		channelId:  channelId,
		channelURL: url + channelId,
	}
}

func (s *Service) TriggerFlow(msisdn, requestedPath string) error {
	url := s.baseUrl + s.channelId + "/receive"

	reqPath := path.Base(requestedPath)
	if reqPath == "in" {
		reqPath = fmt.Sprintf("%s/%s", s.o.trigger, reqPath)
	}

	now := time.Now().Format(time.RFC3339)
	params := map[string]string{
		"from": msisdn,
		"text": string(reqPath),
		"date": now,
	}

	//TODO: Clean this up, please
	return s.Send(url, nil, params)
}

func (s *Service) Ack(msg RapidProResponse, success bool) error {
	if !success {
		failedURL := s.channelURL + "/failed"
		failedParams := map[string]string{
			"id": msg.ID,
		}
		if err := s.Send(failedURL, strings.NewReader(""), failedParams); err != nil {
			return err
		}
	}

	sentURL := s.channelURL + "/sent"
	sentParams := map[string]string{
		"id": msg.ID,
	}
	if err := s.Send(sentURL, strings.NewReader(""), sentParams); err != nil {
		return err
	}

	deliveredURL := s.channelURL + "/delivered"
	if err := s.Send(deliveredURL, strings.NewReader(""), sentParams); err != nil {
		return err
	}

	return nil
}

func (s *Service) Send(baseURL string, body io.Reader, params map[string]string) error {
	log := *s.o.logger

	queryParams := url.Values{}
	if len(params) > 0 {
		for key, value := range params {
			queryParams.Add(key, value)
		}
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, body)
	if err != nil {
		return fmt.Errorf("Failed to create new rapid pro request, :%v", err)
	}

	req.URL.RawQuery = queryParams.Encode()

	client := &http.Client{
		Timeout: s.o.timeout,
	}

	req.Header.Add("Authorization", s.o.token)
	req.Header.Add("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request :%v", err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info("Sent to Rapid Pro", "url", baseURL, "with response", string(res))

	return nil
}
