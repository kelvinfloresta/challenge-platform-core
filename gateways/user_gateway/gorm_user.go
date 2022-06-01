package user_gateway

import (
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
)

type GormUserGatewayFacade struct {
	DB *database.Database
}

func (g GormUserGatewayFacade) GetById(id string) (*GetByIdUserGatewayOutput, error) {
	data := &models.User{}
	result := g.DB.Conn.Limit(1).Find(data, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	email := ""
	if data.Email != nil {
		email = *data.Email
	}

	adapted := &GetByIdUserGatewayOutput{
		ID:          data.ID,
		Name:        data.Name,
		Password:    data.Password,
		Email:       email,
		Phone:       data.Phone,
		BirthDate:   data.BirthDate,
		Document:    data.Document,
		JobPosition: data.JobPosition,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		DeletedAt:   data.DeletedAt.Time,
	}

	return adapted, nil
}

func (g GormUserGatewayFacade) GetByEmail(email string) (*GetByEmailUserGatewayOutput, error) {
	data := &models.User{}
	result := g.DB.Conn.Limit(1).Find(data, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	if data.Email == nil {
		return nil, nil
	}

	adapted := &GetByEmailUserGatewayOutput{
		ID:          data.ID,
		Name:        data.Name,
		Password:    data.Password,
		Email:       *data.Email,
		Phone:       data.Phone,
		BirthDate:   data.BirthDate,
		Document:    data.Document,
		JobPosition: data.JobPosition,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		DeletedAt:   data.DeletedAt.Time,
	}

	return adapted, nil
}

func (g GormUserGatewayFacade) ChangePassword(input ChangePasswordInput) (bool, error) {
	updated :=
		g.DB.Conn.Model(&models.User{}).Where("id = ?", input.UserID).Update("password", input.Password)

	return updated.RowsAffected > 0, updated.Error
}
