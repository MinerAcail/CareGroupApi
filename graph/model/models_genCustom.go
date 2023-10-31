package model

import (
	"database/sql/driver"
	"regexp"
	"strings"
	"time"

	// "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/kobbi/vbciapi/mypkg"
	"github.com/lib/pq"
	// "github.com/kobbi/vbciapi/mypkg"
)

type Post struct {
	ID   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Tags pq.StringArray `gorm:"type:text[]"`
}
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
	Token       *string      `json:"token,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	SubChurches []*SubChurch `json:"subChurches,omitempty" `
}
type MyType struct {
	ID      uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	MyArray mypkg.Myarray `gorm:"type:text[]" json:"MyArray"`
}

type MyArrayType struct {
	Myarray mypkg.Myarray
}

func (t *MyArrayType) Scan(value interface{}) error {
	// Implement the Scan method to convert the database value into Myarray
	if bytes, ok := value.([]byte); ok {
		// Convert bytes to your custom Myarray type
		t.Myarray = mypkg.Myarray(strings.Split(string(bytes), ","))
	}
	return nil
}

func (t MyArrayType) Value() (driver.Value, error) {
	// Implement the Value method to convert Myarray into a database-compatible value
	return strings.Join(t.Myarray, ","), nil
}

type MyArr struct {
	ID      uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()"`
	MyArray mypkg.Myarray `gorm:"type:text[]" json:"MyArray"`
}

type Member struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phoneNumber,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Day         string    `json:"day"`
	Password    *string   `json:"password,omitempty"`
	// Types            *string       `json:"types"`
	Types pq.StringArray `gorm:"type:text[]" `

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
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	LastComment string     `json:"lastComment"`
	Absence     bool       `json:"absence"`
	Present     bool       `json:"present"`
	CreatedAt   time.Time  `json:"createdAt"`
	LeaderName  *string    `json:"leaderName,omitempty"`
	Leader      *Member    `json:"leader,omitempty"`
	LeaderID    *string    `json:"leaderID,omitempty"`
	Member      *Member    `json:"member,omitempty"`
	SubChurch   *SubChurch `json:"subChurch,omitempty"`
	SubChurchID *string    `json:"subChurchID,omitempty"`
	MemberID    string     `json:"memberID"`
}
type SubChurch struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Password  *string   `json:"password,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Types     *string   `json:"types,omitempty"`
	Token     *string   `json:"token,omitempty"`
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
