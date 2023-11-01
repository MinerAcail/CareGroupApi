package schemas

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/kobbi/vbciapi/graph/model"
	"github.com/kobbi/vbciapi/jwt/middleware"
	"github.com/lib/pq"
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

// func GetSortingCondition(sortInput *model.SortInput) (string, error) {
// 	// Check for valid sorting fields to avoid SQL injection or incorrect sorting
// 	validSortingFields := []string{"id", "name", "email", "phone_number", "day", "types", "created_at", "updated_at"}

// 	if sortInput == nil {
// 		return "", nil
// 	}
// 	if sortInput.Field == "" {
// 		sortInput.Field = "name"
// 	}
// 	if sortInput.Order == "" {
// 		return "", errors.New(" Invalid sorting input. Both 'field' and 'order' must be provided. ")
// 	}

// 	field := strings.ToLower(sortInput.Field)
// 	order := strings.ToUpper(sortInput.Order)

// 	if order != "ASC" && order != "DESC" {
// 		return "", errors.New(" Invalid sorting order. Must be either 'ASC' or 'DESC'. ")
// 	}

// 	//CheckIf field contain any of the list Strings in validSortingFields
// 	if !Contains(validSortingFields, field) {
// 		return "", errors.New(" Invalid sorting field")
// 	}
// 	// Sort the leaders based on their names(sort.Field)
// 	// Sort the leaders in ascending or descending order(order)
// 	query := fmt.Sprintf("%s %s", field, order)
// 	//the fmt.Sprintf add 2 string to 1 (sort.Field, order)=>(name DESC)

// 	return query, nil
// }

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

	// Check if a member with the same name and phone number already exists
	err := db.Where("name = ? AND phone_number = ?", inputName, inputPhoneNumber).First(&model).Error
	if err == nil {
		return fmt.Errorf("a member with the same name and phone number already exists")
	}
	// else if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return err
	// }

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
func CheckMembersDuplicateRecords(db *gorm.DB, inputName, inputEmail, inputPhoneNumber string) error {
	member := &model.Member{}

	// Check if a member with the same name and phone number already exists
	err := db.Where("name = ? AND phone_number = ?", inputName, inputPhoneNumber).First(member).Error
	if err == nil {
		return fmt.Errorf("a Member with the same name and phone number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if an email already exists
	err = db.Where("email = ?", inputEmail).First(member).Error
	if err == nil {
		return fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err

	}

	// Check if a phone_number already exists
	err = db.Where("phone_number = ?", inputPhoneNumber).First(member).Error
	if err == nil {
		return fmt.Errorf("phone_number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

// Function to calculate and update ReferenceIDCount for leaders.
func UpdateReferenceIDCounts(db *gorm.DB, subChurchID string) error {
	// Query all leaders
	var leaders []model.Member

	// if err := db.Where("types IN (?) ", []string{"Leader", "SubLeader"}).Find(&leaders).Error; err != nil {
	if err := db.Where("types IN (?)", pq.Array([]string{"Leader", "SubLeader"})).Find(&leaders).Error; err != nil {
		return err
	}

	// Loop through leaders and calculate their ReferenceIDCount
	for i := range leaders {
		var count int64 // Change the type to int64

		// Query the number of members associated with this leader
		if err := db.Model(&model.Member{}).
			Where("leader_id = ?", leaders[i].ID).
			Count(&count).
			Error; err != nil {
			return err
		}

		// Update the ReferenceIDCount field
		leaders[i].ReferenceIDCount = (*int)(unsafe.Pointer(&count)) // Convert int64 to *int
		if err := db.Save(&leaders[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
func FindLeaderByID(db *gorm.DB, leaderID string) (*model.Member, error) {
	var leader model.Member
	if err := db.Where("id = ?", leaderID).First(&leader).Error; err != nil {
		// Handle errors here, such as if the leader is not found.
		return nil, err
	}
	return &leader, nil
}

func FindLeaderWithSameDay(ctx context.Context, db *gorm.DB, day string, subChurchID string) (*model.Member, error) {
	// Call UpdateReferenceIDCounts to ensure the counts are up-to-date.
	leaderID, ok := ctx.Value(middleware.IDContextKey).(string)
	if !ok {
		return nil, fmt.Errorf("leaderID not found in request context")
	}
	if err := UpdateReferenceIDCounts(db, subChurchID); err != nil {
		// Handle the error, such as logging or returning an error response.
		return nil, err
	}
	// Define a variable to hold the selected leader with the same day.
	var selectedLeader model.Member

	// Perform a database query to find leaders with the specified day.
	var leadersWithSameDay []model.Member
	

	if err := db.Where("sub_church_id = ? AND types && ?", subChurchID, pq.Array([]string{"SubLeader"})).Find(&leadersWithSameDay).Error; err != nil {
		// Handle the error, such as logging or returning an error response.
		return nil, err
	}

	if len(leadersWithSameDay) == 0 {
		// No leaders found with the same day, return the leader identified by leaderID.
		return FindLeaderByID(db, leaderID) // Replace FindLeaderByID with the actual function to retrieve a leader by ID.
	}

	// Sort the leaders first by ReferenceIDCount in ascending order.
	sort.SliceStable(leadersWithSameDay, func(i, j int) bool {
		// Dereference the pointers and then compare the integer values.
		return *leadersWithSameDay[i].ReferenceIDCount < *leadersWithSameDay[j].ReferenceIDCount
	})

	// Now, sort them by lexicographically smallest name.
	minReferenceIDCount := leadersWithSameDay[0].ReferenceIDCount
	potentialLeaders := make([]model.Member, 0)
	for _, leader := range leadersWithSameDay {
		if leader.ReferenceIDCount == minReferenceIDCount {
			potentialLeaders = append(potentialLeaders, leader)
		} else {
			// Since the leaders are sorted by ReferenceIDCount, we can break the loop when we find a leader with a higher count.
			break
		}
	}

	// If there are potential leaders with the same lowest ReferenceIDCount, sort them by lexicographically smallest name.
	if len(potentialLeaders) > 0 {
		selectedLeader = potentialLeaders[0]
		for _, leader := range potentialLeaders {
			if leader.Name < selectedLeader.Name {
				selectedLeader = leader
			}
		}
	}

	// Return the selected leader with the same day.
	return &selectedLeader, nil
}
func GetLeaderByChurchID(db *gorm.DB, churchID *string) (*model.Member, error) {
	var leader model.Member

	// if err := db.Where("types IN (?) AND sub_church_id = ?",  pq.Array([]string{"Leader", "SubLeader"}), churchID).First(&leader).Error; err != nil {
	// 	return nil, err
	// fixing array query
	// }
	if err := db.Where("sub_church_id = ? AND types && ?", churchID, pq.Array([]string{"Leader"})).Find(&leader).Error; err != nil {
		// Handle the error, such as logging or returning an error response.
		return nil, err
	}
	// if err := r.DB.Where("sub_church_id = ? AND types && ?", churchID, pq.Array([]string{"Leader"})).Find(&leader).Error; err != nil {
	// 	// Handle the error, such as logging or returning an error response.
	// 	return nil, err
	// }

	// // Check if leader.ID is zero (assuming it's of type [16]byte)
	// if leader.ID == [16]byte{} {
	// 	// Handle the case where leader is not found, e.g., return an error or take appropriate action.
	// 	return nil, fmt.Errorf("leader is nil")
	// }

	return &leader, nil
}

// GetMigrationRequestByID retrieves a MigrationRequest by its ID.
func  GetMigrationRequestByID(db *gorm.DB, requestID string) (*model.MigrationRequest, error) {
    var migrationRequest model.MigrationRequest

    // Use GORM to find the migration request by its ID
    if err := db.Where("id = ?", requestID).First(&migrationRequest).Error; err != nil {
        
        return nil, fmt.Errorf(" Migration request not found")
    }

    return &migrationRequest, nil
}

func GetSubChurchIDForLeader(db *gorm.DB, leaderID string) (string, error) {
	var subChurchID string

	// Assuming you have a 'leaders' table with a 'sub_church_id' column
	err := db.Model(&model.Member{}).
		Where("id = ?", leaderID).
		Pluck("sub_church_id", &subChurchID).Error

	if err != nil {
		return "", err
	}

	return subChurchID, nil
}

func CleanPhoneNumber(phoneNumber string) string {
	cleanedPhoneNumber := strings.ReplaceAll(phoneNumber, " ", "")
	cleanedPhoneNumber = strings.TrimPrefix(cleanedPhoneNumber, "+233")
	cleanedPhoneNumber = strings.TrimPrefix(cleanedPhoneNumber, "233")

	// Check if the cleaned phone number has fewer than 10 digits
	if len(cleanedPhoneNumber) < 10 {
		// Add a leading "0" to make it 10 digits
		cleanedPhoneNumber = "0" + cleanedPhoneNumber
	}

	return cleanedPhoneNumber
}

var (
	ErrInvalidInput = errors.New("invalid input data")
	ErrDatabase     = errors.New("database error")
)

func GetWeekNumber(date time.Time) int {
	_, weekNumber := date.ISOWeek()
	return weekNumber
}
