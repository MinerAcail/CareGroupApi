package schemas

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/kobbi/vbciapi/graph/model"
	"github.com/kobbi/vbciapi/jwt/middleware"
	"github.com/lib/pq"
	"github.com/tealeg/xlsx"
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

	// if err := db.Where("types = ? AND sub_church_id = ?", "SubLeader", subChurchID).Find(&leaders).Error; err != nil {
	// 	// Handle errors here, such as if no leaders with the same day are found.
	// 	return err
	// }
	if err := db.Where("sub_church_id = ? AND types && ?", subChurchID, pq.Array([]string{"SubLeader"})).Find(&leaders).Error; err != nil {
		// Handle the error, such as logging or returning an error response.
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

	var leadersWithSameDay []model.Member
	if err := db.Where("types IN (?) AND day = ? AND sub_church_id = ?", pq.Array([]string{"SubLeader"}), day, subChurchID).Find(&leadersWithSameDay).Error; err != nil {
		// Handle errors here, such as if no leaders with the same day are found.
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
func GetMigrationRequestByID(db *gorm.DB, requestID string) (*model.MigrationRequest, error) {
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
func FindSubByID(db *gorm.DB, subChurchID string) ([]model.Member, error) {
	var leadersWithSameDay []model.Member
	// if err := db.Where("types IN (?) AND sub_church_id = ?", []string{"Leader", "SubLeader"}, subChurchID).Find(&leadersWithSameDay).Error; err != nil {
	// 	// Handle errors here, such as if no leaders with the same day are found.
	// 	return nil, err
	// }
	if err := db.Where("sub_church_id = ? AND types && ?", subChurchID, pq.Array([]string{"Leader", "SubLeader"})).Find(&leadersWithSameDay).Error; err != nil {
		// Handle the error, such as logging or returning an error response.
		return nil, err
	}
	return leadersWithSameDay, nil
}
func HandleMissingLeaders(db *gorm.DB, subChurchID string, missingDays []string) error {
	missingLeaders := map[string]*model.Member{}
	LeaderType := "Leader"

	for _, day := range missingDays {
		// Create a new leader for the missing day
		newLeader := &model.Member{
			Day:   day,
			Types: []string{LeaderType},
			// Set the appropriate type for leaders
			SubChurchID: &subChurchID,
		}
		missingLeaders[day] = newLeader
	}
	var members []model.Member

	for _, day := range missingDays {
		leader := missingLeaders[day]
		membersWithMissingLeader := make([]model.Member, 0)
		leaderID := leader.ID.String()

		// Loop through members and assign the missing leader
		for _, member := range members {
			if member.Day == day {
				member.LeaderID = &leaderID
				membersWithMissingLeader = append(membersWithMissingLeader, member)
			}
		}

		// Save the leader and members to the database
		if err := db.Create(leader).Error; err != nil {
			return err
		}
		if err := db.Create(&membersWithMissingLeader).Error; err != nil {
			return err
		}
	}

	return nil
}

func CellMatches(cellValue string, expectedColumn string, synonyms []string) bool {
	// Create a regex pattern to match the expected column name and its synonyms, accounting for variations in whitespace and case-insensitivity
	pattern := "(?i)\\s*(" + expectedColumn + "|" + strings.Join(synonyms, "|") + ")\\s*"
	regex := regexp.MustCompile(pattern)

	// Check if the cell value matches the pattern
	return regex.MatchString(cellValue)
}

func ConvertExcelToCSV(fileDataString string) (string, error) {
	xlFile, err := xlsx.OpenBinary([]byte(fileDataString))
	if err != nil {
		return "", err
	}

	var csvData [][]string
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			var csvRow []string
			for _, cell := range row.Cells {
				csvRow = append(csvRow, cell.String())
			}
			csvData = append(csvData, csvRow)
		}
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "csvtempdir")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file within the temporary directory
	tmpFile, err := os.CreateTemp(tempDir, "temp.csv")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	csvWriter := csv.NewWriter(tmpFile)
	for _, row := range csvData {
		csvWriter.Write(row)
	}
	csvWriter.Flush()

	return tmpFile.Name(), nil
}

func ReadMembersFromCSV(ctx context.Context, csvFilePath string, churchID *string, db *gorm.DB) ([]*model.Member, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	// Define variables to store members and duplicate records
	var members []*model.Member
	var duplicateRecords []string

	var leaders []model.Member

	if err := db.Where("sub_church_id = ? AND types && ?", churchID, pq.Array([]string{"Leader", "SubLeader"})).Find(&leaders).Error; err != nil {
		// Handle the error, such as logging or returning an error response.
		return nil, fmt.Errorf("leader not found")
	}
	if len(leaders) >= 2 {
		var leadersMap = make(map[string]*model.Member)

		for _, leader := range leaders {
			leadersMap[leader.Day] = &leader
		}

		// Ensure that there is a leader for each day of the week (Monday to Sunday)
		for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
			if leadersMap[day] == nil {
				return nil, fmt.Errorf(" No leader assigned for %s", day)
			}
		}
	}

	// Read and process the CSV data line by line
	for rowIndex := 0; ; rowIndex++ {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if rowIndex == 0 {
			// Skip the header row
			continue
		}

		if len(record) < 5 {
			// Handle the case where a row does not have enough columns
			return nil, fmt.Errorf(" Row %d has insufficient columns", rowIndex)
		}

		// Extract data from the CSV record
		name := record[0]
		email := record[1]
		day := record[2]
		location := record[3]
		phoneNumber := record[4]

		var lop *model.Member
		if len(leaders) == 1 {
			// // Get the leader for the given churchID
			leader, err := GetLeaderByChurchID(db, churchID)
			if err != nil {
				return nil, err
			}
			lop = leader

		} else if len(leaders) > 1 {

			// // Get the leader for the given churchID
			leader, err := FindLeaderWithSameDay(ctx, db, day, *churchID)
			if err != nil {
				return nil, err
			}
			lop = leader
		}

		leaderID := lop.ID.String()

		if leaderID == "" {
			// Handle the case where lop is nil, e.g., return an error or take appropriate action.
			return nil, fmt.Errorf("leaderID is nil")
		}

		// Clean phone number
		phoneNumber = CleanPhoneNumber(phoneNumber)

		// Create a member object
		member := &model.Member{
			Name:        name,
			Email:       email,
			PhoneNumber: &phoneNumber,
			Day:         day,
			Location:    &location,
			SubChurchID: churchID,
		}

		// Call the checking Duplicate Records function if already in the database
		if err := CheckMembersDuplicateRecords(db, name, email, phoneNumber); err != nil {
			// Collect the record with duplicate data
			duplicateRecords = append(duplicateRecords, fmt.Sprintf("Row %d - Name: %s, Email: %s, Phone: %s", rowIndex, name, email, phoneNumber))

			// Skip creating this member and continue to the next iteration
			continue
		}

		members = append(members, member)
	}

	if len(duplicateRecords) > 0 {
		errMessage := "Duplicate records found:\n" + strings.Join(duplicateRecords, "\n")
		return nil, fmt.Errorf(errMessage)
	}

	return members, nil
}

func ProcessCSVFile(ctx context.Context, csvFilePath string, churchID string, db *gorm.DB) ([]*model.Member, error) {
	// Open the CSV file for reading.
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close()

	// Create a CSV reader to read the file.
	csvReader := csv.NewReader(file)
	// Read all records from the CSV.
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	columnSynonyms := map[string][]string{
		"Name": {
			"Full Name", "Full_Name", "First Name", "MEMBERS NAME",
			"members Name", "MEMBERS NAME", "NAMES", " NAME", "Name", "Names",
		},
		"Email": {
			"Email", "Email Address", "E-mail", "Emails", "Email",
		},
		"Phone Number": {
			"Phone", "TELEPHONE NUMBER", "Telephone Number", "TELEPHONE NUMBER",
			"Contact", "Contacts", "numbers", "number",
			"members number", "members numbers", "member numbers",
		},
		"Day": {
			"Day born", "Schedule Day", "Day", "Days", "Day borns",
		},
		"Location": {
			"Locations", "Location", "LOCATION",
		},
	}

	// Determine the indices of the header columns by name.
	header := records[0]
	nameIndex := findIndex(header, "Name", columnSynonyms)
	emailIndex := findIndex(header, "Email", columnSynonyms)
	locationIndex := findIndex(header, "Location", columnSynonyms)
	dayIndex := findIndex(header, "Day", columnSynonyms)
	phoneNumberIndex := findIndex(header, "Phone Number", columnSynonyms)

	// Read and process CSV data.
	var members []*model.Member

	for rowNumber, record := range records {
		// Skip the header row (if it exists).
		if rowNumber == 0 {
			continue
		}

		// Assuming the CSV structure is: name, email, phoneNumber, day, location
		if len(record) != 5 {
			return nil, fmt.Errorf(" Invalid CSV format at row %d", rowNumber)
		}

		name := record[nameIndex]
		email := record[emailIndex]
		location := record[locationIndex]
		day := record[dayIndex]
		phoneNumber := record[phoneNumberIndex]

		var leaders []model.Member

		if err := db.Where("sub_church_id = ? AND types && ?", churchID, pq.Array([]string{"Leader", "SubLeader"})).Find(&leaders).Error; err != nil {
			// Handle the error, such as logging or returning an error response.
			return nil, err
		}

		var lop *model.Member
		if len(leaders) == 1 {
			// // Get the leader for the given churchID
			leader, err := GetLeaderByChurchID(db, &churchID)
			if err != nil {
				return nil, err
			}
			lop = leader

		} else if len(leaders) <= 6 {
			// // Get the leader for the given churchID
			leader, err := GetLeaderByChurchID(db, &churchID)
			if err != nil {
				return nil, err
			}
			lop = leader
		}else if len(leaders) >= 7 {
			fmt.Printf("day of leaders: %s\n", day)
			fmt.Printf("location of leaders: %s\n", location)
			fmt.Printf("churchID of leaders: %s\n", churchID)

			// // Get the leader for the given churchID
			leader, err := FindLeaderWithSameDay(ctx, db, day, churchID)
			if err != nil {
				return nil, err
			}
			lop = leader
		}

		leaderID := lop.ID.String()

		// Create a Member instance.
		member := &model.Member{
			Name:        name,
			Email:       email,
			PhoneNumber: &phoneNumber,
			Day:         day,
			Location:    &location,
			LeaderID:    &leaderID, // Assign the LeaderID to the selected leader's ID

			SubChurchID: &churchID,
		}
		// Create the member in the database.
		if err := db.Create(member).Error; err != nil {
			return nil, err
		}

		// Append the member to the list.
		members = append(members, member)
	}

	// Return the list of imported members.
	return members, nil
}
func findIndex(header []string, target string, columnSynonyms map[string][]string) int {
    // Try to find an exact match first.
    for i, name := range header {
        if name == target {
            return i
        }
    }

    // If an exact match is not found, look for synonyms.
    for i, name := range header {
        for _, synonym := range columnSynonyms[target] {
            if name == synonym {
                return i
            }
        }
    }

    return -1 // Not found
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
