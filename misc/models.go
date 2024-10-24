package misc

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string
	Password            string
	Session_id          string
	ActiveRealWorkoutID int
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
	OwnerID     int
	Owner       User
	Active      bool
	Finished    bool
	CurrentWeek int
	TemplateID  int
	Template    Workout
	Time_start  time.Time
	Time_end    time.Time
}

type RealSeance struct {
	gorm.Model
	Week                   int
	OwnerID                int
	Owner                  User
	Active                 bool
	Finished               bool
	TemplateID             int
	Template               Seance
	CorrespondingWorkoutID int
	CorrespondingWorkout   RealWorkout
}

type RealExercise struct {
	gorm.Model
	OwnerID               int
	Owner                 User
	TemplateID            int
	Template              Exercise
	Finished              bool
	CorrespondingSeanceID int
	CorrespondingSeance   RealSeance
}

type RealSet struct {
	gorm.Model
	OwnerID                 int
	Owner                   User
	CorrespondingExerciseID int
	CorrespondingExercise   RealExercise
	Reps                    int
	Weight                  int
}
