package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Repository struct {
	connection *gorm.DB
}

func initRepository(conn *gorm.DB) *Repository {
	return &Repository{
		connection: conn,
	}
}

type TrxFunc func(trx *gorm.DB) error

func (repo *Repository) Transaction(trxFunc TrxFunc) error {

	err := repo.connection.Transaction(trxFunc)

	return err
}

func (repo *Repository) Save(val any) error {
	result := repo.connection.Create(val)
	if result.Error != nil {
		log.Printf("could not save %T value to database due to error: %s", val, result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) Update(val any) error {
	result := repo.connection.Omit(clause.Associations).Updates(val)
	if result.Error != nil {
		log.Printf("could not update %T value to database due to error: %s", val, result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) UpdatesWithMap(val any, updatesMap map[string]interface{}) error {
	result := repo.connection.Model(val).Omit(clause.Associations).Updates(updatesMap)
	if result.Error != nil {
		log.Printf("could not update %T value to database due to error: %s", val, result.Error)
		return result.Error
	}

	return nil
}

func (repo *Repository) Delete(val any) error {
	result := repo.connection.Delete(val)
	if result.Error != nil {
		log.Printf("could not delete %T value from database due to error: %s", val, result.Error)
		return result.Error
	}

	return nil
}
