package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wyvern-api/models"
)

// UserRepo struct
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo initiate TransactionRepo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByID is method to find user by id
func (repo *UserRepo) FindByID(ID int64) (models.User, error) {
	var user models.User
	result := repo.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, ID)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

// Update is method to update user
func (repo *UserRepo) Update(db *gorm.DB, user models.User) (models.User, error) {
	result := db.Save(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
