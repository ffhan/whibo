package main

import "time"

type Result struct {
	Failed   map[string]error `json:"failed" yaml:"failed"`
	Projects []*Project       `json:"projects" yaml:"projects"`
}

func NewResult() *Result {
	return &Result{
		Failed:   make(map[string]error),
		Projects: make([]*Project, 0, 4),
	}
}

func (r *Result) AddProject(project *Project) {
	r.Projects = append(r.Projects, project)
}

type Project struct {
	Name     string    `json:"name" yaml:"name"`
	Path     string    `json:"path" yaml:"path"`
	Branches []*Branch `json:"branches" yaml:"branches"`
}

func NewProject(name string, path string) *Project {
	return &Project{Name: name, Path: path, Branches: make([]*Branch, 0, 4)}
}

func (p *Project) AddBranch(branch *Branch) {
	p.Branches = append(p.Branches, branch)
}

type Branch struct {
	Name    string    `json:"name" yaml:"name"`
	Commits []*Commit `json:"commits" yaml:"commits"`
}

func NewBranch(name string) *Branch {
	return &Branch{Name: name, Commits: make([]*Commit, 0, 4)}
}

func (b *Branch) AddCommit(commit *Commit) {
	b.Commits = append(b.Commits, commit)
}

type Commit struct {
	Name   string    `json:"name" yaml:"name"`
	Author string    `json:"author" yaml:"author"`
	Date   time.Time `json:"date" yaml:"date"`
	Hash   string    `json:"hash" yaml:"hash"`
}

func NewCommit(name string, author string, date time.Time, hash string) *Commit {
	return &Commit{Name: name, Author: author, Date: date, Hash: hash}
}
