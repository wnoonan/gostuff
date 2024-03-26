package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SentryProject struct {
	Name string
	Id   string
}

type PagerdutyService struct {
	Name string
	Id   string
}

type DatadogService struct {
	Name string
	Id   string
}

func GetSentryProjects() ([]SentryProject, error) {
	authToken := os.Getenv("SENTRY_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("SENTRY_AUTH_TOKEN environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sentry.io/api/0/organizations/teamsnap/projects/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sentryProjects []SentryProject

	var response []struct {
		Slug string `json:"slug"`
		Id   string `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, project := range response {
		sentryProject := SentryProject{
			Name: project.Slug,
			Id:   project.Slug,
		}
		sentryProjects = append(sentryProjects, sentryProject)
	}

	return sentryProjects, nil
}

func GetPagerdutyServices() ([]PagerdutyService, error) {
	authToken := os.Getenv("PAGERDUTY_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("PAGERDUTY_TOKEN environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.pagerduty.com/services", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var pagerdutyServices []PagerdutyService

	var response struct {
		Services []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"services"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, service := range response.Services {
		pagerdutyService := PagerdutyService{
			Name: service.Name,
			Id:   service.Id,
		}
		pagerdutyServices = append(pagerdutyServices, pagerdutyService)
	}

	return pagerdutyServices, nil
}

func GetDatadogServices() ([]DatadogService, error) {
	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("DD_API_KEY environment variable not set")
	}
	appKey := os.Getenv("DD_APP_KEY")
	if appKey == "" {
		return nil, fmt.Errorf("DD_APP_KEY environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.datadoghq.com/api/v2/services/definitions", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("DD-API-KEY", apiKey)
	req.Header.Set("DD-APPLICATION-KEY", appKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var datadogServices []DatadogService

	var response struct {
		Data []struct {
			Attributes struct {
				Schema struct {
					Service string `json:"dd-service"`
				} `json:"schema"`
			} `json:"attributes"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	for _, service := range response.Data {
		datadogService := DatadogService{
			Name: service.Attributes.Schema.Service,
			Id:   service.Attributes.Schema.Service,
		}
		datadogServices = append(datadogServices, datadogService)
	}

	return datadogServices, nil
}
