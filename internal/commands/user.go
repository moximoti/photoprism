package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/pkg/txt"
	"github.com/urfave/cli"
	"os"
	"strings"
)

// PasswdCommand updates a password.
var UserCommand = cli.Command{
	Name:  "user",
	Usage: "Manage Users from CLI",
	Subcommands: []cli.Command{
		{
			Name:   "create",
			Usage:  "creates a new user. Provide at least username and password",
			Action: userCreate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "fullname, n",
					Usage: "full name of the new user",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "unique username",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "sets the users password",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "deletes user by username",
			Action:    userDelete,
			ArgsUsage: "takes username as argument",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "execute deletion",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "prints a list of all users",
			Action: userList,
		},
	},
}

func userCreate(ctx *cli.Context) error {
	return withDependencies(ctx, func(conf *config.Config) error {

		var newUser = entity.User{
			RoleAdmin:    true, // TODO change back to false when implementing access control
			RoleGuest:    false,
			UserDisabled: false,
		}

		if ctx.String("username") != "" && ctx.String("password") != "" {
			log.Debugf("creating user in non-interactive mode")
			newUser.FullName = strings.TrimSpace(ctx.String("fullname"))
			newUser.UserName = strings.TrimSpace(ctx.String("username"))
			newUser.Password = strings.TrimSpace(ctx.String("password"))
		}

		if ctx.String("fullname") == "" && ctx.String("username") == "" {
			fmt.Printf("please enter full name: ")
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			newUser.FullName = strings.TrimSpace(text)
		}

		if ctx.String("username") == "" {
			fmt.Printf("please enter a username: ")
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			newUser.UserName = strings.TrimSpace(text)
		}

		if ctx.String("password") == "" {
			for {
				fmt.Printf("please enter a new password for %s (at least 4 characters)\n", txt.Quote(newUser.UserName))
				pw := getPassword("New password: ")
				if confirm := getPassword("Confirm password: "); confirm == pw {
					newUser.Password = pw
					break
				} else {
					log.Infof("passwords did not match or too short. please try again\n")
				}
			}
		}

		if err := newUser.CreateAndValidate(conf); err != nil {
			return err
		}
		return nil
	})
}

func userDelete(ctx *cli.Context) error {
	return withDependencies(ctx, func(conf *config.Config) error {
		username := ctx.Args()[0]
		if !ctx.Bool("force") {
			user := entity.FindUserByName(username)
			if user != nil {
				log.Infof("found user %s with uid: %s. Use -f to perform actual deletion", user.UserName, user.UserUID)
				return nil
			}
			return errors.New("user not found")
		}
		err := entity.DeleteUserByName(username)
		if err != nil {
			log.Errorf("%s", err)
			return nil
		}
		log.Infof("sucessfully deleted %s", username)
		return nil
	})
}

func userList(ctx *cli.Context) error {
	return withDependencies(ctx, func(conf *config.Config) error {
		users := entity.AllUsers()
		for _, user := range users {
			//fmt.Printf("%s ", user.UserUID)
			if user.UserName != "" {
				fmt.Printf("%s ", user.UserName)
			} else {
				fmt.Printf("[%s]", user.FullName)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("total users found: %v\n", len(users))
		return nil
	})
}

func withDependencies(ctx *cli.Context, f func(conf *config.Config) error) error {
	conf := config.NewConfig(ctx)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := conf.Init(); err != nil {
		return err
	}

	conf.InitDb()

	// command is executed here
	if err := f(conf); err != nil {
		return err
	}

	conf.Shutdown()
	return nil
}
