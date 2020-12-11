package main

import "time"

type Result struct {
	Failed   map[string]error `json:"failed" yaml:"failed"`
	Projects []Project        `json:"projects" yaml:"projects"`
}

func NewResult() *Result {
	return &Result{
		Failed:   make(map[string]error),
		Projects: make([]Project, 0, 4),
	}
}

func (r *Result) AddProject(name, path string) *Project {
	r.Projects = append(r.Projects, Project{
		Name:     name,
		Path:     path,
		Branches: make([]Branch, 0, 4),
	})
	return &r.Projects[len(r.Projects)-1]
}

type Project struct {
	Name     string   `json:"name" yaml:"name"`
	Path     string   `json:"path" yaml:"path"`
	Branches []Branch `json:"branches" yaml:"branches"`
}

func (p *Project) AddBranch(name string) *Branch {
	p.Branches = append(p.Branches, Branch{
		Name:    name,
		Commits: make([]Commit, 0, 4),
	})
	return &p.Branches[len(p.Branches)-1]
}

type Branch struct {
	Name    string   `json:"name" yaml:"name"`
	Commits []Commit `json:"commits" yaml:"commits"`
}

func (b *Branch) AddCommit(name, author, hash string, date time.Time) *Commit {
	b.Commits = append(b.Commits, Commit{
		Name:   name,
		Author: author,
		Date:   date,
		Hash:   hash,
	})
	return &b.Commits[len(b.Commits)-1]
}

type Commit struct {
	Name   string    `json:"name" yaml:"name"`
	Author string    `json:"author" yaml:"author"`
	Date   time.Time `json:"date" yaml:"date"`
	Hash   string    `json:"hash" yaml:"hash"`
}
