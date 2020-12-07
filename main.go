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
	logRgx    = regexp.MustCompile("Author: (.*)\\nDate:[ ]+(.*)[ \\n]+(.*)")
	branchRgx = regexp.MustCompile("[ ]+(.*)")

	sinceFlag   = flag.String("since", "7", "how many days before today")
	authorsFlag = flag.String("authors", "", "authors separated by a comma")
)

const (
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func main() {
	flag.Parse()

	failed := make(map[string]error)

	since := setupSince()

	path := setupPath()
	allAuthors := setupAuthors()

	dir, err := ioutil.ReadDir(path)
	must(err)

dirLoop:
	for _, d := range dir {
		if d.IsDir() && time.Now().Sub(d.ModTime()) <= since {
			dirName := d.Name()
			fmt.Printf("-----------------------%s-----------------------\n", dirName)

			target := filepath.Join(path, dirName)
			branches, err := gitBranches(target)
			if !handleError(dirName, err, failed) {
				continue
			}

			for _, branch := range branches {
				output, err := gitLog(since, target, branch)
				if !handleError(dirName, err, failed) {
					continue dirLoop
				}
				outputAuthorCommits(branch, output, allAuthors)
			}
		}
	}

	printFailed(failed)
}

func gitBranches(target string) ([]string, error) {
	cmd := exec.Command("git", "branch", "-l")
	cmd.Dir = target

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := branchRgx.FindAllStringSubmatch(string(output), -1)

	result := make([]string, len(branches))
	for i, branch := range branches {
		result[i] = branch[1]
	}
	return result, nil
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

func gitLog(since time.Duration, target, branch string) ([]byte, error) {
	cmd := exec.Command("git", "log", branch, fmt.Sprintf("--since=\"%d days ago\"", int(math.Round(since.Hours()/24))))
	cmd.Dir = target
	output, err := cmd.Output()
	return output, err
}

func outputAuthorCommits(branch string, output []byte, allAuthors []string) {
	commits := logRgx.FindAllStringSubmatch(string(output), -1)

	for _, match := range commits {
		author := match[1]
		date := match[2]
		commitName := match[3]

		if isAuthor(allAuthors, author) {
			fmt.Println("\t- ", colorGreen, branch, colorReset, date, colorYellow, commitName, colorReset)
		}
	}
}

func printFailed(failed map[string]error) {
	fmt.Println(colorRed, "\nFailed: ")
	for name, err := range failed {
		fmt.Printf("\t* %s: %v\n%s", name, err, colorReset)
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

// if true - no error, false - error occurred
func handleError(target string, err error, failed map[string]error) bool {
	if err != nil {
		failed[target] = err
		return false
	}
	return true
}
