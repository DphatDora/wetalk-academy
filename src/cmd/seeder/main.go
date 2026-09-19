package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"wetalk-academy/config"
	"wetalk-academy/seeders/data"
	"wetalk-academy/seeders/migration"
	"wetalk-academy/seeders/runner"
	"wetalk-academy/seeders/validation"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "seeder:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usageError()
	}
	command := args[0]
	if command == "migrate" {
		if len(args) != 2 || (args[1] != "up" && args[1] != "down") {
			return usageError()
		}
		return withDatabase(func(ctx context.Context, db *mongo.Database) error {
			if args[1] == "up" {
				return migration.Up(ctx, db)
			}
			return migration.Down(ctx, db)
		})
	}

	if command != "validate" && command != "seed" && command != "clear" {
		return usageError()
	}
	selected, keys, err := parseSelection(command, args[1:])
	if err != nil {
		return err
	}
	if command == "validate" {
		if err := validation.Validate(selected); err != nil {
			return err
		}
		printSummary("validated", selected)
		return nil
	}

	return withDatabase(func(ctx context.Context, db *mongo.Database) error {
		seedRunner := runner.New(db)
		if command == "seed" {
			if err := seedRunner.Seed(ctx, selected); err != nil {
				return err
			}
			printSummary("seeded", selected)
			return nil
		}
		if err := seedRunner.Clear(ctx, keys); err != nil {
			return err
		}
		fmt.Printf("cleared %d topic dataset(s): %s\n", len(keys), strings.Join(keys, ", "))
		return nil
	})
}

func parseSelection(command string, args []string) ([]data.TopicDataset, []string, error) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	all := flags.Bool("all", false, "select all topic datasets")
	only := flags.String("only", "", "select one topic dataset by key")
	if err := flags.Parse(args); err != nil {
		return nil, nil, err
	}
	if *all == (*only != "") {
		return nil, nil, fmt.Errorf("use exactly one of --all or --only <topic-key>")
	}
	selected, err := data.Select(*only)
	if err != nil {
		return nil, nil, err
	}
	keys := data.Keys()
	if *only != "" {
		keys = []string{*only}
	}
	return selected, keys, nil
}

func withDatabase(action func(context.Context, *mongo.Database) error) error {
	conf := config.GetConfig()
	if strings.TrimSpace(conf.Database.URI) == "" || strings.TrimSpace(conf.Database.Name) == "" {
		return fmt.Errorf("MONGO_URI and MONGO_DB_NAME are required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(conf.Database.URI))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return action(ctx, client.Database(conf.Database.Name))
}

func printSummary(action string, datasets []data.TopicDataset) {
	var lessons, contents, quizzes int
	for _, dataset := range datasets {
		lessons += len(dataset.Lessons)
		contents += len(dataset.Contents)
		quizzes += len(dataset.Quizzes)
	}
	fmt.Printf("%s topics=%d lessons=%d contents=%d quizzes=%d\n", action, len(datasets), lessons, contents, quizzes)
}

func usageError() error {
	return fmt.Errorf("usage: seeder migrate <up|down> | seeder <validate|seed|clear> <--all|--only topic-key>")
}
