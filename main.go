package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	apiKey       = os.Getenv("OPENAI_API_KEY")
	client       *openai.Client
	systemCard   string
	preTaskCard  string
	postTaskCard string
)

// Define a root command
var rootCmd = &cobra.Command{
	Use:   "juno",
	Short: "Juno is a simple CLI app",
}

// Define a clone command
// Define a clone command
var cloneCmd = &cobra.Command{
	Use:   "clone [url]",
	Short: "Clone a git repository",
	Args:  cobra.ExactArgs(1), // Require exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		// Get the repository URL from the argument
		url := args[0]

		// Get the current working directory
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}

		// Create a folder in the current working directory with the name of the repository
		repoName := path.Base(url)                      // Get the last element of the URL
		repoName = strings.TrimSuffix(repoName, ".git") // Remove the .git suffix if any
		dir = path.Join(dir, repoName)                  // Join the directory and the repository name
		err = os.Mkdir(dir, 0755)                       // Create the folder with read-write-execute permissions for owner and read-execute permissions for group and others
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}

		fmt.Println("Cloning", url, "to", dir)
		repo, err := git.PlainClone(dir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Println("Error cloning repository:", err)
			return
		}

		fmt.Println("Cloned successfully:", repo)
	},
}

// Define a do command
// Define a do command
var doCmd = &cobra.Command{
	Use:   "do [command]",
	Short: "Print the command",
	Args:  cobra.ExactArgs(1), // Require exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		// Get the command from the argument
		command := args[0]

		// Print the command
		fmt.Println("Command:", command)

		executeTask(client, systemCard, preTaskCard+"\n"+command+"\n"+postTaskCard)
	},
}

func main() {
	systemCard = readFile("cards/00-system-card.txt")
	preTaskCard = readFile("cards/01-pre-task.txt")
	postTaskCard = readFile("cards/02-post-task.txt")

	client = openai.NewClient(apiKey)

	// Add the sub commands to the root command
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.Run = doCmd.Run

	// Execute the root command
	rootCmd.Execute()
}
