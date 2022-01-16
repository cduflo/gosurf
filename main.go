package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SubCondition struct {
	MinHeight float32
	MaxHeight float32
	Rating string 
	ForecastDay string
	TimeOfDay string
}

type Condition struct {
	Am SubCondition
	Pm SubCondition
	ForecastDay string
}

type Conditions struct {
	Conditions []Condition
}

type DataResponse struct {
	Data Conditions
}

func getConditionsResponse() *DataResponse {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, "https://services.surfline.com/kbyg/regions/forecasts/conditions?subregionId=58581a836630e24c4487915a&accesstoken=5f79115fe609028d579946049802a8b91861c034", nil)

	if err != nil {
		log.Fatal(err)
	}
	
	// TODO: insert Auth header
	// req.Header.Set("User-Agent", "spacecount-tutorial")

	resp, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	conditions := new(DataResponse)

	json.Unmarshal([]byte(body), conditions)

	fmt.Println(conditions.Data.Conditions[0].Pm.MaxHeight)

	return conditions
}

func filterWaveHeight(conditions []Condition, minHeight float32) []SubCondition {
    filtered := []SubCondition{}

    for i := range conditions {
        if conditions[i].Am.MinHeight > minHeight {
			newVal := conditions[i].Am
			newVal.ForecastDay = conditions[i].ForecastDay
			newVal.TimeOfDay = "AM"
            filtered = append(filtered, newVal)
        }
		if conditions[i].Pm.MinHeight > minHeight {
			newVal := conditions[i].Pm
			newVal.ForecastDay = conditions[i].ForecastDay
			newVal.TimeOfDay = "PM"
            filtered = append(filtered, newVal)
        }
    }

	return filtered
}

func filterRating(conditions *[]Condition, minRating string) {

}

func formatMessage(filtered []SubCondition) string {
	//https://www.twilio.com/docs/glossary/what-sms-character-limit
	stringBase := "Go Surf!"

	for i := range filtered {
		stringBase = fmt.Sprint(stringBase, "\n", filtered[i].ForecastDay, " ", filtered[i].TimeOfDay, ": ", filtered[i].MinHeight, "-", filtered[i].MaxHeight, " ", filtered[i].Rating)
	}

	return stringBase
}

func sendMessage(client *twilio.RestClient, message string) {
    from := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
    to := os.Getenv("TWILIO_TO_PHONE_NUMBER")

    params := &openapi.CreateMessageParams{}
    params.SetTo(to)
    params.SetFrom(from)
    params.SetBody(message)

    resp, err := client.ApiV2010.CreateMessage(params)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        response, _ := json.Marshal(*resp)
        fmt.Println("Response: " + string(response))
    }
}

func loadEnvVars() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}


func main() {
	loadEnvVars()

	client := twilio.NewRestClient()

	rsp := getConditionsResponse()

	filtered := filterWaveHeight(rsp.Data.Conditions, 2)

	if len(filtered) > 0 {
		message := formatMessage(filtered)
		fmt.Printf("message", message);
		sendMessage(client, message)
	}
}

// make httpService (getJSON)
// make dataService (filterConditions)
// make textService (notify)