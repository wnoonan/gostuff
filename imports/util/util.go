package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/wnoonan/gostuff/imports/services"
	"github.com/wnoonan/gostuff/imports/users"
)

type User struct {
	ModuleName            string   `json:"module_name"`
	Email                 string   `json:"email"`
	GithubUsername        string   `json:"github_username"`
	GithubMaintainerTeams []string `json:"github_maintainer_teams"`
	GithubMemberTeams     []string `json:"github_member_teams"`
	SentryUser            users.SentryUser
	PagerdutyUser         users.PagerdutyUser
	DatadogUser           users.DatadogUser
}

type Service struct {
	Name             string `json:"name"`
	SentryProject    services.SentryProject
	PagerdutyService services.PagerdutyService
	DatadogService   services.DatadogService
}

func LoadUsers(file string) ([]User, error) {
	var users []User

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func LoadServices(file string) ([]Service, error) {
	var services []Service

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func MatchUsers(users []User, sentryUsers []users.SentryUser, pagerDutyUsers []users.PagerdutyUser, datadogUsers []users.DatadogUser) []User {
	for i, user := range users {
		for _, sentryUser := range sentryUsers {
			if user.Email == sentryUser.Email {
				users[i].SentryUser = sentryUser
			}
		}
		for _, pagerDutyUser := range pagerDutyUsers {
			if user.Email == pagerDutyUser.Email {
				users[i].PagerdutyUser = pagerDutyUser
			}
		}
		for _, datadogUser := range datadogUsers {
			if user.Email == datadogUser.Email {
				users[i].DatadogUser = datadogUser
			}
		}
	}

	return users
}

func MatchServices(services []Service, sentryProjects []services.SentryProject, pagerdutyServices []services.PagerdutyService, datadogServices []services.DatadogService) []Service {
	for i, service := range services {
		serviceSlug := strings.ReplaceAll(strings.ToLower(service.Name), " ", "-")
		serviceNames := []string{service.Name, serviceSlug}

		for _, sentryProject := range sentryProjects {
			if slices.Contains(serviceNames, sentryProject.Name) {
				services[i].SentryProject = sentryProject
			}
		}
		for _, pagerdutyService := range pagerdutyServices {
			if slices.Contains(serviceNames, pagerdutyService.Name) {
				services[i].PagerdutyService = pagerdutyService
			}
		}
		for _, datadogService := range datadogServices {
			if slices.Contains(serviceNames, datadogService.Name) {
				services[i].DatadogService = datadogService
			}
		}
	}

	return services
}

func WriteUserImports(users []User, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	for _, user := range users {
		w.WriteString(
			fmt.Sprintf(
				`
# User: %s

import {
  id = "teamsnap:%s"
  to = module.%s.module.teamsnap_member[0].github_membership.member
}
				
				`,
				user.Email,
				user.GithubUsername,
				user.ModuleName,
			),
		)

		for _, team := range user.GithubMaintainerTeams {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s:%s"
  to = module.%s.module.teamsnap_member[0].github_team_membership.maintainer["%s"]
}

					 `,
					team,
					user.GithubUsername,
					user.ModuleName,
					team,
				),
			)
		}

		for _, team := range user.GithubMemberTeams {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s:%s"
  to = module.%s.module.teamsnap_member[0].github_team_membership.member["%s"]
}
					`,

					team,
					user.GithubUsername,
					user.ModuleName,
					team,
				),
			)
		}
		if user.SentryUser.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "teamsnap/%s"
  to = module.%s.sentry_organization_member.sentry-organization-member[0]
}

					 `,
					user.SentryUser.Id,
					user.ModuleName,
				),
			)
		}
		if user.PagerdutyUser.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s"
  to = module.%s.pagerduty_user.pagerduty-user[0]
}

					 `,
					user.PagerdutyUser.Id,
					user.ModuleName,
				),
			)
		}
		if user.DatadogUser.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s"
  to = module.%s.datadog_user.datadog-user[0]
}

					 `,
					user.DatadogUser.Id,
					user.ModuleName,
				),
			)
		}
	}

	w.Flush()
	f.Close()

	return nil
}

func WriteServiceImports(services []Service, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	for _, service := range services {
		w.WriteString(
			fmt.Sprintf(
				`
# Service %s\n
				`,
				service.Name,
			),
		)

		if service.SentryProject.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "teamsnap/%s"
  to = module.%s.module.sentry[0].sentry_project.this
}

					 `,
					service.SentryProject.Id,
					service.Module(),
				),
			)
		}
		if service.PagerdutyService.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s"
  to = module.%s.pagerduty_service.service[0]
}

					 `,
					service.PagerdutyService.Id,
					service.Module(),
				),
			)
		}
		if service.DatadogService.Id != "" {
			w.WriteString(
				fmt.Sprintf(
					`
import {
  id = "%s"
  to = module.%s.datadog_service_definition_yaml.service_definition[0]
}

					 `,
					service.DatadogService.Id,
					service.Module(),
				),
			)
		}
	}

	w.Flush()
	f.Close()

	return nil
}

func (s *Service) Module() string {
	return strings.ReplaceAll(strings.ToLower(s.Name), "-", "_")
}
