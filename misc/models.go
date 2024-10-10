package misc

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string
	Password   string
	Session_id string
}

type Workout struct {
	gorm.Model
	Name    string
	Seances []Seance `gorm:"many2many:workout_seances;"`
	OwnerID int
	Owner   User
}

type Seance struct {
	gorm.Model
	Name      string
	Exercises []Exercise `gorm:"many2many:seance_exercises;"`
	OwnerID   int
	Owner     User
}

type Exercise struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex"`
	OwnerID int
	Owner   User
	Seances []Seance `gorm:"many2many:seance_exercises;"`
}

type RealWorkout struct {
	gorm.Model
	Template   Workout
	Time_start time.Time
	Time_end   time.Time
}
