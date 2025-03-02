
package main

import (
	"Aggregator/internal/config"
	"Aggregator/internal/database"
	"database/sql"
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
	cmds.register("addfeed", handlerfeed)

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