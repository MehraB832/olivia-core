package modules

import (
	"fmt"
	"strings"

	"github.com/MehraB832/olivia_core/language"

	"github.com/MehraB832/olivia_core/util"

	"github.com/MehraB832/olivia_core/user"

	"github.com/MehraB832/olivia_core/language/date"
)

var (
	// ReminderSetterTag is the intent tag for its module
	ReminderSetterTag = "reminder setter"
	// ReminderGetterTag is the intent tag for its module
	ReminderGetterTag = "reminder getter"
)

// ReminderSetterReplacer replaces the pattern contained inside the response by the date of the reminder
// and its reason.
// See modules/modules.go#Module.Replacer() for more details.
func ReminderSetterReplacer(locale, entry, response, token string) (string, string) {
	// Search the time and
	sentence, date := date.SearchTime(locale, entry)
	reason := language.SearchReason(locale, sentence)

	// Format the date
	formattedDate := date.Format("01/02/2006 03:04")

	// Add the reminder inside the user's information
	user.UpdateUserProfile(token, func(information user.UserProfile) user.UserProfile {
		information.ImportantDates = append(information.ImportantDates, user.UserReminder{
			ReminderDetails: reason,
			ReminderDate:   formattedDate,
		})

		return information
	})

	return ReminderSetterTag, fmt.Sprintf(response, reason, formattedDate)
}

// ReminderGetterReplacer gets the reminders in the user's information and replaces the pattern in the
// response patterns by the current reminders
// See modules/modules.go#Module.Replacer() for more details.
func ReminderGetterReplacer(locale, _, response, token string) (string, string) {
	reminders := user.RetrieveUserProfile(token).ImportantDates
	var formattedReminders []string

	// Iterate through the reminders and parse them
	for _, reminder := range reminders {
		formattedReminder := fmt.Sprintf(
			util.SelectRandomMessage(locale, "reminder"),
			reminder.ReminderDetails,
			reminder.ReminderDate,
		)
		formattedReminders = append(formattedReminders, formattedReminder)
	}

	// If no reminder has been found
	if len(formattedReminders) == 0 {
		return ReminderGetterTag, util.SelectRandomMessage(locale, "no reminders")
	}

	return ReminderGetterTag, fmt.Sprintf(response, strings.Join(formattedReminders, " "))
}
