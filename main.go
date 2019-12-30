package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/intenthq/golukay"
)

var version = "0.4-CHRISTMAS-EDITION.1"

type slackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type verb struct {
	Past    string
	Present string
}

var verbs = []verb{
	// Business Verbs
	verb{"Activated", "Activate"},
	verb{"Aggregated", "Aggregate"},
	verb{"Benchmarked", "Benchmark"},
	verb{"Considered", "Consider"},
	verb{"Connected", "Connect"},
	verb{"Consulted", "Consult"},
	verb{"Deep dove", "Deep dive"},
	verb{"Deployed", "Deploy"},
	verb{"Designed", "Design"},
	verb{"Discontinued", "Discontinue"},
	verb{"Discussed", "Discuss"},
	verb{"Disrupted", "Disrupt"},
	verb{"Drilled down", "Drill down"},
	verb{"Engaged", "Engage"},
	verb{"Ignited", "Ignite"},
	verb{"Incentivised", "Incentivise"},
	verb{"Integrated", "Integrate"},
	verb{"Kick started", "Kick start"},
	verb{"Optimised", "Optimise"},
	verb{"Planned", "Plan"},
	verb{"Positioned", "Position"},
	verb{"Prioritised", "Prioritise"},
	verb{"Researched", "Research"},
	verb{"Reviewed", "Review"},
	verb{"Roadmapped", "Roadmap"},
	verb{"Rolled out", "Roll-out"},
	verb{"Scaled", "Scale"},
	verb{"Scoped", "Scope"},
	verb{"Specced out", "Spec out"},
	verb{"Sued", "Sue"},
	verb{"Synergized", "Synergize"},
	verb{"Upgraded", "Upgrade"},

	// Cat Verbs
	verb{"Attacked", "Attack"},
	verb{"Curiously pawed", "Curiously paw"},
	verb{"Curled up by", "Curl up by"},
	verb{"Hid from", "Hide from"},
	verb{"Hunted", "Hunt"},
	verb{"Clawed at", "Hiss at"},
	verb{"Licked", "Lick"},
	verb{"Meowed at", "Meow at"},
	verb{"Nibbled", "Nibble"},
	verb{"Presented belly to", "Present belly to"},
	verb{"Ran away from", "Run away from"},
	verb{"Stared nerviously at", "Stare nervously at"},
	verb{"Stared out the window at", "Stare out the window at"},
	verb{"Coughed hairball onto", "Cough hairball onto"},
	verb{"Knocked over", "Knock over"},
	verb{"Vomited on", "Vomit on"},
	verb{"Hissed at", "Hiss at"},
	verb{"Kneaded", "Knead"},
	verb{"Sat on", "Sit on"},

	// Christmas Verbs
	verb{"Decorated", "Decorate"},
	verb{"Prepared", "Prepare"},
	verb{"Ate", "Eat"},
	verb{"Gift wrapped", "Gift wrap"},
	verb{"Sung at", "Sing at"},
	verb{"Worshipped", "Worship"},
}

var defaultSubjects = []string{
	"our",
	"our",
	"our",
	"our", // A simple way to increase the frequency of occurence
	"client’s",
	"new client’s",
	"prospective client’s",
	"hade’s",
	"literally anyone’s",
	"the government’s",
	"my",
	"vendor's",
}

var defaultObjects = []string{
	"naughty list",
	"nice list",
	"elves",
	"sleigh",
	"workshop",
	"presents",
	"carol singers",
	"merriment",
	"christmas tree",
	"christmas cake",
	"christmas miracle",
	"mulled wine",
	"snowman",
	"mince pies",
	"three wise men",
	"stocking",
	"turkey",
	"wrapping paper",
}

func isBankHoliday(holidays []golukay.BankHoliday, date time.Time) bool {
	for _, holiday := range holidays {
		y1, m1, d1 := holiday.Date.Date()
		y2, m2, d2 := date.Date()

		if y1 == y2 && m1 == m2 && d1 == d2 {
			return true
		}
	}
	return false
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest() {
	fmt.Println("Running Business Cat Bot V" + version)

	_, dryRun := os.LookupEnv("DRY_RUN")
	if dryRun {
		fmt.Println("Running in dry run mode, will not post message to Slack")
	}

	today := time.Now()

	// Checks if it is a bank holiday (no standup should be posted)
	holidays, err := golukay.GetHolidays()

	if err == nil {
		if isBankHoliday(holidays.EnglandAndWales.Events, today) {
			fmt.Println("Bank holiday today, no standup message needed")
			return
		}
	}

	var lastWorkday time.Time
	var lastWorkdayName string

	if today.Weekday() == time.Monday {
		lastWorkday = today.AddDate(0, 0, -3)
		lastWorkdayName = "F"
		fmt.Println("It is a Monday")
	} else {
		lastWorkday = today.AddDate(0, 0, -1)
		lastWorkdayName = "Y"
		fmt.Println("It is not a Monday")
	}

	standupMessage := fmt.Sprintf("%s: %s\nT: %s", lastWorkdayName, getStandup(lastWorkday, true), getStandup(today, false))

	fmt.Println("Generated standup message:\n---\n" + standupMessage + "\n---")

	slackHook := os.Getenv("SLACK_TOKEN")
	payload := slackResponse{
		ResponseType: "channel",
		Text:         standupMessage,
	}

	jsonValue, _ := json.Marshal(payload)

	if dryRun {
		fmt.Println("Skipping posting to Slack")

		return
	}

	fmt.Println("Posting standup to Slack...")
	response, err := http.Post(slackHook, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println("Network error encountered:")
		fmt.Println(err)
	} else if response.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		errorMessage := buf.String()
		fmt.Println("Failed to send standup. Response from Slack\n" + errorMessage)
	} else {
		fmt.Println("Successfully sent standup")
	}
}

func getStandup(date time.Time, past bool) string {
	seed := int64(date.Year() * int(date.Month()) * date.Day())
	rand.Seed(seed)

	verbForms := verbs[rand.Intn(len(verbs))]
	var verb string

	if past {
		verb = verbForms.Past
	} else {
		verb = verbForms.Present
	}

	var subjects, objects []string

	customSubjects, hasCustomSubjects := os.LookupEnv("EXTRA_SUBJECTS")

	if hasCustomSubjects {
		fmt.Println("Adding custom subjects to standup message generator")
		subjects = append(defaultSubjects, strings.Split(customSubjects, ",")...)
	} else {
		subjects = defaultSubjects
	}

	customObjects, hasCustomObjects := os.LookupEnv("EXTRA_OBJECTS")

	if hasCustomObjects {
		fmt.Println("Adding custom objects to standup message generator")
		objects = append(defaultObjects, strings.Split(customObjects, ",")...)
	} else {
		objects = defaultObjects
	}

	return verb + " " + getRandom(subjects) + " " + getRandom(objects)
}

func getRandom(array []string) string {
	return array[rand.Intn(len(array))]
}
