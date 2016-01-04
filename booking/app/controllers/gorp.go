package controllers

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
	"github.com/uileyar/samples/booking/app/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	fmt.Println("InitDB in")
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

	t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true

	t = Dbm.AddTable(models.Hotel{}).SetKeys(true, "HotelId")

	t = Dbm.AddTable(models.Booking{}).SetKeys(true, "BookingId")
	t.ColMap("User").Transient = true
	t.ColMap("Hotel").Transient = true
	t.ColMap("CheckInDate").Transient = true
	t.ColMap("CheckOutDate").Transient = true

	Dbm.TraceOn("[gorp]", r.INFO)

	err := Dbm.CreateTables()
	if err == nil {
		insertDemoData()
	}
}

func insertDemoData() {
	fmt.Println("insertDemoData in")

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

	hotels := []*models.Hotel{
		&models.Hotel{0, "111 Marriott Courtyard", "Tower Pl, Buckhead", "Atlanta", "GA", "30305", "USA", 120},
		&models.Hotel{0, "222 W Hotel", "Union Square, Manhattan", "New York", "NY", "10011", "USA", 450},
		&models.Hotel{0, "333 Hotel Rouge", "1315 16th St NW", "Washington", "DC", "20036", "USA", 250},
	}
	for n, hotel := range hotels {
		fmt.Printf("%d: hotel = %v\n", n+1, hotel)
		if err := Dbm.Insert(hotel); err != nil {
			panic(err)
		}
	}
	fmt.Println("insertDemoData out")
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	fmt.Println("Begin in")

	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	fmt.Println("Begin out")
	return nil
}

func (c *GorpController) Commit() r.Result {
	fmt.Println("Commit in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	fmt.Println("Commit out")
	return nil
}

func (c *GorpController) Rollback() r.Result {
	fmt.Println("Rollback in")
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	fmt.Println("Rollback out")
	return nil
}
