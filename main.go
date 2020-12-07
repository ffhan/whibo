package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	logRgx = regexp.MustCompile("Author: (.*)\\nDate:[ ]+(.*)[ \\n]+(.*)")

	sinceFlag   = flag.String("since", "7", "how many days before today")
	authorsFlag = flag.String("authors", "", "authors separated by a comma")
)

func main() {
	flag.Parse()

	failed := make(map[string]error)

	since := setupSince()

	path := setupPath()
	allAuthors := setupAuthors()

	dir, err := ioutil.ReadDir(path)
	must(err)

	for _, d := range dir {
		if d.IsDir() && time.Now().Sub(d.ModTime()) <= since {
			fmt.Printf("-----------------------%s-----------------------\n", d.Name())

			target := filepath.Join(path, d.Name())
			output, err := gitLog(since, target)
			if err != nil {
				failed[d.Name()] = err
				continue
			}
			outputAuthorCommits(output, allAuthors)
		}
	}

	printFailed(failed)
}

func setupSince() time.Duration {
	days, err := strconv.Atoi(*sinceFlag)
	must(err)
	return time.Duration(days) * 24 * time.Hour
}

func setupPath() string {
	path := flag.Arg(0)

	if path == "" {
		wd, err := os.Getwd()
		must(err)
		path = wd
		log.Printf("path not set - using CD %s\n", path)
	}
	return path
}

func setupAuthors() []string {
	allAuthors := strings.Split(*authorsFlag, ",")
	if len(allAuthors) == 0 {
		output, err := exec.Command("git", "config", "user.name").Output()
		must(err)
		allAuthors = append(allAuthors, string(output))
		log.Println("authorsFlag not set, setting the default git username")
	}
	return allAuthors
}

func gitLog(since time.Duration, target string) ([]byte, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("--since=\"%d days ago\"", int(math.Round(since.Hours()/24))))
	cmd.Dir = target
	output, err := cmd.Output()
	return output, err
}

func outputAuthorCommits(output []byte, allAuthors []string) {
	commits := logRgx.FindAllStringSubmatch(string(output), -1)

	for _, match := range commits {
		author := match[1]
		date := match[2]
		commitName := match[3]

		if isAuthor(allAuthors, author) {
			fmt.Println("\t- ", date, commitName)
		}
	}
}

func printFailed(failed map[string]error) {
	fmt.Println("\nFailed: ")
	for name, err := range failed {
		fmt.Printf("\t* %s: %v\n", name, err)
	}
}

func isAuthor(authors []string, currentAuthor string) bool {
	for _, author := range authors {
		if strings.Contains(strings.ToUpper(currentAuthor), strings.ToUpper(author)) {
			return true
		}
	}
	return false
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
