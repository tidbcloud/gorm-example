package main

import (
	"github.com/google/uuid"
	"testing"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver, tidb

func TestGORM(t *testing.T) {
	user := User{Name: uuid.NewString()}

	DB.Create(&user)

	var result User
	if err := DB.First(&result, user.ID).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	}
}

func TestUniqueKey(t *testing.T) {

	name := uuid.NewString()
	user := User{Name: name}
	if err := DB.Create(&user).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	}
	if err := DB.Create(&user).Error; err == nil {
		t.Error("Should return error because of the same name")
	}

	var result User
	if err := DB.First(&result, user.ID).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	}
}
