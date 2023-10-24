package main

import (
	"CMS/backendServices/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/manveru/faker"
)

var fixtureList []string
var teamScores map[string]int = make(map[string]int)

func main() {
	rand.Seed(time.Now().UnixNano())
	baseURL := "http://localhost:8080"

	// Welcome Page
	welcomePageURL := baseURL + "/"
	welcomePageResponse, err := http.Get(welcomePageURL)
	if err != nil {
		fmt.Println("Welcome Page GET request error:", err)
		return
	}
	defer welcomePageResponse.Body.Close()
	responseBody, err := io.ReadAll(welcomePageResponse.Body)
	if err != nil {
		fmt.Println("Response body read error:", err)
		return
	}

	fmt.Println(string(responseBody))

	// Create Squads
	teamList := []string{}
	matchNumber := 1
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 4; i++ {
		playerList := []string{}
		for j := 0; j < 15; j++ {
			// Create a player
			playerID := createPlayer(baseURL)
			playerList = append(playerList, playerID)
		}

		playingXI := playerList[:11]

		teamName := fake.CompanyName()
		// Create a team
		teamID := createTeam(baseURL, teamName, playerList, playingXI)
		teamList = append(teamList, teamID)
	}

	// Create Fixtures
	createFixtures(baseURL, teamList, &matchNumber)

	// Create Tournament
	for i := 0; i < len(fixtureList); i++ {
		fixture := getFixtures(fixtureList[i])
		teamA := fixture["teamA"]
		teamB := fixture["teamB"]
		simulateMatch(teamA.(string), teamB.(string))
	}

	for i := 0; i < len(teamList); i++ {
		team := getTeams(teamList[i])
		value := team["matchPoints"]
		teamScores[teamList[i]] = int(value.(float64))
	}

	// Sort Teams by Scores
	sortedTeams := sortTeamsByScores(teamList, teamScores)
	topTeams := sortedTeams[:2]
	fmt.Println("\nTop 2 Teams for the Final Match:")
	simulateMatch(topTeams[0], topTeams[1])
}

func sortTeamsByScores(teams []string, scores map[string]int) []string {
	sortedTeams := make([]string, len(teams))
	copy(sortedTeams, teams)

	// Custom sorting function to sort teams by scores
	sort.Slice(sortedTeams, func(i, j int) bool {
		return scores[sortedTeams[i]] > scores[sortedTeams[j]]
	})

	return sortedTeams
}

func simulateMatch(teamA, teamB string) {
	// Get Teams
	rand.Seed(time.Now().UnixNano())
	teamAResp := getTeams(teamA)
	teamBResp := getTeams(teamB)

	// Update Match Stats
	teamAResp["matchesPlayed"] = teamAResp["matchesPlayed"].(float64) + 1
	teamBResp["matchesPlayed"] = teamBResp["matchesPlayed"].(float64) + 1

	matchStat := rand.Intn(2)
	if matchStat == 1 {
		teamAResp["matchesTied"] = teamAResp["matchesTied"].(float64) + 1
		teamBResp["matchesTied"] = teamBResp["matchesTied"].(float64) + 1
		teamAResp["matchPoints"] = teamAResp["matchPoints"].(float64) + 1
		teamBResp["matchPoints"] = teamBResp["matchPoints"].(float64) + 1
	} else if matchStat == 2 {
		teamAResp["matchesWon"] = teamAResp["matchesWon"].(float64) + 1
		teamBResp["matchesLost"] = teamBResp["matchesLost"].(float64) + 1
		teamAResp["matchPoints"] = teamAResp["matchPoints"].(float64) + 2
	} else {
		teamAResp["matchesLost"] = teamAResp["matchesLost"].(float64) + 1
		teamBResp["matchesWon"] = teamBResp["matchesWon"].(float64) + 1
		teamBResp["matchPoints"] = teamBResp["matchPoints"].(float64) + 2
	}

	// Update Teams
	updateTeam(teamA, teamAResp)
	updateTeam(teamB, teamBResp)
}

func updateTeam(teamId string, teamReq map[string]interface{}) {
	// Update Team

	team := models.Team{
		MatchesPlayed: int(teamReq["matchesPlayed"].(float64)),
		MatchesWon:    int(teamReq["matchesWon"].(float64)),
		MatchesLost:   int(teamReq["matchesLost"].(float64)),
		MatchesTied:   int(teamReq["matchesTied"].(float64)),
		MatchPoints:   int(teamReq["matchPoints"].(float64)),
	}
	println(encodeJSON(team))
	req, err := http.NewRequest("PUT", "http://localhost:8080/update-team"+teamId, encodeJSON(team))
	if err != nil {
		log.Println("Error in request")
		return
	}
	teamResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error in request")
		return
	}
	defer teamResponse.Body.Close()

	var teamJsonResponse map[string]interface{}
	if err := json.NewDecoder(teamResponse.Body).Decode(&teamJsonResponse); err != nil {
		fmt.Println("Update Team response decoding error:", err)
		return
	}

	_, dataExists := teamJsonResponse["data"].(map[string]interface{})
	if !dataExists {
		fmt.Println("Data field not found in Update Team response")
		return
	}
}

func getTeams(id string) map[string]interface{} {
	// Get Teams
	teamResponse, err := http.Get("http://localhost:8080/get-team/" + id)
	if err != nil {
		fmt.Println("Get Team GET request error:", err)
		return nil
	}
	defer teamResponse.Body.Close()

	var teamJsonResponse map[string]interface{}
	if err := json.NewDecoder(teamResponse.Body).Decode(&teamJsonResponse); err != nil {
		fmt.Println("Get Team response decoding error:", err)
		return nil
	}

	data, dataExists := teamJsonResponse["data"].(map[string]interface{})
	if !dataExists {
		fmt.Println("Data field not found in Get Team response")
		return nil
	}
	teamResp, _ := data["data"].(map[string]interface{})
	return teamResp
}

func getFixtures(id string) map[string]interface{} {
	// Get Fixtures
	fixtureResponse, err := http.Get("http://localhost:8080/get-fixture/" + id)
	if err != nil {
		fmt.Println("Get Fixture GET request error:", err)
		return nil
	}
	defer fixtureResponse.Body.Close()

	var fixtureJsonResponse map[string]interface{}
	if err := json.NewDecoder(fixtureResponse.Body).Decode(&fixtureJsonResponse); err != nil {
		fmt.Println("Get Fixture response decoding error:", err)
		return nil
	}

	data, dataExists := fixtureJsonResponse["data"].(map[string]interface{})
	if !dataExists {
		fmt.Println("Data field not found in Get Fixture response")
		return nil
	}
	fixtureResp, _ := data["data"].(map[string]interface{})
	return fixtureResp
}

func createPlayer(baseURL string) string {
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	typeValue := map[string][]string{
		"Batsman":     {"Wicket Keeper", "Top Order", "Middle Order", "Lower Order"},
		"Bowler":      {"Leg Spin", "Off Spin", "Leg Break", "Off Break", "Medium Pace", "Fast Bowler"},
		"All Rounder": {"Batting All Rounder", "Bowling All Rounder"},
	}

	typeValueArray := []string{"Batsman", "Bowler", "All Rounder"}
	randomType := typeValueArray[rand.Intn(len(typeValueArray))]
	value := typeValue[randomType][rand.Intn(len(typeValue[randomType]))]

	payload := models.Player{
		Name:          fake.Name(),
		Age:           rand.Intn(23) + 18,
		DominantHand:  []string{"Left Handed", "Right Handed"}[rand.Intn(2)],
		Role:          []string{randomType, value},
		BattingAvg:    rand.Float32()*(50-20) + 20,
		StrikeRate:    rand.Float32()*(200-100) + 100,
		Economy:       rand.Float32()*(10-5) + 5,
		MatchesPlayed: rand.Intn(91) + 10,
		TotalRuns:     rand.Intn(501) + 500,
		TotalWickets:  rand.Intn(51) + 50,
	}

	playerResponse, err := http.Post(baseURL+"/create-player", "application/json", encodeJSON(payload))
	if err != nil {
		fmt.Println("Create Player POST request error:", err)
		return ""
	}
	defer playerResponse.Body.Close()

	var playerJsonResponse map[string]interface{}
	if err := json.NewDecoder(playerResponse.Body).Decode(&playerJsonResponse); err != nil {
		fmt.Println("Create Player response decoding error:", err)
		return ""
	}

	data, dataExists := playerJsonResponse["data"].(map[string]interface{})
	if !dataExists {
		fmt.Println("Data field not found in Create Player response")
		return ""
	}
	respData, _ := data["data"].(map[string]interface{})
	// Check for existence of InsertedID in data field
	insertedID, idExists := respData["InsertedID"].(string)
	if !idExists {
		fmt.Println("InsertedID field not found in Create Player response")
		return ""
	}

	return insertedID
}

func createTeam(baseURL, teamName string, playerList, playingXI []string) string {
	teamPayload := models.Team{
		Name:          teamName,
		Members:       playerList,
		MatchesPlayed: 0,
		MatchesWon:    0,
		MatchesLost:   0,
		MatchesTied:   0,
		MatchPoints:   0,
		Captain:       playingXI[rand.Intn(len(playingXI))],
		PlayingXI:     playingXI,
	}

	teamResponse, err := http.Post(baseURL+"/create-team", "application/json", encodeJSON(teamPayload))
	if err != nil {
		fmt.Println("Create Team POST request error:", err)
		return ""
	}
	defer teamResponse.Body.Close()

	var teamJsonResponse map[string]interface{}
	if err := json.NewDecoder(teamResponse.Body).Decode(&teamJsonResponse); err != nil {
		fmt.Println("Create Team response decoding error:", err)
		return ""
	}

	data, dataExists := teamJsonResponse["data"].(map[string]interface{})
	if !dataExists {
		fmt.Println("Data field not found in Create Team response")
		return ""
	}
	teamResp, _ := data["data"].(map[string]interface{})
	// Check for existence of InsertedID in data field
	insertedID, idExists := teamResp["InsertedID"].(string)
	if !idExists {
		fmt.Println("InsertedID field not found in Create Team response")
		return ""
	}

	return insertedID

}

func createFixtures(baseURL string, teamList []string, matchNumber *int) {
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	for i := 0; i+1 < 4; i++ {
		for j := i + 1; j < 4; j++ {
			venue := fake.City()
			matchDate := time.Now().Add(time.Duration(rand.Intn(30)) * 24 * time.Hour)

			payload := models.Fixture{
				MatchNumber:  *matchNumber,
				TeamA:        teamList[i],
				TeamB:        teamList[j],
				Venue:        venue,
				MatchDate:    matchDate,
				IsFinalMatch: false,
			}

			matchResponse, err := http.Post(baseURL+"/create-fixture", "application/json", encodeJSON(payload))
			if err != nil {
				fmt.Println("Create Fixture POST request error:", err)
				return
			}
			defer matchResponse.Body.Close()

			var teamJsonResponse map[string]interface{}
			if err := json.NewDecoder(matchResponse.Body).Decode(&teamJsonResponse); err != nil {
				fmt.Println("Create Fixture response decoding error:", err)
				break
			}

			data, dataExists := teamJsonResponse["data"].(map[string]interface{})
			if !dataExists {
				fmt.Println("Data field not found in Create Fixture response")
				break
			}

			matchResp, _ := data["data"].(map[string]interface{})
			// Check for existence of InsertedID in data field
			insertedID, idExists := matchResp["InsertedID"].(string)
			if !idExists {
				fmt.Println("InsertedID field not found in Create Fixture response")
				break
			}
			fixtureList = append(fixtureList, insertedID)
			*matchNumber++
		}
	}
}

func encodeJSON(data interface{}) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	_ = json.NewEncoder(buffer).Encode(data)
	return buffer
}
