package commands

import (
	"fmt"
	"log"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/spf13/cobra"

	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/models"
)

var addUserCmd = &cobra.Command{
	Use:   "add_user [id] [email] [firstname] [lastname] [password] [admin]",
	Short: "Add a new user",
	Long:  `Insert a new user in database.`,
	Run:   addUser,
}

func addUser(cmd *cobra.Command, args []string) {
	if len(args) < 5 {
		cmd.Usage()
		log.Fatalln("Missing arguments")
	}

	dbSession := models.NewDBSession()

	if user := dbSession.FindUser(args[0]); user != nil {
		log.Fatalln("There is already a user with that id: " + args[0])
	}

	if user := dbSession.FindUserByEmail(args[1]); user != nil {
		log.Fatalln("There is already a user with that email: " + args[1])
	}

	password, err := bcrypt.GenerateFromPassword([]byte(args[4]), bcrypt.DefaultCost)
	if err != nil {
		panic("Arg")
	}

	isAdmin := (len(args) >= 5) && (args[4] == "true")

	user := &models.User{
		ID:        args[0],
		Email:     args[1],
		FirstName: args[2],
		LastName:  args[3],
		Admin:     isAdmin,
		Status:    models.UserStatusActive,
		Lang:      core.DefaultLang,
		Password:  string(password),
	}

	if err := dbSession.CreateUser(user); err != nil {
		log.Fatalln(fmt.Sprintf("Failed to create user: %v", err))
	}
}
