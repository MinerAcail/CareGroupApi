package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kobbi/vbciapi/graph/model"
	"github.com/kobbi/vbciapi/graph/schemas"
	"github.com/kobbi/vbciapi/jwt/helpers"
	"github.com/kobbi/vbciapi/jwt/middleware"
	"gorm.io/gorm"
)

// Aggregations is the resolver for the Aggregations field.
func (r *leaderResolver) Aggregations(ctx context.Context, obj *model.Leader) (*model.LeaderAggregations, error) {
	// Calculate the aggregations for the leader

	totalStudents := len(obj.Students)

	// Calculate max students per leader
	maxStudentsPerLeader := 0
	for _, student := range obj.Students {
		if len(student.Registrations) > maxStudentsPerLeader {
			maxStudentsPerLeader = len(student.Registrations)
		}
	}

	return &model.LeaderAggregations{
		TotalLeaders:             1, // Since we're aggregating for one leader at a time
		AverageStudentsPerLeader: float64(totalStudents) / float64(1),
		MaxStudentsPerLeader:     maxStudentsPerLeader,
	}, nil
}

// ID is the resolver for the id field.
func (r *leaderResolver) ID(ctx context.Context, obj *model.Leader) (string, error) {
	id := obj.ID.String()
	return id, nil
}

// CallerCentarUpdateRegister is the resolver for the CallerCentarUpdateRegister field.
func (r *mutationResolver) CallerCentarUpdateRegister(ctx context.Context, input model.CreateRegistrationInput, leaderIDs []*string) ([]*model.Registration, error) {
	// Retrieve the current week's registrations
	/* currentWeekRegistrations, err := r.getCurrentWeekRegistrations()
	if err != nil {
		return nil, err
	}
	for _, reg := range currentWeekRegistrations {
		log.Printf("Registration ID: %s, CallerID: %v", reg.ID, reg.CallerID)
	} */
	// Shuffle the registrations
	// shuffledRegistrations := shuffleRegistrations(currentWeekRegistrations)

	updatedRegistrations := []*model.Registration{}

	/* // Distribute shuffled registrations to leaders
	for _, leaderID := range leaderIDs {
		leader, err := r.getLeaderByID(*leaderID)
		if err != nil {
			return nil, err
		}

		// Calculate how many registrations this leader should receive
		registrationsPerLeader := len(shuffledRegistrations) / len(leaderIDs)
		remainingRegistrations := len(shuffledRegistrations) % len(leaderIDs)

		// Calculate the range of registrations for this leader
		startIdx := len(updatedRegistrations)
		endIdx := startIdx + registrationsPerLeader
		if len(updatedRegistrations)+registrationsPerLeader > len(shuffledRegistrations) {
			endIdx += remainingRegistrations
		}

		// Assign the registrations to this leader
		for j := startIdx; j < endIdx; j++ {
			shuffledRegistrations[j].CallerID = append(shuffledRegistrations[j].CallerID, *leaderID)
			updatedRegistrations = append(updatedRegistrations, shuffledRegistrations[j])
		}
	}

	// Save the distributed registrations to the database
	if err := r.DB.Save(updatedRegistrations).Error; err != nil {
		return nil, err
	} */

	// Return the list of updated registrations
	return updatedRegistrations, nil
}

// CreateLeader is the resolver for the createLeader field.
func (r *mutationResolver) CreateLeader(ctx context.Context, input model.CreateLeaderInput) (*model.Leader, error) {
	// Check if  email already exists
	Leaders := &model.Leader{}
	// Call the checking Duplicate Records function if already in db
	if err := schemas.CheckDuplicateRecords(r.DB, Leaders, input.Name, *input.Email, input.PhoneNumber); err != nil {
		return nil, err
	}

	// Hashing Password to the DB  using the input Passowrd

	password, err := helpers.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", schemas.ErrInvalidInput)
	}

	// Generate and assign tokens to the Leader object

	leader := &model.Leader{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Day:         input.Day,
		Password:    password,
		Types:       input.Types,
	}

	// Save the leader to the database using your preferred ORM
	if err := r.DB.Create(leader).Error; err != nil {
		return nil, fmt.Errorf("failed to save leader to the database: %w", schemas.ErrDatabase)
	}

	token, err := helpers.GenerateToken(input.Types, (leader.ID).String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update the leader's token with the generated token
	leader.Token = token

	// Save the updated leader with the token to the database
	if err := r.DB.Save(leader).Error; err != nil {
		return nil, fmt.Errorf("failed to save leader's token to the database: %w", err)
	}

	return leader, nil
}

// CreateRegistration is the resolver for the createRegistration field.
func (r *mutationResolver) CreateRegistration(ctx context.Context, studentID string, input model.CreateRegistrationInput) (*model.Registration, error) {
	// Check if studentID is a valid UUID
	_, err := uuid.Parse(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid studentID: %w", err)
	}
	// Create a new Registration object using the input data
	registration := &model.Registration{
		LastComment: *input.LastComment,
		Absence:     *input.Absence,
		Present:     *input.Present,
		StudentID:   studentID,
	}

	// Save the registration to the database
	if err := r.DB.Create(registration).Error; err != nil {
		return nil, err
	}

	// Return the created registration
	return registration, nil
}

// CreateStudent is the resolver for the createStudent field.
func (r *mutationResolver) CreateStudent(ctx context.Context, input model.CreateStudentInput) (*model.Student, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	leaderID, ok := ctx.Value(middleware.LeaderIDContextKey).(string)
	if !ok {
		return nil, fmt.Errorf("LeaderID not found in request context")
	}
	students := &model.Student{}
	// Call the checking Duplicate Records function if already in db
	if err := schemas.CheckDuplicateRecords(r.DB, students, input.Name, *input.Email, input.PhoneNumber); err != nil {
		return nil, err
	}

	student := &model.Student{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Day:         input.Day,
		LeaderID:    leaderID, // Set the LeaderID field using the extracted value.
	}

	err := r.DB.Create(student).Error
	if err != nil {
		return nil, err
	}

	// Eager load the associated Leader using Preload
	err = r.DB.Model(student).Preload("Registrations").Find(student).Error
	if err != nil {
		return nil, err
	}

	return student, nil
}

// DeleteStudent is the resolver for the deleteStudent field.
func (r *mutationResolver) DeleteStudent(ctx context.Context, studentID string) (bool, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return false, err
	}

	// Retrieve the student from the database based on the provided studentID
	student := &model.Student{}
	if err := r.DB.Where("id = ?", studentID).First(student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("student with ID %s not found", studentID)
		}
		return false, err
	}

	// Delete the student from the database
	if err := r.DB.Delete(student).Error; err != nil {
		return false, err
	}

	// Return success
	return true, nil
}

// DaleteRegistration is the resolver for the daleteRegistration field.
func (r *mutationResolver) DaleteRegistration(ctx context.Context, registrationID string) (bool, error) {
	registration := &model.Registration{}

	if err := r.DB.Where("id = ?", registrationID).First(&registration).Error; err != nil {
		return false, err
	}

	// Delete the registration from the database
	if err := r.DB.Delete(registration).Error; err != nil {
		return false, err
	}
	return true, nil
}

// DistributeRegistrationsToLeaders is the resolver for the distributeRegistrationsToLeaders field.
func (r *mutationResolver) DistributeRegistrationsToLeaders(ctx context.Context, leaderIds []string) ([]*model.LeaderRegistrationsDistribution, error) {
	// Fetch the current week's registrations
	currentWeekRegistrations, err := r.Query().CurrentWeekRegistrations(ctx)
	if err != nil {
		return nil, err
	}

	// Shuffle the leader IDs using your shuffling algorithm
	schemas.ShuffleLeaders(leaderIds)

	// Create a map to store registration IDs for each leader
	leaderRegistrationMap := make(map[string][]string)

	// Distribute registration IDs among shuffled leader IDs
	for i, reg := range currentWeekRegistrations {
		leaderIndex := i % len(leaderIds)
		leaderID := leaderIds[leaderIndex]
		leaderRegistrationMap[leaderID] = append(leaderRegistrationMap[leaderID], reg.ID.String())

		// Update the registration's leader_id
		reg.LeaderID = leaderID // Update this line with the correct field name from your Registration model
		if err := r.DB.Save(&reg).Error; err != nil {
			return nil, err
		}
	}

	// Create the final response structure
	var leaderDistribution []*model.LeaderRegistrationsDistribution
	for leaderID, regIDs := range leaderRegistrationMap {
		leaderDist := &model.LeaderRegistrationsDistribution{
			LeaderID:        leaderID,
			RegistrationIDs: regIDs,
		}
		leaderDistribution = append(leaderDistribution, leaderDist)
	}

	return leaderDistribution, nil
}

// LoginLeader is the resolver for the loginLeader field.
func (r *mutationResolver) LoginLeader(ctx context.Context, input model.LoginLeaderInput) (*model.Leader, error) {
	// Fetch the leader from the database based on the provided email
	leader := &model.Leader{}
	if err := r.DB.Where("phone_number = ?", input.PhoneNumber).First(leader).Error; err != nil {
		return nil, fmt.Errorf("leader not found")
	}

	// Verify the provided password against the hashed password in the database
	if err := helpers.VerifyPassword(leader.Password, input.Password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	// Generate a new token for the authenticated leader
	token, err := helpers.GenerateToken(leader.Types, leader.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	// Update the leader's token with the newly generated token
	leader.Token = token

	// Save the updated leader with the token to the database
	if err := r.DB.Save(leader).Error; err != nil {
		return nil, fmt.Errorf("failed to save leader's token to the database")
	}

	return leader, nil
}

// UpdateLeader is the resolver for the updateLeader field.
func (r *mutationResolver) UpdateLeader(ctx context.Context, input model.UpdateLeaderInput) (*model.Leader, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}

	leaderID, ok := ctx.Value(middleware.LeaderIDContextKey).(string)
	if !ok {
		return nil, fmt.Errorf("leaderID not found in request context")
	}

	// Retrieve the leader from the database based on the provided ID
	leader := &model.Leader{}
	if err := r.DB.Where("id = ?", leaderID).First(leader).Error; err != nil {
		return nil, err
	}
	// Call the checking Duplicate Records function if already in db
	if err := schemas.CheckDuplicateRecords(r.DB, leader, *input.Name, *input.Email, *input.PhoneNumber); err != nil {
		return nil, err
	}

	// Update the leader's fields with the input values
	if input.Name != nil {
		leader.Name = *input.Name
	}
	if input.Email != nil {
		leader.Email = input.Email
	}
	if input.PhoneNumber != nil {
		leader.PhoneNumber = *input.PhoneNumber
	}
	if input.Day != nil {
		leader.Day = *input.Day
	}
	if input.Types != nil {
		leader.Types = *input.Types
	}

	// Save the updated leader to the database
	if err := r.DB.Save(leader).Error; err != nil {
		return nil, err
	}

	// Return the updated leader
	return leader, nil
}

// UpdateRegistration is the resolver for the updateRegistration field.
func (r *mutationResolver) UpdateRegistration(ctx context.Context, input model.CreateRegistrationInput, registrationID string) (*model.Registration, error) {
	// Fetch the existing registration from the database
	registration := &model.Registration{}

	if err := r.DB.Where("id = ?", registrationID).First(&registration).Error; err != nil {
		return nil, err
	}

	// Update the fields with the new values
	if input.Absence != nil {
		registration.Absence = *input.Absence
	}
	if input.Present != nil {
		registration.Present = *input.Present
	}
	if input.LastComment != nil {
		registration.LastComment = *input.LastComment
	}

	// Save the updated registration to the database
	if err := r.DB.Save(registration).Error; err != nil {
		return nil, err
	}

	// Convert the GORM model to the GraphQL model and return
	return registration, nil
}

// UpdateRegistrationByLeader is the resolver for the updateRegistrationByLeader field.
func (r *mutationResolver) UpdateRegistrationByLeader(ctx context.Context, input model.CreateRegistrationInput, registrationID string, leaderID string) (*model.Registration, error) {
	register := &model.Registration{}

	if err := r.DB.Where("id =?", registrationID).First(&register).Error; err != nil {
		return nil, err
	}

	// Check if the registration is associated with the provided leader ID
	if register.LeaderID != leaderID {
		return nil, fmt.Errorf("registration is not associated with the provided leader")
	}

	// Update the registration's attributes based on the input
	if input.Absence != nil {
		register.Absence = *input.Absence
	}

	if input.Present != nil {
		register.Present = *input.Present
	}
	if input.LastComment != nil {
		register.LastComment = *input.LastComment
	}

	// Save the updated registration
	if err := r.DB.Save(register).Error; err != nil {
		return nil, err
	}

	return register, nil
}

// UpdateStudent is the resolver for the updateStudent field.
func (r *mutationResolver) UpdateStudent(ctx context.Context, input model.UpdateStudentInput, studentID string) (*model.Student, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}

	Student := &model.Student{}
	// Retrieve the Student from the database based on the provided ID
	if err := r.DB.Where("id = ?", studentID).First(Student).Error; err != nil {
		return nil, err
	}
	// Check if the provided email is already used by another student
	if err := r.DB.Where("email = ?", input.Email).Not("id = ?", studentID).First(&model.Student{}).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if the provided phone number is already used by another student
	if err := r.DB.Where("phone_number = ?", input.PhoneNumber).Not("id = ?", studentID).First(&model.Student{}).Error; err == nil {
		return nil, fmt.Errorf("phone number already exists")
	}
	// Update the Student's fields with the input values
	if input.Name != nil {
		Student.Name = *input.Name
	}
	if input.Email != nil {
		Student.Email = input.Email
	}
	if input.PhoneNumber != nil {
		Student.PhoneNumber = *input.PhoneNumber
	}
	if input.Day != nil {
		Student.Day = *input.Day
	}

	// Save the Student leader to the database
	if err := r.DB.Save(Student).Error; err != nil {
		return nil, err
	}

	// Return the updated Student
	return Student, nil
}

// GetgroupBy is the resolver for the GetgroupBy field.
func (r *queryResolver) GetgroupBy(ctx context.Context, groupByField string, tableName string) (*model.GroupByResults, error) {
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}
	var totalRecords int

	// Query the total count of records
	query := "SELECT COUNT(*) FROM " + tableName
	row := r.DB.Raw(query).Row()
	if err := row.Scan(&totalRecords); err != nil {
		return nil, err
	}

	var groupByResults []*model.GroupBy

	// Construct the raw SQL query
	query = fmt.Sprintf("SELECT %s, COUNT(*) FROM %s GROUP BY %s;", groupByField, tableName, groupByField)

	// Execute the query
	rows, err := r.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// Initialize variables to store values from the current row
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return nil, err
		}
		// Create a GroupBy object and append it to the results

		groupByResult := &model.GroupBy{
			Key:   &key,
			Count: &count,
		}
		groupByResults = append(groupByResults, groupByResult)
	}
	// Check for scanning errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Create the final GroupByResults object
	finalGroupByResults := &model.GroupByResults{
		Total:          &totalRecords,
		GroupByResults: groupByResults,
	}

	return finalGroupByResults, nil
}

// GetRegistrations is the resolver for the GetRegistrations field.
func (r *queryResolver) GetRegistrations(ctx context.Context) ([]*model.WeeklyResults, error) {
	var registrations []*model.Registration
	if err := r.DB.Find(&registrations).Error; err != nil {
		return nil, err
	}

	registrationsByWeek := make(map[int][]*model.Registration)

	for _, reg := range registrations {
		if reg.Absence { // Filter only registrations with absence true
			weekNumber := schemas.GetWeekNumber(reg.CreatedAt)
			registrationsByWeek[weekNumber] = append(registrationsByWeek[weekNumber], reg)
		}
	}

	var monthlyResults []*model.WeeklyResults
	for _, regs := range registrationsByWeek {
		year, month, day := regs[0].CreatedAt.Date()
		weekOfMonth := (day-1)/7 + 1
		monthStr := month.String() + " " + strconv.Itoa(year) // Get the month string
		monthlyResults = append(monthlyResults, &model.WeeklyResults{
			Date: &model.DateInfo{
				Month:       &monthStr, // Use a pointer to the month string
				WeekOfMonth: &weekOfMonth,
			},
			Registrations: regs,
		})
	}

	return monthlyResults, nil
}

// CurrentWeekRegistrations is the resolver for the currentWeekRegistrations field.
func (r *queryResolver) CurrentWeekRegistrations(ctx context.Context) ([]*model.Registration, error) {
	var registrations []*model.Registration
	if err := r.DB.Find(&registrations).Error; err != nil {
		return nil, err
	}
	currentWeekNumber := schemas.GetWeekNumber(time.Now()) // Get the current week number

	var currentWeekRegistrations []*model.Registration
	for _, reg := range registrations {
		if reg.Absence && schemas.GetWeekNumber(reg.CreatedAt) == currentWeekNumber {
			currentWeekRegistrations = append(currentWeekRegistrations, reg)
		}
	}
	return currentWeekRegistrations, nil
}

// Leader is the resolver for the leader field.
func (r *queryResolver) Leader(ctx context.Context, id string) (*model.Leader, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}
	// Retrieve the leader from the database based on the provided ID
	leader := &model.Leader{}
	// Assuming you have a GORM database connection named "db"

	if err := r.DB.Preload("Students.Registrations").Where("id = ?", id).First(&leader).Error; err != nil {
		return nil, err
	}

	// Return the leader
	return leader, nil
}

// Leaders is the resolver for the leaders field.
func (r *queryResolver) Leaders(ctx context.Context, sort *model.SortInput, groupBy []string) ([]*model.Leader, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	/* err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	} */
	// Retrieve the list of leaders from the database
	var leaders []*model.Leader

	// Assuming you have a GORM database connection named "db"
	query := r.DB.Preload("Students.Registrations")

	// Apply grouping based on the specified fields in the groupBy argument
	if len(groupBy) > 0 {
		// Validate and add valid group fields to the GROUP BY clause
		validGroupFields := make([]string, 0)
		for _, groupField := range groupBy {
			if schemas.IsValidGroupField(groupField) {
				validGroupFields = append(validGroupFields, groupField)
			} else {
				// Handle error when the groupField is not valid or not supported for grouping
				return nil, fmt.Errorf("invalid groupBy field: %s", groupField)
			}
		}

		if len(validGroupFields) > 0 {
			// Include "id" column in the GROUP BY clause along with other valid group fields
			groupByClause := "id, " + strings.Join(validGroupFields, ", ")

			// Add the GROUP BY clause to the query
			query = query.Group(groupByClause)
		}
	}

	sortCondition, err := schemas.GetSortingCondition(sort)
	if err != nil {
		return nil, err
	}
	if sortCondition != "" {
		query = query.Order(sortCondition)
	}

	if err := query.Find(&leaders).Error; err != nil {
		return nil, err
	}

	// Return the list of leaders with populated students
	return leaders, nil
}

// LeadersByIds is the resolver for the leadersByIds field.
func (r *queryResolver) LeadersByIds(ctx context.Context, id []*string) ([]*model.Leader, error) {
	// Extract authentication context information using middleware (if needed).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve leaders from the database based on the provided IDs.
	var leaders []*model.Leader
	// Assuming you have a GORM database connection named "db"

	if err := r.DB.Preload("Students.Registrations").Find(&leaders, "id IN (?)", id).Error; err != nil {
		return nil, err
	}

	// Return the list of leaders.
	return leaders, nil
}

// RegistrationsByLeader is the resolver for the registrationsByLeader field.
func (r *queryResolver) RegistrationsByLeader(ctx context.Context, leaderID string) ([]*model.Registration, error) {
	var registrations []*model.Registration
	if err := r.DB.Where("leader_id = ?", leaderID).Find(&registrations).Error; err != nil {
		return nil, err
	}
	return registrations, nil
}

// Student is the resolver for the student field.
func (r *queryResolver) Student(ctx context.Context, id string) (*model.Student, error) {
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}
	// Retrieve the student from the database based on the provided ID
	student := &model.Student{}
	if err := r.DB.First(student, id).Error; err != nil {
		return nil, err
	}

	// Return the student
	return student, nil
}

// Students is the resolver for the students field.
func (r *queryResolver) Students(ctx context.Context, sort *model.SortInput) ([]*model.Student, error) {
	// Extract LeaderID from the request context (provided by AuthenticationMiddleware).
	err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve the list of students from the database
	var students []*model.Student

	query := r.DB.Model(&model.Student{}).Preload("Registrations")

	sortCondition, err := schemas.GetSortingCondition(sort)
	if err != nil {
		return nil, err
	}
	if sortCondition != "" {
		query = query.Order(sortCondition)
	}

	// Assuming you have a GORM database connection named "db"
	err = query.Find(&students).Error
	if err != nil {
		return nil, err
	}

	// Return the list of students with populated registrations
	return students, nil
}

// StudentsByLeader is the resolver for the studentsByLeader field.
func (r *queryResolver) StudentsByLeader(ctx context.Context, leaderID string) ([]*model.Student, error) {
	/* err := middleware.ExtractCTXinfo(ctx)
	if err != nil {
		return nil, err
	} */
	// Retrieve the list of students by the provided leaderID from the database
	var students []*model.Student
	// Use Preload to fetch students with eager-loaded registrations
	err := r.DB.Model(&model.Student{}).Where("leader_id = ?", leaderID).Preload("Registrations").Find(&students).Error
	if err != nil {
		return nil, err
	}

	// Return the list of students
	return students, nil
}

// StudentRegistrations is the resolver for the studentRegistrations field.
func (r *queryResolver) StudentRegistrations(ctx context.Context, studentID string) ([]*model.Registration, error) {
	// Fetch all registrations associated with the studentID from the database
	var registrations []*model.Registration
	if err := r.DB.Where("student_id = ?", studentID).Find(&registrations).Error; err != nil {
		return nil, err
	}

	return registrations, nil
}

// ID is the resolver for the id field.
func (r *registrationResolver) ID(ctx context.Context, obj *model.Registration) (string, error) {
	id := obj.ID.String()
	// Log the ID

	return id, nil
}

// ID is the resolver for the id field.
func (r *studentResolver) ID(ctx context.Context, obj *model.Student) (string, error) {
	id := obj.ID.String()
	// Log the ID

	return id, nil
}

// Leader returns LeaderResolver implementation.
func (r *Resolver) Leader() LeaderResolver { return &leaderResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Registration returns RegistrationResolver implementation.
func (r *Resolver) Registration() RegistrationResolver { return &registrationResolver{r} }

// Student returns StudentResolver implementation.
func (r *Resolver) Student() StudentResolver { return &studentResolver{r} }

type leaderResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type registrationResolver struct{ *Resolver }
type studentResolver struct{ *Resolver }
