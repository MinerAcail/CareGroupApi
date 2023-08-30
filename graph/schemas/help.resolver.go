package schemas

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kobbi/vbciapi/graph/model"
	"gorm.io/gorm"
)

// Utility function to check if a string is present in a slice of strings.
func Contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
func GetSortingCondition(sortInput *model.SortInput) (string, error) {
	// Check for valid sorting fields to avoid SQL injection or incorrect sorting
	validSortingFields := []string{"id", "name", "email", "phone_number", "day", "types", "created_at", "updated_at"}

	if sortInput == nil {
		return "", nil
	}
	if sortInput.Field == "" {
		sortInput.Field = "name"
	}
	if sortInput.Order == "" {
		return "", errors.New(" Invalid sorting input. Both 'field' and 'order' must be provided. ")
	}

	field := strings.ToLower(sortInput.Field)
	order := strings.ToUpper(sortInput.Order)

	if order != "ASC" && order != "DESC" {
		return "", errors.New(" Invalid sorting order. Must be either 'ASC' or 'DESC'. ")
	}

	//CheckIf field contain any of the list Strings in validSortingFields
	if !Contains(validSortingFields, field) {
		return "", errors.New(" Invalid sorting field")
	}
	// Sort the leaders based on their names(sort.Field)
	// Sort the leaders in ascending or descending order(order)
	query := fmt.Sprintf("%s %s", field, order)
	//the fmt.Sprintf add 2 string to 1 (sort.Field, order)=>(name DESC)

	return query, nil
}

/*
order := "ASC"

	if strings.ToUpper(sortInput.Order) == "DESC" {
		order = "DESC"
	}
*/
func IsValidGroupField(field string) bool {
	// List of valid fields for grouping in the Leader model
	validGroupFields := []string{"name", "email", "phoneNumber", "day", "types", "location"}

	// Check if the provided field is in the list of valid group fields
	for _, validField := range validGroupFields {
		if field == validField {
			return true
		}
	}

	return false
}

func ShuffleLeaders(leaders []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(leaders) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		leaders[i], leaders[j] = leaders[j], leaders[i]
	}
}

func CheckDuplicateRecords(db *gorm.DB, model interface{}, inputName, inputEmail, inputPhoneNumber string) error {

	// Check if a student with the same name and phone number already exists
	err := db.Where("name = ? AND phone_number = ?", inputName, inputPhoneNumber).First(&model).Error
	if err == nil {
		return fmt.Errorf("a student with the same name and phone number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if an email already exists
	err = db.Where("email = ?", inputEmail).First(&model).Error
	if err == nil {
		return fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if a phone_number already exists
	err = db.Where("phone_number = ?", inputPhoneNumber).First(&model).Error
	if err == nil {
		return fmt.Errorf("phone_number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrDatabase     = errors.New("database error")
)

func GetWeekNumber(date time.Time) int {
	_, weekNumber := date.ISOWeek()
	return weekNumber
}
