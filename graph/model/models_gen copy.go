package model

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

type LeaderWithRegistrations struct {
	ID              string   `json:"id"`
	RegistrationIds []string `json:"registrationIds"`
}

type LeaderAggregations struct {
	TotalLeaders             int64   `json:"totalLeaders"`
	AverageStudentsPerLeader float64 `json:"averageStudentsPerLeader"`
	MaxStudentsPerLeader     int     `json:"maxStudentsPerLeader"`
}
type LeaderStats struct {
    LeaderID         uuid.UUID
    ReferenceIDCount int
}
type Church struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string       `json:"name"`
	Password    *string      `json:"password,omitempty"`
	Email       *string      `json:"email,omitempty"`
	Types       *string      `json:"types,omitempty"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	CreatedAt   time.Time    `json:"createdAt"`
	Token       *string      `json:"token,omitempty"`
	SubChurches []*SubChurch `json:"subChurches,omitempty" `
}
type Member struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name             string          `json:"name"`
	Email            string          `json:"email"`
	PhoneNumber      *string         `json:"phoneNumber,omitempty"`
	Location         *string         `json:"location,omitempty"`
	Day              string          `json:"day"`
	Password         *string         `json:"password,omitempty"`
	Types            *string         `json:"types,omitempty"`
	Token            *string         `json:"token,omitempty"`
	LeaderID         *string         `json:"LeaderID,omitempty"`
	ReferenceIDCount *int            `json:"ReferenceIDCount,omitempty"`
	Registrations    []*Registration `json:"registrations,omitempty"`
	SubChurch        *SubChurch      `json:"subChurch,omitempty"`
	SubChurchID      *string         `json:"subChurchID,omitempty"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	CreatedAt        time.Time       `json:"createdAt"`
}

type Registration struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	LastComment string    `json:"lastComment"`
	Absence     bool      `json:"absence"`
	Present     bool      `json:"present"`
	CreatedAt   time.Time `json:"createdAt"`
	LeaderName  *string   `json:"leaderName,omitempty"`
	Leader      *Member   `json:"leader,omitempty"`
	LeaderID    *string   `json:"leaderID,omitempty"`
	Member      *Member   `json:"member,omitempty"`
	MemberID    string    `json:"memberID"`
}
type SubChurch struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Church    *Church   `json:"church,omitempty"`
	ChurchID  string    `json:"churchId"`
	Leaders   []*Member `json:"leaders,omitempty"`
	Members   []*Member `json:"members,omitempty"`
}

type Returns interface {
	IsReturns()
	// *Member
	// *Church
	// *SubChurch
}

func (Member) IsReturns() {}

func (Church) IsReturns()    {}
func (SubChurch) IsReturns() {}

// Token     *string   `json:"token,omitempty"`
//
//	Password  *string   `json:"password,omitempty"`
//	Types     *string   `json:"types,omitempty"`
//	Email     *string   `json:"email,omitempty"`
func IsNumeric(s string) bool {
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return numericRegex.MatchString(s)
}
func RemoveSpacesFromNumber(s string) string {
	var result string
	for _, char := range s {
		if char != ' ' {
			result += string(char)
		}
	}
	return result
}
func IsValidEmail(email string) bool {
	// Regular expression for a basic email format check
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
