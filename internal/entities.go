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
	ID                 uint `gorm:"primary_key"`
	ApplicationId      uint
	EnvName            string `gorm:"unique;not null"`
	WriteWithoutPrefix bool   `gorm:"default:0;not null"`
	Active             bool   `gorm:"default:1;not null"`
	Variables          []Variable
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Variable struct {
	ID            uint `gorm:"primary_key"`
	EnvironmentId uint
	VarName       string `gorm:"unique;not null"`
	VarValue      string `gorm:"unique;not null"`
	Active        bool   `gorm:"default:1;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
