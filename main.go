package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FabulousBernt/gator/internal/config"
	"github.com/FabulousBernt/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login requires a username argument")
	}

	// Check if the user exists in the database
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Only set the user if they exist
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("User set to:", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	// 1. Check if a name was provided
	if len(cmd.args) == 0 {
		return errors.New("register requires a username argument")
	}

	// 2. Create the user in the database
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	// 3. Set the current user in the config
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w", err)
	}

	// 4. Print the created user
	fmt.Println("User created successfully!")
	fmt.Printf("ID: %v, Name: %v, CreatedAt: %v\n", user.ID, user.Name, user.CreatedAt)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't reset database: %w", err)
	}

	fmt.Println("Users successfully deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get all users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("addfeed requires a name and a url")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Printf("Feed created successfully!\n")
	fmt.Printf("ID: %v\nName: %v\nURL: %v\nUserID: %v\n", feed.ID, feed.Name, feed.Url, feed.UserID)
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func main() {
	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Open database connection
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create database queries instance
	dbQueries := database.New(db)

	// Store both in state
	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: not enough arguments, a command is required")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
