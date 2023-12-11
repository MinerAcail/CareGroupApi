package model

import (
	"database/sql/driver"
	// "encoding/json"
	// "encoding/json"
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
	Types            pq.StringArray              `gorm:"type:text[]" `
	Pwood            *string                     `json:"pwood,omitempty"`
	ChurchMinistries []*MemberChurchMinistryRole `json:"churchMinistries,omitempty"`

	Token            *string         `json:"token,omitempty"`
	Leader           *Member         `json:"leader,omitempty"`
	LeaderID         *string         `json:"LeaderID,omitempty"`
	ReferenceIDCount *int            `json:"ReferenceIDCount,omitempty"`
	Registrations    []*Registration `json:"registrations,omitempty"`
	SubChurch        *SubChurch      `json:"subChurch,omitempty"`
	SubChurchID      *string         `json:"subChurchID,omitempty"`
	PersonalInfor    *FamilyInfo     `json:"personalInfor,omitempty"`
	PersonalInforID  *string         `json:"personalInforId,omitempty"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	CreatedAt        time.Time       `json:"createdAt"`
}

type FamilyInfo struct {
	ID                       uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()"`
	LastName                 *string           `json:"lastName,omitempty"`
	SpouseID                 *string           `json:"spouseId,omitempty"`
	DateOfBirth              *string           `json:"dateOfBirth,omitempty"`
	Member                   *Member           `gorm:"foreignKey:MemberID;references:ID" json:"member,omitempty"`
	MemberID                 string            `json:"memberID"`
	Spouse                   *Member           `json:"spouse,omitempty"`
	SpouseNameNotVbci        *string           `json:"spouseNameNotVbci,omitempty"`
	SpousePhoneNumberNotVbci *string           `json:"spousePhoneNumberNotVbci,omitempty"`
	ChildrenID               *string           `json:"childrenId,omitempty"`
	// Children     []*Member `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children                []Member `gorm:"many2many:family_info_children;"`
	Relationship             *string           `json:"relationship,omitempty"`
	NextOfKin                *string           `json:"nextOfKin,omitempty"`
	Occupation               *JobInfo          `json:"occupation,omitempty"`
	OccupationID             *string           `json:"occupationId,omitempty"`
	Education                *string           `json:"education,omitempty"`
	EmergencyContact         *EmergencyContact `json:"emergencyContact,omitempty"`
	EmergencyContactID       *string           `json:"emergencyContactId,omitempty"`
	UpdatedAt                time.Time         `json:"updatedAt"`
	CreatedAt                time.Time         `json:"createdAt"`
}
type MemberChildren struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	ParentID   string    `json:"ParentID"`
	ChildrenID string    `json:"childrenId"`
}



// func (fi FamilyInfo) Value() (driver.Value, error) {
// 	return json.Marshal(fi.Children)
// }

//	func (fi *FamilyInfo) Scan(value interface{}) error {
//		return json.Unmarshal(value.([]byte), &fi.Children)
//	}
type ChurchMinistryRole struct {
	ID   uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4()"`
	Role *ChurchMinistryRolesEnum `json:"role,omitempty"`
}
type MemberChurchMinistryRole struct {
	ID                   uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4()"`
	MemberID             string              `json:"memberID"`
	ChurchMinistryRoleID string              `json:"churchMinistryRoleID"`
	ChurchMinistryRole   *ChurchMinistryRole `json:"ChurchMinistryRole,omitempty"`
}
type EmergencyContact struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        *string   `json:"name,omitempty"`
	PhoneNumber *string   `json:"phoneNumber,omitempty"`

	Relation  *string   `json:"relation,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}
type JobInfo struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	TypeOfWork     *string   `json:"typeOfWork,omitempty"`
	Position       *string   `json:"position,omitempty"`
	Company        *string   `json:"company,omitempty"`
	WorkExperience *string   `json:"workExperience,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}
type Finance struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string       `json:"name"`
	Password    *string      `json:"password,omitempty"`
	Email       *string      `json:"email,omitempty"`
	Types       *string      `json:"types,omitempty"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	CreatedAt   time.Time    `json:"createdAt"`
	Token       *string      `json:"token,omitempty"`
	SubChurchID *string      `json:"subChurchID,omitempty"`
	SubChurches []*SubChurch `json:"subChurches"`
}

type Registration struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	LastComment  string     `json:"lastComment"`
	Absence      bool       `json:"absence"`
	Present      bool       `json:"present"`
	CreatedAt    time.Time  `json:"createdAt"`
	LeaderName   *string    `json:"leaderName,omitempty"`
	MemberName   *string    `json:"memberName,omitempty"`
	Leader       *Member    `json:"leader,omitempty"`
	LeaderID     *string    `json:"leaderID,omitempty"`
	Member       *Member    `json:"member,omitempty"`
	SubChurch    *SubChurch `json:"subChurch,omitempty"`
	SubChurchID  *string    `json:"subChurchID,omitempty"`
	MemberID     string     `json:"memberID"`
	Report       *bool      `json:"report,omitempty"`
	Done         *bool      `json:"done,omitempty"`
	TempLeaderID *string    `json:"tempLeaderID,omitempty"`
}
type RegistrationByCallAgent struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CallAgentID string    `json:"callAgentId"`
	CallAgent   *Member   `json:"callAgent,omitempty"`
	// SubChurch        *SubChurch      `json:"subChurch,omitempty"`
	// SubChurchID      *string         `json:"subChurchID,omitempty"`
	Day pq.StringArray `gorm:"type:text[]" `
	// Registrations []*Registration `json:"registrations,omitempty"`
}
type SubChurch struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name         string      `json:"name"`
	Password     *string     `json:"password,omitempty"`
	Email        *string     `json:"email,omitempty"`
	Types        *string     `json:"types,omitempty"`
	Token        *string     `json:"token,omitempty"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	CreatedAt    time.Time   `json:"createdAt"`
	Church       *Church     `json:"church,omitempty"`
	ChurchID     string      `json:"churchId"`
	Leaders      []*Member   `json:"leaders,omitempty"`
	Members      []*Member   `json:"members,omitempty"`
	IsLocal      *bool       `json:"isLocal,omitempty"`
	CallCenterID *string     `json:"callCenterId,omitempty"`
	CallCenter   *CallCenter `json:"callCenter,omitempty"`
}
type CallCenter struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Password  *string   `json:"password,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Types     *string   `json:"types,omitempty"`
	Token     *string   `json:"token,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	// SubChurches []*SubChurch `json:"subChurches,omitempty" `
	SubChurches []*SubChurch `gorm:"foreignKey:CallCenterID" json:"subChurches,omitempty"`
}
type MigrationRequest struct {
	ID                  uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	LocationFrom        *string    `json:"locationFrom,omitempty"`
	LocationEnd         *string    `json:"locationEnd,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	MemberID            *string    `json:"memberID,omitempty"`
	MemberName          *string    `json:"MemberName,omitempty"`
	DestinationChurchID string     `json:"destinationChurchID,omitempty"`

	StatusID *string          `json:"statusID,omitempty"`
	Status   *MigrationStatus `json:"status,omitempty"`
}
type Returns interface {
	IsReturns()
	// *Member
	// *Church
	// *SubChurch
}

func (Member) IsReturns() {}

func (Church) IsReturns()     {}
func (SubChurch) IsReturns()  {}
func (CallCenter) IsReturns() {}

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
