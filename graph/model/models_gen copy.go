package model

import (
	"time"

	"github.com/google/uuid"
)



type Leader struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name         string     `json:"name"`
	Email        *string    `json:"email,omitempty"`
	PhoneNumber  string     `json:"phoneNumber"`
	Day          string     `json:"day"`
	Password     string     `json:"password"`
	Types        string     `json:"types"`
	Token        string     `json:"token"`
	RefreshToken string     `json:"refreshToken"`
	Location     string     `json:"location"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Students     []*Student `json:"students"`
}

type LeaderWithRegistrations struct {
	ID              string   `json:"id"`
	RegistrationIds []string `json:"registrationIds"`
}
type Registration struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `json:"createdAt"`
	LastComment string    `json:"lastComment"`
	Absence     bool      `json:"absence"`
	LeaderID    *string   `json:"leaderID,omitempty"`
	Leader      *Leader   `json:"leader,omitempty"`
	Present     bool      `json:"present"`
	StudentID   string    `json:"studentID"`
	Student     *Student  `json:"student,omitempty"`
}

type Student struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"student_id"`
	Name          string          `json:"name"`
	Email         *string         `json:"email,omitempty"`
	PhoneNumber   string          `json:"phoneNumber"`
	Day           string          `json:"day"`
	LeaderID      string          `json:"leaderID"`
	Registrations []*Registration `json:"registrations,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type LeaderAggregations struct {
	TotalLeaders             int64   `json:"totalLeaders"`
	AverageStudentsPerLeader float64 `json:"averageStudentsPerLeader"`
	MaxStudentsPerLeader     int     `json:"maxStudentsPerLeader"`
}
type GroupByResults struct {
	Total          *int       `json:"total,omitempty"`
	GroupByResults []*GroupBy `json:"groupByResults,omitempty"`
}
