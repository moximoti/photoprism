package commands

import (
	"bufio"
	"context"
	"errors"
	//"fmt"
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
			Name:   "delete",
			Usage:  "deletes user by username or email",
			Action: userDelete,
		},
		{
			Name:   "list, ls",
			Usage:  "prints a list of all users",
			Action: userList,
		},
	},
}

func userCreate(ctx *cli.Context) error {
	return withDependencies(ctx, func(conf *config.Config) error {

		var newUser = entity.User{
			RoleAdmin:    false,
			RoleGuest:    false,
			UserDisabled: false,
		}

		if ctx.String("fullname") == "" && ctx.String("username") == "" {
			log.Infof("please enter full name: ")
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			newUser.FullName = strings.TrimSpace(text)
		} else {
			newUser.FullName = strings.TrimSpace(ctx.String("fullname"))
		}

		if ctx.String("username") == "" {
			log.Infof("please enter a username: ")
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			newUser.UserName = strings.TrimSpace(text)
		} else {
			newUser.UserName = strings.TrimSpace(ctx.String("username"))
		}

		newUser.Password = ""
		if ctx.String("password") == "" {
			for {
				log.Infof("please enter a new password for %s (at least 4 characters)\n", txt.Quote(newUser.UserName))
				pw := getPassword("New password: ")
				if confirm := getPassword("Confirm password: "); confirm == pw {
					newUser.Password = pw
					break
				} else {
					log.Infof("passwords did not match or too short. please try again\n")
				}
			}
		} else {
			newUser.Password = strings.TrimSpace(ctx.String("password"))
		}

		if err := newUser.CreateAndValidate(conf.Settings().Users.PasswordPolicy == config.PolicyNone); err != nil {
			return err
		}
		return nil
	})
}

func userDelete(ctx *cli.Context) error {
	return errors.New("not implemented")
}

func userList(ctx *cli.Context) error {
	return errors.New("not implemented")
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
