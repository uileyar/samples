package jobs

import (
	"fmt"

	"github.com/uileyar/samples/booking/app/controllers"
	"github.com/uileyar/samples/booking/app/models"
)

// Periodically count the bookings in the database.
type BookingCounter struct{}

func (c BookingCounter) Run() {
	bookings, err := controllers.Dbm.Select(models.Booking{},
		`select * from Booking`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d bookings.\n", len(bookings))
}

type UserCounter struct{}

func (u UserCounter) Run() {
	users, err := controllers.Dbm.Select(models.User{},
		`select * from user`)
	if err != nil {
		panic(err)
	}

	for n, v := range users {
		fmt.Printf("There are %d user %v\n", n+1, v)
	}
}

func init() {
	/*
		revel.OnAppStart(func() {
			jobs.Schedule("@every 1m", BookingCounter{})
		})
		revel.OnAppStart(func() {
			jobs.Schedule("@every 1m", UserCounter{})
		})
	*/
}
