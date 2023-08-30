package database

import (
	"github.com/kobbi/vbciapi/database/models"
	"github.com/kobbi/vbciapi/graph/model"
	"gorm.io/gorm"
)

func CreateLeader(db *gorm.DB, input model.CreateLeaderInput) (*models.Leader, error) {
	leader := &models.Leader{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Day:         input.Day,
		Password:    input.Password,
	}

	if err := db.Create(leader).Error; err != nil {
		return nil, err
	}

	return leader, nil
}
