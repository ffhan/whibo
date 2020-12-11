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
	logRgx    = regexp.MustCompile("commit (\\w+).*\\nAuthor: (.*)\\nDate:[ ]+(.*)[ \\n]+(.*)")
	branchRgx = regexp.MustCompile("[ ]+(.*)")

	sinceFlag     = flag.String("since", "7", "how many days before today")
	authorsFlag   = flag.String("authors", "", "authors separated by a comma")
	outputFlag    = flag.String("o", "shell", "output format (supported shell, json, yaml)")
	groupByBranch = flag.Bool("group-by-branch", false, groupByBranchDesc)
)

const (
	groupByBranchDesc = "group by branches - if false commits will be sorted by date without grouping by branch"
	colorGreen        = "\033[32m"
	colorReset        = "\033[0m"
	colorYellow       = "\033[33m"
	colorRed          = "\033[31m"
	colorCyan         = "\033[36m"
)

func main() {
	flag.Parse()

	since := setupSince()

	path := setupPath()
	allAuthors := setupAuthors()

	dir, err := ioutil.ReadDir(path)
	must(err)

	result := NewResult()

dirLoop:
	for _, d := range dir {
		if d.IsDir() && time.Now().Sub(d.ModTime()) <= since {
			dirName := d.Name()

			target := filepath.Join(path, dirName)
			project := result.AddProject(dirName, target)

			branches, err := gitBranches(target)
			if !handleError(dirName, err, result.Failed) {
				continue
			}

			for _, branchName := range branches {
				branch := project.AddBranch(branchName)
				output, err := gitLog(since, target, branchName)
				if !handleError(dirName, err, result.Failed) {
					continue dirLoop
				}
				parseCommits(result, project, branch, output, allAuthors)
			}
		}
	}

	must(WriteOutput(result, os.Stdout, OutputType(*outputFlag)))
}

func gitBranches(target string) ([]string, error) {
	cmd := exec.Command("git", "branch", "-l")
	cmd.Dir = target

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, strings.ReplaceAll(string(output), "\n", "; "))
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
	if *authorsFlag == "" {
		output, err := exec.Command("git", "config", "user.name").Output()
		must(err)
		username := string(output)
		log.Println("authorsFlag not set, setting the default git username", username)
		return []string{username}
	}
	allAuthors := strings.Split(*authorsFlag, ",")
	return allAuthors
}

func gitLog(since time.Duration, target, branch string) ([]byte, error) {
	cmd := exec.Command("git", "log", branch, fmt.Sprintf("--since=\"%d days ago\"", int(math.Round(since.Hours()/24))))
	cmd.Dir = target
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, strings.ReplaceAll(string(output), "\n", "; "))
	}
	return output, nil
}

func parseCommits(result *Result, project *Project, branch *Branch, output []byte, allAuthors []string) {
	commits := logRgx.FindAllStringSubmatch(string(output), -1)

	for _, match := range commits {
		hash := match[1]
		author := match[2]
		date := match[3]
		commitName := match[4]

		parsedDate, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", date)
		if !handleError(project.Name, err, result.Failed) {
			return
		}

		if isAuthor(allAuthors, author) {
			branch.AddCommit(commitName, author, hash, parsedDate)
		}
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
