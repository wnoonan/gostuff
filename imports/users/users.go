package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SentryUser struct {
	Email string
	Id    string
}

type PagerdutyUser struct {
	Email string
	Id    string
}

type DatadogUser struct {
	Email string
	Id    string
}

func GetSentryUsers() ([]SentryUser, error) {
	authToken := os.Getenv("SENTRY_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("SENTRY_AUTH_TOKEN environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sentry.io/api/0/organizations/teamsnap/members/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sentryUsers []SentryUser

	var response []struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, user := range response {
		sentryUser := SentryUser{
			Email: user.Email,
			Id:    user.Id,
		}
		sentryUsers = append(sentryUsers, sentryUser)
	}

	return sentryUsers, nil
}

func GetPagerdutyUsers() ([]PagerdutyUser, error) {
	authToken := os.Getenv("PAGERDUTY_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("PAGERDUTY_TOKEN environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.pagerduty.com/users", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pagerdutyUsers []PagerdutyUser

	var response struct {
		Users []struct {
			Email string `json:"email"`
			Id    string `json:"id"`
		} `json:"users"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, user := range response.Users {
		pagerdutyUser := PagerdutyUser{
			Email: user.Email,
			Id:    user.Id,
		}
		pagerdutyUsers = append(pagerdutyUsers, pagerdutyUser)
	}

	return pagerdutyUsers, nil
}

func GetDatadogUsers() ([]DatadogUser, error) {
	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("DD_API_KEY environment variable not set")
	}

	appKey := os.Getenv("DD_APP_KEY")
	if appKey == "" {
		return nil, fmt.Errorf("DD_APP_KEY environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.datadoghq.com/api/v2/users", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page[size]", "100")
	q.Add("filter[status]", "Active")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("DD-API-KEY", apiKey)
	req.Header.Set("DD-APPLICATION-KEY", appKey)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var datadogUsers []DatadogUser

	var response struct {
		Data []struct {
			Id         string `json:"id"`
			Attributes struct {
				Email string `json:"email"`
			} `json:"attributes"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, user := range response.Data {

		datadogUser := DatadogUser{
			Email: user.Attributes.Email,
			Id:    user.Id,
		}
		datadogUsers = append(datadogUsers, datadogUser)
	}

	return datadogUsers, nil
}
