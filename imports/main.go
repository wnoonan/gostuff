package main

import (
	"flag"
	"fmt"

	"github.com/wnoonan/gostuff/imports/services"
	"github.com/wnoonan/gostuff/imports/users"
	"github.com/wnoonan/gostuff/imports/util"
)

func main() {
	usersFile := flag.String("users-file", "../tmp/users.json", "The file to load users from")
	servicesFile := flag.String("services-file", "../tmp/services.json", "The file to load services from")
	sentryUsers, err := users.GetSentryUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	pagerdutyUsers, err := users.GetPagerdutyUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	datadogUsers, err := users.GetDatadogUsers()
	if err != nil {
		fmt.Println(err)
		return
	}

	sentryProjects, err := services.GetSentryProjects()
	if err != nil {
		fmt.Println(err)
		return
	}

	pagerdutyServices, err := services.GetPagerdutyServices()
	if err != nil {
		fmt.Println(err)
		return
	}

	datadogServices, err := services.GetDatadogServices()
	if err != nil {
		fmt.Println(err)
		return
	}

	loadedUsers, err := util.LoadUsers(*usersFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	loadedServices, err := util.LoadServices(*servicesFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	matchedUsers := util.MatchUsers(loadedUsers, sentryUsers, pagerdutyUsers, datadogUsers)
	matchedServices := util.MatchServices(loadedServices, sentryProjects, pagerdutyServices, datadogServices)

	util.WriteUserImports(matchedUsers, "../user_imports.tf")
	util.WriteServiceImports(matchedServices, "../service_imports.tf")

	for _, service := range loadedServices {
		fmt.Println(service.Name)
	}
	// sentryUtil, err := util.NewSentryUtil()
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
