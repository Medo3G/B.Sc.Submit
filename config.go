package main

import (
	// "log"
	// "encoding/json"
	"os"
	"strings"

	"github.com/mostafa-alaa-494/b.sc.submit/config"
)

func init() {

	config.SubmitName = os.Getenv("CONFIG_SUBMIT_NAME")
	config.AdminPassword = os.Getenv("CONFIG_ADMIN_PASSWORD")
	config.TestPassword = os.Getenv("CONFIG_TEST_PASSWORD")
	config.SubmissionDeadline = os.Getenv("CONFIG_SUBMISSION_DEADLINE")
	config.TeamNameFormat = os.Getenv("CONFIG_TEAM_NAME_FORMAT")
	config.FeaturesEnabled = map[string]bool{
		"proposals":    os.Getenv("CONFIG_FEATURE_ENABLED_PROPOSALS") == "1",
		"submissions":  os.Getenv("CONFIG_FEATURE_ENABLED_SUBMISSIONS") == "1",
		"reservations": os.Getenv("CONFIG_FEATURE_ENABLED_EVALUATIONS") == "1",
		"embed":        os.Getenv("CONFIG_FEATURE_ENABLED_EMBED") == "1",
		"settings":     os.Getenv("CONFIG_FEATURE_ENABLED_SETTINGS") == "1",
	}

	config.GoogleAPIClientSecret = os.Getenv("CONFIG_GOOGLE_API_CLIENT_SECRET")
	config.GoogleAPIClientToken = os.Getenv("CONFIG_GOOGLE_API_CLIENT_TOKEN")
	config.StudentsSheetID = os.Getenv("CONFIG_STUDENTS_SHEET_ID")
	// json.Unmarshal([]byte(os.Getenv("CONFIG_SUBMISSIONS_ITEMS")), &config.SubmissionsItems)
	// config.SubmissionsFolderID = os.Getenv("CONFIG_SUBMISSIONS_FOLDER_ID")
	// config.SubmissionsMetaDescription = os.Getenv("CONFIG_SUBMISSIONS_META_DESCRIPTION")
	config.EvaluationsCellRange = os.Getenv("CONFIG_EVALUATIONS_CELL_RANGE")

	config.EvaluationsCalendarID = os.Getenv("CONFIG_EVALUATIONS_CALENDAR_ID")
	config.EvaluationsWeekStart = os.Getenv("CONFIG_EVALUATIONS_WEEK_START")
	config.EvaluationsWeekEnd = os.Getenv("CONFIG_EVALUATIONS_WEEK_END")
	config.EvaluationsCalendarEmbed = os.Getenv("CONFIG_EVALUATIONS_CALENDAR_EMBED")
	config.ReservationDaysAhead = os.Getenv("CONFIG_RESERVATION_DAYS_AHEAD")

	config.VRStudetsSheetID = os.Getenv("CONFIG_VR_STUDENTS_SHEET_ID")
	config.VRCalendarID = os.Getenv("CONFIG_VR_CALENDAR_ID")
	config.BscVRWeeklyMinutes = os.Getenv("CONFIG_BCS_HOLO_WEEKLY_MINUTES")

	config.HoloStudentsSheetID = os.Getenv("CONFIG_HOLO_STUDENTS_SHEET_ID")
	config.HoloCalendarID = os.Getenv("CONFIG_HOLO_CALENDAR_ID")
	config.BscHoloWeeklyMinutes = os.Getenv("CONFIG_BCS_HOLO_WEEKLY_MINUTES")

	config.SlackTestToken = os.Getenv("CONFIG_SLACK_TEST_TOKEN")
	config.SlackUserToken = os.Getenv("CONFIG_SLACK_USER_TOKEN")
	config.SlackBotToken = os.Getenv("CONFIG_SLACK_BOT_TOKEN")
	config.SlackWebhookToken = os.Getenv("CONFIG_SLACK_WEBHOOK_TOKEN")
	config.SlackAdmins = strings.Split(os.Getenv("CONFIG_SLACK_ADMINS"), ",")
}
