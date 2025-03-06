
package main

import (
	"Aggregator/internal/config"
	"Aggregator/internal/database"
	"database/sql"
	"context"
	"log"
	"os"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}


func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("")
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db: dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", reset)
	cmds.register("users", listUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", midddlewareLoggedIn(handlerAddfeed))
	cmds.register("feeds", listfeeds)
	cmds.register("follow", midddlewareLoggedIn(follow))
	cmds.register("following", midddlewareLoggedIn(listfeedsfollow))
	cmds.register("unfollow", midddlewareLoggedIn(unfollow))
	cmds.register("browse", midddlewareLoggedIn(browse))

	if len(os.Args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
		return 
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}

func midddlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}