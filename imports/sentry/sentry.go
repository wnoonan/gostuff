package sentry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PagerdutyServiceIntegration struct {
	ServiceId     int
	IntegrationId string
	ServiceName   string
}

func GetPagerDutyIntegration() ([]PagerdutyServiceIntegration, error) {
	authToken := os.Getenv("SENTRY_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("SENTRY_AUTH_TOKEN environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sentry.io/api/0/organizations/teamsnap/integrations/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response []struct {
		Id       string `json:"id"`
		Provider struct {
			Name string `json:"name"`
		} `json:"provider"`
		Config struct {
			Services []struct {
				Id      int    `json:"id"`
				Service string `json:"service"`
			} `json:"service_table"`
		} `json:"configData"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var serviceIntegrations []PagerdutyServiceIntegration

	for _, integration := range response {
		if integration.Provider.Name == "PagerDuty" {
			for _, service := range integration.Config.Services {
				if strings.Contains(service.Service, "(terraform)") {
					svcInt := PagerdutyServiceIntegration{
						ServiceId:     service.Id,
						IntegrationId: integration.Id,
						ServiceName:   service.Service,
					}
					serviceIntegrations = append(serviceIntegrations, svcInt)
				}
			}
		}
	}

	return serviceIntegrations, nil
}
