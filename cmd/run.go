package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dikmit/rakup"
	"github.com/spf13/cobra"
)

type runOptions struct {
	Verbose  bool
	Kuchisu  int
	MaxWord  int
	WordFile string
	User     string
	Password string
}

// NewRunCmd creates a new `rakup run` command
func NewRunCmd() *cobra.Command {
	options := runOptions{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the search",
		PreRun: func(cmd *cobra.Command, args []string) {
			if options.Verbose {
				log.SetOutput(os.Stderr)
			} else {
				log.SetOutput(ioutil.Discard)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRun(cmd, options)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.User, "user", "u", "", "login user")
	flags.StringVarP(&options.Password, "password", "p", "", "login password")
	flags.BoolVarP(&options.Verbose, "verbose", "v", false, "Print a verbose message")
	flags.IntVarP(&options.Kuchisu, "kuchisu", "n", 30, "Search max kuchisu")
	flags.IntVar(&options.MaxWord, "max-word", 1000, "Save max word line")
	flags.StringVar(&options.WordFile, "word-file", "words.txt", "Search word file")

	return cmd
}

func runRun(cmd *cobra.Command, options runOptions) error {
	log.Println("run command")
	user := conf.User
	if options.User != "" {
		user = options.User
	}
	password := conf.Password
	if options.Password != "" {
		password = options.Password
	}
	if user == "" || password == "" {
		return fmt.Errorf("please set config")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := rakup.NewBrowser()
	if err != nil {
		return err
	}
	defer b.Close()

	if err := b.Login(user, password); err != nil {
		return err
	}

	words, err := rakup.ReadWords(options.WordFile)
	if err != nil {
		log.Printf("words read error: %s\n", err)
	}

	log.Println("トレンドワード取得開始")
	trendWords, err := b.GetTrend(ctx)
	if err != nil {
		log.Printf("get trend error: %s\n", err)
	}
	if len(trendWords) > 0 {
		log.Printf("トレンドワード %v\n", trendWords)
		// トレンドワードを保存
		if len(words) < options.MaxWord {
			wordExists := func(s string, words []string) bool {
				for _, word := range words {
					if s == word {
						return true
					}
				}
				return false
			}
			for _, tWord := range trendWords {
				if !wordExists(tWord, words) {
					if err := rakup.SaveWord(tWord, options.WordFile); err != nil {
						cmd.PrintErrf("word write error: %s\n", err)
					}
				}
			}
		}
		words = append(words, trendWords...)
	}

	rand.Seed(time.Now().UnixNano())
	var searchCount int
	// MAX口数になるまで検索
	for searchCount < options.Kuchisu {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ctx canceld: %w", ctx.Err())
		default:
			word := words[rand.Intn(len(words))]
			log.Printf("search word: %s\n", word)
			i, err := b.Search(word)
			if err != nil {
				log.Println(err)
				continue
			}
			searchCount = i
		}
		log.Printf("search loop kuchi: %d", searchCount)
	}
	cmd.Printf("Searched %d kuchi\n", searchCount)
	return nil
}
