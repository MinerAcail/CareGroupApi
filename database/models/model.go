package models

import (
	"time"
)



type Registration struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LastComment string    `json:"lastComment"`
	Absence     bool      `json:"absence"`
	Present     bool      `json:"present"`
	LeaderID    *string   `json:"leaderID,omitempty"`
	MemberID    uint      `gorm:"index" json:"memberID"`
	Member      Member    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Church struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `json:"name"`
	SubChurches []*SubChurch `json:"subChurches"`
}


type SubChurch struct {
    ID       uint       `gorm:"primaryKey" json:"id"`
    Name     string     `json:"name"`
    ChurchID uint       `gorm:"index" json:"churchID"`
    Church   Church     `gorm:"foreignKey:ChurchID" json:"church"`
    Leaders  []*Member  `gorm:"many2many:subchurch_leaders;" json:"leaders"`
    Members  []*Member  `gorm:"many2many:subchurch_members;" json:"members"`
}

type Member struct {
    ID           uint         `gorm:"primaryKey" json:"id"`
    Name         string       `json:"name"`
    Email        string       `json:"email"`
    PhoneNumber  *string      `json:"phoneNumber,omitempty"`
    Location     *string      `json:"location,omitempty"`
    Day          string       `json:"day"`
    Types        string       `json:"types"`
    Token        string       `json:"token"`
    Registrations []*Registration `gorm:"foreignKey:MemberID" json:"registrations,omitempty"`
    ChurchID     uint         `gorm:"index" json:"churchID"`
    Church       Church       `gorm:"foreignKey:ChurchID" json:"church"`
    UpdatedAt    time.Time    `json:"updatedAt"`
    CreatedAt    time.Time    `json:"createdAt"`
}
