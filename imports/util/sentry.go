package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/wnoonan/gostuff/imports/services"
)

type SentryUtil struct {
	client    *http.Client
	authToken string
}

type SentryIssueAlert struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Projects []string `json:"projects"`
	Actions  []struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		ServiceId string `json:"service"`
		Enabled   bool   `json:"enabled"`
	} `json:"actions"`
}

type SentryIntegration struct {
	ID       string `json:"id"`
	Provider struct {
		Key     string `json:"key"`
		Aspects struct {
		} `json:"aspects"`
	} `json:"provider"`
	ConfigData struct {
		ServiceTable []struct {
			Service        string `json:"service"`
			IntegrationKey string `json:"integration_key"`
			Id             int    `json:"id"`
		} `json:"service_table"`
	} `json:"configData,omitempty"`
}

func NewSentryUtil() (*SentryUtil, error) {
	authToken := os.Getenv("SENTRY_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("SENTRY_AUTH_TOKEN environment variable not set")
	}

	return &SentryUtil{
		client:    &http.Client{},
		authToken: authToken,
	}, nil
}

func (s *SentryUtil) PagerdutyIntegration() (*SentryIntegration, error) {
	var sentryPagerdutyIntegrations []SentryIntegration

	req, err := http.NewRequest("GET", "https://sentry.io/api/0/organizations/teamsnap/integrations/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.authToken))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&sentryPagerdutyIntegrations)
	if err != nil {
		return nil, err
	}

	var pagerdutyIntegration SentryIntegration
	for _, integration := range sentryPagerdutyIntegrations {
		if integration.Provider.Key == "pagerduty" && len(integration.ConfigData.ServiceTable) > 0 {
			pagerdutyIntegration = integration
		}
	}

	return &pagerdutyIntegration, nil
}

func (s *SentryUtil) ProjectsPagerdutyIssueAlertRules(projects *[]services.SentryProject) (*[]SentryIssueAlert, error) {

	var sentryIssueAlerts []SentryIssueAlert

	for _, project := range *projects {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://sentry.io/api/0/projects/teamsnap/%s/rules/", project.Name), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.authToken))

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		var pagerdutyAlerts []SentryIssueAlert
		err = json.NewDecoder(resp.Body).Decode(&pagerdutyAlerts)
		if err != nil {
			return nil, err
		}

		for _, alert := range pagerdutyAlerts {
			for _, action := range alert.Actions {
				if strings.Contains(action.Id, "PagerDuty") {
					sentryIssueAlerts = append(sentryIssueAlerts, alert)
				}
			}
		}
	}

	return &sentryIssueAlerts, nil
}

func WriteAlertRulesToFile(alerts *[]SentryIssueAlert, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	for index, alert := range *alerts {
		w.WriteString(
			fmt.Sprintf(

				`
	Alert #: %v
	Alert Name: %s
	Pagerduty Integration Service Id: %s
	Action Name: %s
	Url: https://teamsnap.sentry.io/alerts/rules/%s/%s/
				`, index, alert.Name, alert.Actions[0].ServiceId, alert.Actions[0].Name, alert.Projects[0], alert.Id,
			),
		)
	}

	w.Flush()
	f.Close()

	return nil
}
