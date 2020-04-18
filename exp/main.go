package main

import (
	"fmt"

	"lenslocked.com/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
}
type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"dbname=%s sslmode=disable",
	// 	host, port, user, dbname)
	// db, err := gorm.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// db.LogMode(true)

	// user := models.User{
	// 	Name:  "Michael Scott",
	// 	Email: "michael@dundermifflin.com",
	// }
	// user.Name = "Updated Name"
	// us, err := models.NewUserService(psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer us.Close()
	// us.DestructiveReset()
	// if err := us.Create(&user); err != nil {
	// 	panic(err)
	// }

	// user.Name = "Updated Name"
	// if err := us.Update(&user); err != nil {
	// 	panic(err)
	// }
	// foundUser, err := us.ByEmail("michael@dundermifflin.com")
	// if err := us.Delete(foundUser.ID); err != nil {
	// 	panic(err)
	// }
	// // Verify the user is deleted
	// _, err = us.ByID(foundUser.ID)
	// if err != models.ErrNotFound {
	// 	panic("user was not deleted!")
	// }
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
