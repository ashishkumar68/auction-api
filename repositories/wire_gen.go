// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package repositories

import (
	"gorm.io/gorm"
)

// Injectors from di_container.go:

func NewUserRepository(conn *gorm.DB) *UserRepository {
	userRepository := initUserRepository(conn)
	return userRepository
}

func NewItemRepository(conn *gorm.DB) *ItemRepository {
	itemRepository := initItemRepository(conn)
	return itemRepository
}

func NewBidRepository(conn *gorm.DB) *BidRepository {
	bidRepository := initBidRepository(conn)
	return bidRepository
}
