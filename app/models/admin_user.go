package models

import "time"

// AdminUser is table
type AdminUser struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func migrateAdminUser() {
	DbConnection.AutoMigrate(&AdminUser{})
}
