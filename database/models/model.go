package models

import "time"

type Leader struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Email       *string    `json:"email,omitempty"`
	PhoneNumber string     `gorm:"not null" json:"phoneNumber"`
	Day         string     `gorm:"not null" json:"day"`
	Password    string     `gorm:"not null" json:"password"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	Students    []*Student `gorm:"foreignKey:LeaderID" json:"students"`
}

type Registration struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	LastComment string    `gorm:"not null" json:"lastComment"`
	Absence      bool      `gorm:"not null" json:"absence"`
	Present      bool      `gorm:"not null" json:"present"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	StudentID    int       `gorm:"index" json:"studentId"`
}

type Student struct {
	ID            int             `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"not null" json:"name"`
	Email         *string         `json:"email,omitempty"`
	PhoneNumber   string          `gorm:"not null" json:"phoneNumber"`
	Day           string          `gorm:"not null" json:"day"`
	LeaderID      int             `gorm:"index" json:"leaderID"`
	Leader        Leader          `gorm:"foreignKey:LeaderID" json:"leader"`
	Registrations []*Registration `gorm:"foreignKey:StudentID" json:"registrations"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
}
