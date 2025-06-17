package database

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Ubah tipe kolom ID di tabel users
	if err := db.Exec("ALTER TABLE users MODIFY COLUMN id CHAR(36)").Error; err != nil {
		return err
	}

	// Tambahkan default UUID untuk kolom ID
	if err := db.Exec("ALTER TABLE users MODIFY COLUMN id CHAR(36) DEFAULT (UUID())").Error; err != nil {
		return err
	}

	// Ubah tipe kolom ID dan CreatedBy di tabel categories
	if err := db.Exec("ALTER TABLE categories MODIFY COLUMN id CHAR(36)").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE categories MODIFY COLUMN created_by CHAR(36)").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE categories MODIFY COLUMN id CHAR(36) DEFAULT (UUID())").Error; err != nil {
		return err
	}

	return nil
}
