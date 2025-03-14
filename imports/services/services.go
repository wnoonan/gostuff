package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tidwall/pretty"
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

type DataDogMonitorResponse struct {
	Metadata struct {
		TotalCount int `json:"total_count"`
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		PageCount  int `json:"page_count"`
	} `json:"metadata"`
	Counts struct {
		Status []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"status"`
		Type []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"type"`
		Tag []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"tag"`
		Muted []struct {
			Name  bool `json:"name"`
			Count int  `json:"count"`
		} `json:"muted"`
	} `json:"counts"`
	Monitors []DataDogMonitor `json:"monitors"`
}

type DataDogMonitor struct {
	ID                   int    `json:"id"`
	OrgID                int    `json:"org_id"`
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Classification       string `json:"classification"`
	Status               string `json:"status"`
	OverallStateModified int    `json:"overall_state_modified"`
	Metrics              []any  `json:"metrics"`
	Scopes               []any  `json:"scopes"`
	Notifications        []struct {
		Name   string `json:"name"`
		Handle string `json:"handle"`
	} `json:"notifications"`
	Creator struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Handle string `json:"handle"`
	} `json:"creator"`
	MutedUntilTs       any    `json:"muted_until_ts"`
	LastTriggeredTs    int    `json:"last_triggered_ts"`
	Query              string `json:"query"`
	RestrictedRoles    []any  `json:"restricted_roles"`
	Priority           any    `json:"priority"`
	Created            int    `json:"created"`
	Modified           int    `json:"modified"`
	SuggestionMetadata struct {
		MonitorHash          string `json:"monitor_hash"`
		SuggestionIds        []any  `json:"suggestion_ids"`
		MonitorQualityIssues []any  `json:"monitor_quality_issues"`
	} `json:"suggestion_metadata"`
	SchedulingType       string   `json:"scheduling_type"`
	EvaluationWindowType string   `json:"evaluation_window_type"`
	Tags                 []string `json:"tags"`
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

func FindDatadogMonitors(query string) (*DataDogMonitorResponse, error) {
	var monitorResp DataDogMonitorResponse

	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("DD_API_KEY environment variable not set")
	}
	appKey := os.Getenv("DD_APP_KEY")
	if appKey == "" {
		return nil, fmt.Errorf("DD_APP_KEY environment variable not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.datadoghq.com/api/v1/monitor/search?query=%s&per_page=1000", url.QueryEscape(query)), nil)
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

	err = json.NewDecoder(resp.Body).Decode(&monitorResp)
	if err != nil {
		return nil, err
	}

	return &monitorResp, nil
}

func (m *DataDogMonitor) Replace(match string, replacement string) error {
	monitorJson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	monitorJson = bytes.ReplaceAll(monitorJson, []byte(match), []byte(replacement))

	err = json.Unmarshal(monitorJson, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *DataDogMonitor) PrettyPrintTerminal() error {
	mjson, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(pretty.Color(mjson, nil)))

	return nil
}

func (m *DataDogMonitor) HighlightMatch(match string) error {
	mjson, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	highlightedJson := strings.ReplaceAll(string(pretty.Color(mjson, nil)), match, fmt.Sprintf("\033[1;31m%s\033[0m", match))

	fmt.Println(highlightedJson)

	return nil
}

func (m *DataDogMonitor) Update() error {
	apiKey := os.Getenv("DD_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("DD_API_KEY environment variable not set")
	}
	appKey := os.Getenv("DD_APP_KEY")
	if appKey == "" {
		return fmt.Errorf("DD_APP_KEY environment variable not set")
	}

	monitorJson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.datadoghq.com/api/v1/monitor/%d", m.ID), bytes.NewReader(monitorJson))
	if err != nil {
		return err
	}

	req.Header.Set("DD-API-KEY", apiKey)
	req.Header.Set("DD-APPLICATION-KEY", appKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// print response body
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())

	return nil
}
