package internal

import "time"

type Application struct {
	ID           uint   `gorm:"primary_key"`
	Name         string `gorm:"unique;not null"`
	Active       bool   `gorm:"default:1;not null"`
	Environments []Environment
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Environment struct {
	ID                 uint   `gorm:"primary_key"`
	ApplicationId      uint   `gorm:"not null;index:env_unique"`
	EnvName            string `gorm:"not null;index:env_unique"`
	WriteWithoutPrefix bool   `gorm:"default:0;not null"`
	Active             bool   `gorm:"default:1;not null"`
	Variables          []Variable
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Variable struct {
	ID            uint   `gorm:"primary_key"`
	EnvironmentId uint   `gorm:"not null;index:var_unique"`
	VarName       string `gorm:"not null;index:var_unique"`
	VarValue      string `gorm:"not null"`
	Active        bool   `gorm:"default:1;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
