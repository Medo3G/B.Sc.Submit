package google

import (
	"log"
	"strings"
	"fmt"
	"time"
	"strconv"


	"github.com/mostafa-alaa-494/b.sc.submit/config"
	calendar "google.golang.org/api/calendar/v3"
)

var (
	_calendarService *calendar.Service
)

func calendarService() (*calendar.Service, error) {
	if _calendarService == nil {
		c, err := googleClient()
		if err != nil {
			return nil, err
		}

		_calendarService, err = calendar.New(c)
		if err != nil {
			return nil, err
		}
	}

	return _calendarService, nil
}

// CalendarFreeSlots func
func CalendarFreeSlots() ([]*calendar.Event, error) {
	service, err := calendarService()
	if err != nil {
		return nil, err
	}

	timeMin := CalendarStartOfTheWeek(time.Now())//time.Parse(time.RFC3339, config.EvaluationsWeekStart)
	// timeMax, _ := time.Parse(time.RFC3339, config.EvaluationsWeekEnd)
	daysAhead,_ := strconv.Atoi(config.ReservationDaysAhead)
	timeMax := timeMin.AddDate(0, 0, daysAhead)//.Format(time.RFC3339)
	timeNow := time.Now()
	if timeNow.After(timeMin) {
		timeMin = timeNow
	}

	slots, err := service.Events.
		List(config.EvaluationsCalendarID).
		SingleEvents(true).
		OrderBy("startTime").
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		Q("FREE").
		Do()
	if err != nil {
		return nil, err
	}

	categorySlots := []*calendar.Event{}
	for _,e := range slots.Items {

		// time.Date
		// if(e.End.DateTime - e.Start.DateTime > 1){
		categorySlots = append(categorySlots, e)
		// }

	}
	return slots.Items, nil
}

// CalendarTeamSlot func
func CalendarTeamSlot(teamName string) (*calendar.Event, error) {
	service, err := calendarService()
	if err != nil {
		return nil, err
	}

	slots, err := service.Events.
		List(config.EvaluationsCalendarID).
		SingleEvents(true).
		MaxResults(1).
		TimeMin(config.EvaluationsWeekStart).
		TimeMax(config.EvaluationsWeekEnd).
		Q(teamName).
		Do()
	if err != nil {
		return nil, err
	}

	if len(slots.Items) == 0 {
		return nil, nil
	}

	return slots.Items[0], nil
}

// CalendarStartOfTheWeek func
func CalendarStartOfTheWeek(weekDay time.Time) (time.Time){

	start := weekDay
	
	for(start.Weekday() != time.Saturday){
		start = start.AddDate(0,0,-1)
	}
	for(start.Hour() != 0){
		start = start.Add(time.Hour*-1)
	}
	return start

}

// CalendarAllTeamSlotsInWeek func
func CalendarAllTeamSlotsInWeek(teamName string, weekEvent *calendar.Event) (*calendar.Events, error) {
	service, err := calendarService()
	if err != nil {
		return nil, err
	}

	s,_ := time.Parse(time.RFC3339, weekEvent.Start.DateTime)
	start := CalendarStartOfTheWeek(s)
	end := start.AddDate(0, 0, 7)
	
	startString := strings.Replace(start.String(), " ","T",1)
	endString := strings.Replace(end.String(), " ","T",1)
	startString = strings.Replace(startString, " ","",1)
	endString = strings.Replace(endString, " ","",1)
	startString = strings.Replace(startString, " EET","",1)
	endString = strings.Replace(endString, " EET","",1)

	log.Println("Start: ",startString)
	log.Println("End: ",endString)

	
	slots, err := service.Events.
		List(config.EvaluationsCalendarID).
		SingleEvents(true).
		// MaxResults(1).
		OrderBy("startTime").
		TimeMin(startString).
		TimeMax(endString).
		Q(teamName).
		Do()
	if err != nil {
		return nil, err
	}

	if len(slots.Items) == 0 {
		return nil, nil
	}

	return slots, nil
}

// CalendarAllTeamSlot func
func CalendarAllTeamSlot(teamName string) (*calendar.Events, error) {
	service, err := calendarService()
	if err != nil {
		return nil, err
	}

	slots, err := service.Events.
		List(config.EvaluationsCalendarID).
		SingleEvents(true).
		// MaxResults(1).
		OrderBy("startTime").
		TimeMin(config.EvaluationsWeekStart).
		TimeMax(config.EvaluationsWeekEnd).
		Q(teamName).
		Do()
	if err != nil {
		return nil, err
	}

	if len(slots.Items) == 0 {
		return nil, nil
	}

	return slots, nil
}

// CalendarUnReserveTeamSlot func
func CalendarUnReserveTeamSlot(teamName, slotID string) error {

	service, err := calendarService()
	if err != nil {
		return err
	}
	oldSlot, err := service.Events.Get(config.EvaluationsCalendarID, slotID).Do()
	if err != nil {
		return err
	}

		oldSlotID := oldSlot.Id

		oldSlot.Summary = strings.Replace(oldSlot.Summary, teamName, "FREE", 1)
		// oldSlot.ColorId = "0"
		// oldSlot = &calendar.Event{
		// 	Summary: "FREE",
		// 	ColorId: "0",
		// }
		_, err = service.Events.Patch(config.EvaluationsCalendarID, oldSlotID, oldSlot).Do()
		return err


}

// CalendarReserveTeamSlot func
func CalendarReserveTeamSlot(teamName, slotID string) error {
	service, err := calendarService()
	if err != nil {
		return err
	}

	newSlot, err := service.Events.Get(config.EvaluationsCalendarID, slotID).Do()
	if err != nil {
		return err
	}

	if !strings.Contains(newSlot.Summary, "FREE") {
		return fmt.Errorf("slot completely reserved")
	}

	if strings.Contains(newSlot.Summary, teamName) {
		return fmt.Errorf("slot already reserved")
	}

	slots,err := CalendarAllTeamSlotsInWeek(teamName,newSlot)// CalendarAllTeamSlot(teamName)
	teamTotalReservedTime := 0.0
	if(slots != nil){
		for _,e := range slots.Items {

			endTime, _ := time.Parse(time.RFC3339, e.End.DateTime)
			startTime, _ := time.Parse(time.RFC3339, e.Start.DateTime)
			teamTotalReservedTime += (endTime.Sub(startTime)).Minutes()
		}
	}
	endTime, _ := time.Parse(time.RFC3339, newSlot.End.DateTime)
	startTime, _ := time.Parse(time.RFC3339, newSlot.Start.DateTime)

	max, _ := strconv.ParseFloat(config.BscVRWeeklyMinutes, 64)
	if teamTotalReservedTime + (endTime.Sub(startTime)).Minutes() > max {
		return fmt.Errorf("You already resevred a total of %.0f minutes. Adding this slot will exceed your maximum allocated %.0f minutes",teamTotalReservedTime, max)
	}
	log.Println("Total Slots: ",teamTotalReservedTime)


	// oldSlot, _ := CalendarTeamSlot(teamName)	
	
	newSlot.Summary = strings.Replace(newSlot.Summary, "FREE", teamName, 1)
	// newSlot.ColorId = "2"
	// newSlot = &calendar.Event{
	// 	Summary: teamName,
	// 	ColorId: "5",
	// }
	if _, err := service.Events.Patch(config.EvaluationsCalendarID, slotID, newSlot).Do(); err != nil {
		return err
	}

	// if oldSlot != nil {
	// 	oldSlotID := oldSlot.Id

	// 	oldSlot.Summary = strings.Replace(oldSlot.Summary, teamName, "FREE", 1)
	// 	oldSlot.ColorId = "0"
	// 	// oldSlot = &calendar.Event{
	// 	// 	Summary: "FREE",
	// 	// 	ColorId: "0",
	// 	// }
	// 	_, err = service.Events.Patch(config.EvaluationsCalendarID, oldSlotID, oldSlot).Do()
	// 	return err
	// }

	return nil
}
