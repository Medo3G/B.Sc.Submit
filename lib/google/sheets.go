package google

import (
	"fmt"
	"strings"

	"github.com/medo3g/b.sc.submit/config"
	// "github.com/medo3g/b.sc.submit/lib/util"
	sheets "google.golang.org/api/sheets/v4"
)

const (
	studentsCellRange = "'Students'!A:G"
	proposalCellRange = "'Proposals'!A:H"
)

var (
	_sheetsService *sheets.Service
)

func sheetsService() (*sheets.Service, error) {
	if _sheetsService == nil {
		c, err := googleClient()
		if err != nil {
			return nil, err
		}

		_sheetsService, err = sheets.New(c)
		if err != nil {
			return nil, err
		}
	}

	return _sheetsService, nil
}

// SheetsSubmit func
// func SheetsSubmit(teamName string, url string) error {
// 	service, err := sheetsService()
// 	if err != nil {
// 		return err
// 	}

// 	cellRange := fmt.Sprintf(config.EvaluationsCellRange, util.ParseTeamName(teamName))

// 	valueRange := &sheets.ValueRange{
// 		Values: [][]interface{}{[]interface{}{url}},
// 	}
// 	_, err = service.Spreadsheets.Values.Update(config.StudentsSheetID, cellRange, valueRange).ValueInputOption("RAW").Do()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// SheetsUserInfoBy func
func SheetsUserInfoBy(field, identifier string) (map[string]string, error) {

	service, err := sheetsService()
	if err != nil {
		return nil, err
	}

	valueRange, err := service.Spreadsheets.Values.Get(config.StudentsSheetID, studentsCellRange).Do()
	if err != nil {
		return nil, err
	}

	for _, valueRow := range valueRange.Values {
		userData := map[string]string{
			"ID":       valueRow[0].(string),
			"UserName": strings.SplitN(valueRow[5].(string), "@", 2)[0],
			"FullName": valueRow[1].(string),
			"Email":    valueRow[5].(string),
			"Group":    valueRow[2].(string),
			// "Team":      util.FormatTeamName(valueRow[3]),
			"Team":      valueRow[3].(string),
			"TeamGroup": valueRow[4].(string),
			"Category":  valueRow[6].(string),
		}
		if strings.ToLower(userData[field]) == strings.ToLower(identifier) {
			return userData, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find %s: %s", field, identifier)
}

// SheetsTeamMembers func
func SheetsTeamMembers(teamName string) ([]map[string]string, error) {
	service, err := sheetsService()
	if err != nil {
		return nil, err
	}

	valueRange, err := service.Spreadsheets.Values.Get(config.StudentsSheetID, studentsCellRange).Do()
	if err != nil {
		return nil, err
	}

	// teamID := util.TrimTeamName(teamName)
	teamID := strings.ToLower(teamName)
	members := []map[string]string{}
	for _, valueRow := range valueRange.Values {
		// if util.TrimTeamName(valueRow[3]) == teamID {
		if strings.ToLower(valueRow[3].(string)) == teamID {
			members = append(members, map[string]string{
				"ID":        valueRow[0].(string),
				"FullName":  valueRow[1].(string),
				"UserName":  strings.SplitN(valueRow[5].(string), "@", 2)[0],
				"Group":     valueRow[2].(string),
				"TeamGroup": valueRow[4].(string),
			})
		}
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("Couldn't find %s", teamName)
	}

	return members, nil
}

// SheetsTeamProposal func
func SheetsTeamProposal(teamName string) (map[string]interface{}, error) {
	service, err := sheetsService()
	if err != nil {
		return nil, err
	}

	valueRange, err := service.Spreadsheets.Values.Get(config.StudentsSheetID, proposalCellRange).Do()
	if err != nil {
		return nil, err
	}

	// teamID := util.TrimTeamName(teamName)
	teamID := strings.ToLower(teamName)
	for _, valueRow := range valueRange.Values[1:] {
		// if util.TrimTeamName(valueRow[0]) == teamID {
		if strings.ToLower(valueRow[0].(string)) == teamID {
			proposal := map[string]interface{}{
				"QAs":      [][]string{},
				"Notes":    "",
				"Late":     "",
				"Approved": false,
			}

			valueRow = valueRow[1:]
			for i, valueCol := range valueRow {
				switch i {
				case len(valueRow) - 3:
					proposal["Notes"] = valueCol.(string)
				case len(valueRow) - 2:
					if valueCol != "NO" {
						proposal["Late"] = valueCol.(string)
					}
				case len(valueRow) - 1:
					proposal["Approved"] = valueCol == "YES"
				default:
					proposal["QAs"] = append(proposal["QAs"].([][]string), []string{
						valueRange.Values[0][i+1].(string),
						valueCol.(string),
					})
				}
			}

			return proposal, nil
		}
	}

	return nil, fmt.Errorf("Couldn't find %s", teamName)
}
