package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type OutputType string

const (
	outputShell OutputType = "shell"
	outputJson  OutputType = "json"
	outputYaml  OutputType = "yaml"
)

type OutputWriter interface {
	SetWriter(writer io.Writer)
	Write(result *Result) error
}

func WriteOutput(result *Result, outWriter io.Writer, outType OutputType) error {
	writer, ok := writers[outType]
	if !ok {
		return errors.New("no writer for the provided output type")
	}
	writer.SetWriter(outWriter)
	return writer.Write(result)
}

var (
	writers = map[OutputType]OutputWriter{
		outputYaml:  &yamlWriter{},
		outputJson:  &jsonWriter{},
		outputShell: &shellWriter{},
	}
)

type baseWriter struct {
	writer io.Writer
}

func (b *baseWriter) SetWriter(writer io.Writer) {
	b.writer = writer
}

type jsonWriter struct {
	baseWriter
}

func (j *jsonWriter) Write(result *Result) error {
	encoder := json.NewEncoder(j.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

type yamlWriter struct {
	baseWriter
}

func (y *yamlWriter) Write(result *Result) error {
	return yaml.NewEncoder(y.writer).Encode(result)
}

func getTerminalWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	split := strings.Split(string(out), " ")

	if len(split) != 2 {
		return 0, errors.New("unexpected stty size output")
	}
	return strconv.Atoi(split[0])
}

type shellWriter struct {
	baseWriter
}

func (s *shellWriter) printf(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(s.writer, format, args...)
	return err
}

func (s *shellWriter) print(arg interface{}) error {
	_, err := fmt.Fprint(s.writer, arg)
	return err
}

func (s *shellWriter) printProjectHeader(project *Project, width int) error {
	dashes := width - len(project.Name)
	if dashes < 0 {
		return s.print(project.Name)
	}
	if err := s.printDashes(dashes); err != nil {
		return err
	}
	if err := s.print(project.Name); err != nil {
		return err
	}
	if err := s.printDashes(dashes); err != nil {
		return err
	}
	if dashes%2 == 0 {
		return s.print("\n")
	}
	return s.print("-\n")
}

func (s *shellWriter) printDashes(dashes int) error {
	for i := 0; i < dashes/2; i++ {
		if err := s.printf("-"); err != nil {
			return err
		}
	}
	return nil
}

const oldMonitorStandardWidth = 80

func (s *shellWriter) Write(result *Result) error {
	width, err := getTerminalWidth()
	if err != nil {
		width = oldMonitorStandardWidth
		_, _ = fmt.Fprintf(os.Stderr, "couldn't get terminal width, assumed width of 80: %v\n", err)
	}
	for _, project := range result.Projects {
		if err = s.printProjectHeader(&project, width); err != nil {
			return err
		}

		if *groupByBranch {
			if err := s.printGroupedByBranch(&project); err != nil {
				return err
			}
		} else {
			if err := s.printSortedByDate(&project); err != nil {
				return err
			}
		}
	}
	return s.printFailed(result.Failed)
}

func (s *shellWriter) printGroupedByBranch(project *Project) error {
	for _, branch := range project.Branches {
		if _, err := fmt.Fprintln(s.writer, "\t+ ", colorGreen, branch.Name, colorReset); err != nil {
			return err
		}
		for _, commit := range branch.Commits {
			if _, err := fmt.Fprintln(s.writer, "\t\t- ",
				colorCyan, commit.Author, colorReset, commit.Date,
				colorYellow, commit.Name, colorReset); err != nil {
				return err
			}
		}
	}
	return nil
}

type ExpandedCommit struct {
	Commit
	BranchName string
}

func (s *shellWriter) printSortedByDate(project *Project) error {
	commits := s.getCommits(project)
	for _, commit := range commits {
		if _, err := fmt.Fprintln(s.writer, "\t- ",
			colorCyan, commit.Author, colorGreen, commit.BranchName, colorReset, commit.Date,
			colorYellow, commit.Name, colorReset); err != nil {
			return err
		}
	}
	return nil
}

func (s *shellWriter) getCommits(project *Project) []ExpandedCommit {
	numberOfCommits := 0
	for _, branch := range project.Branches {
		numberOfCommits += len(branch.Commits)
	}
	commits := make([]ExpandedCommit, 0, numberOfCommits)
	for _, branch := range project.Branches {
		for _, commit := range branch.Commits {
			commits = append(commits, ExpandedCommit{
				Commit:     commit,
				BranchName: branch.Name,
			})
		}
	}
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})
	return commits
}

func (s *shellWriter) printFailed(failed map[string]error) error {
	if _, err := fmt.Fprintln(s.writer, colorRed, "\nFailed: "); err != nil {
		return err
	}
	for name, err := range failed {
		if err := s.printf("\t* %s: %v\n%s", name, err, colorReset); err != nil {
			return err
		}
	}
	return nil
}
