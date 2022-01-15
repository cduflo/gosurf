package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"io/ioutil"
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

func getJSON() *DataResponse {

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


func main() {

	c := getJSON()

	filtered := filterWaveHeight(c.Data.Conditions, 2)




	fmt.Printf("body: %v\n", filtered)

}

// make httpService (getJSON)
// make dataService (filterConditions)
// make textService (notify)