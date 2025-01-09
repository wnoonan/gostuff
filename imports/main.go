package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wnoonan/gostuff/imports/services"
)

func StringPrompt(prompt string) string {
	fmt.Printf("%s: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func InteractivelyReplace(m services.DataDogMonitor) {
	match := StringPrompt("Enter match string")

	m.HighlightMatch(match)

	replacement := StringPrompt("Enter replacement string")

	m.Replace(match, replacement)

	m.PrettyPrintTerminal()

	more := StringPrompt("More replacements? (y/n)")

	if more == "y" {
		InteractivelyReplace(m)
	}
}

func main() {

	// matches := []string{
	// 	"apiv3-controller",
	// 	"apiv3-cache",
	// 	"apiv3-db",
	// 	"apiv3-outbound",
	// 	"apiv3-rack",
	// 	"apiv3-redis",
	// 	"apiv3-sidekiq",
	// 	"apiv3-sinatra",
	// 	"apiv3-hoyle",
	// }

	monResp, err := services.FindDatadogMonitors("(apiv3* OR snapi*) AND (NOT *proxysql*)")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d monitors\n", len(monResp.Monitors))

	for _, m := range monResp.Monitors {
		m.PrettyPrintTerminal()
		skip := StringPrompt("Skip this monitor? (y/n)")
		if skip == "y" {
			continue
		}

		InteractivelyReplace(m)
		update := StringPrompt("Update this monitor? (y/n)")
		if update == "y" {
			if err := m.Update(); err != nil {
				fmt.Println(err)
			}

		}

		fmt.Printf("Monitor %d updated\n", m.ID)
		onward := StringPrompt("Continue? (y/n)")
		if onward == "n" {
			break
		}
	}

	// replace all instances of "apiv3" with "snapi-monolith"

	// usersFile := flag.String("users-file", "../tmp/users.json", "The file to load users from")
	// servicesFile := flag.String("services-file", "../tmp/services.json", "The file to load services from")
	// sentryUsers, err := users.GetSentryUsers()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pagerdutyUsers, err := users.GetPagerdutyUsers()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// datadogUsers, err := users.GetDatadogUsers()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// sentryProjects, err := services.GetSentryProjects()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pagerdutyServices, err := services.GetPagerdutyServices()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// datadogServices, err := services.GetDatadogServices()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// loadedUsers, err := util.LoadUsers(*usersFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// loadedServices, err := util.LoadServices(*servicesFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// matchedUsers := util.MatchUsers(loadedUsers, sentryUsers, pagerdutyUsers, datadogUsers)
	// matchedServices := util.MatchServices(loadedServices, sentryProjects, pagerdutyServices, datadogServices)

	// util.WriteUserImports(matchedUsers, "../user_imports.tf")
	// util.WriteServiceImports(matchedServices, "../service_imports.tf")

	// for _, service := range loadedServices {
	// 	fmt.Println(service.Name)
	// }
	// // sentryUtil, err := util.NewSentryUtil()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// sentryProjects, err := services.GetSentryProjects()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// alertRules, err := sentryUtil.ProjectsPagerdutyIssueAlertRules(&sentryProjects)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = util.WriteAlertRulesToFile(alertRules, "../alert_rules.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
