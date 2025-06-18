package constants

import "golang.org/x/crypto/bcrypt"

type (
	Role          string
	BookingStatus string
)

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"

	PasswordMinLength = 6
	BcryptCost        = bcrypt.DefaultCost

	BookingStatusPending   BookingStatus = "pending"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusCancelled BookingStatus = "cancelled"
)
