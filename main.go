package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var version = "0.2"

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
	verb{"Benchmarked", "Benchmark"},
	verb{"Considered", "Consider"},
	verb{"Consulted", "Consult"},
	verb{"Deep dove", "Deep dive"},
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
	verb{"Prioritised", "Prioritise"},
	verb{"Researched", "Research"},
	verb{"Reviewed", "Review"},
	verb{"Roadmapped", "Roadmap"},
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
	verb{"Clawed at", "Hiss at"},
	verb{"Licked", "Lick"},
	verb{"Meowed at", "Meow at"},
	verb{"Nibbled", "Nibble"},
	verb{"Presented belly to", "Present belly to"},
	verb{"Ran away from", "Run away from"},
	verb{"Stared nerviously at", "Stare nervously at"},
	verb{"Stared out the window at", "Stare out the window at"},
	verb{"Coughed hairball onto", "Cough hairball onto"},
}

var subjects = []string{
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

var objects = []string{
	"data",
	"data lake",
	"dashboard",
	"signal control",
	"AI",
	"ingestion",
	"user flow",
	"exports",
	"workflow",
	"infrastructure",
	"data processes",
	"GDPR readiness",
	"happy path",
	"ISO compliance",
	"demo",
	"all hands",
	"pipelines",
	"interview process",
	"onboarding",
	"roadmap",
	"agenda",
	"churn",
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest() {
	fmt.Println("Running Business Cat Bot V" + version)

	today := time.Now()

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

	return verb + " " + getRandom(subjects) + " " + getRandom(objects)
}

func getRandom(array []string) string {
	return array[rand.Intn(len(array))]
}