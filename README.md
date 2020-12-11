# what-have-i-been-working-on (whibo)

List all your recent work in git repos. Perfect for bad time-trackers and will-do-it-laters.

## Usage

* `go run ./ -since 14 -authors 'author1,author2' ~/work` when in project root
    * lists all your work in the last 14 days, targetting author1 or author2 (case insensitive) git commit authors
* `whibo -since 2 -authors 'author1,author2' ~/work`
* `whibo -since 7 -group-by-branches ~/work` - get author name from git config, group commits by branches
* `whibo -since 14 -authors thisAuthor -o yaml` - get yaml output of commits in the last 14 days, matching case
  insensitive author name
  
Supported outputs:
* shell
* yaml
* json

## Installation

* position to the project root directory
* `go build -o whibo`
* `sudo cp whibo /usr/local/bin/`
* `whibo -h`

## Example output

### shell
![img.png](assets/output_example.png)

```
---------------------------------------------------------------taskio---------------------------------------------------------------
        -   fhancic <fhancic@croz.net>  master  2020-12-03 17:42:49 +0100 CET  removed unused files 
        -   fhancic <fhancic@croz.net>  master  2020-12-03 17:38:01 +0100 CET  updated README 
        -   fhancic <fhancic@croz.net>  master  2020-12-03 17:32:14 +0100 CET  README 
        -   fhancic <fhancic@croz.net>  master  2020-12-03 17:17:30 +0100 CET  Initial commit 
----------------------------------------------------what-have-i-been-working-on-----------------------------------------------------
        -   fhancic <fhancic@croz.net>  master  2020-12-11 19:09:57 +0100 CET  fixed terminal width, improved readme 
        -   fhancic <fhancic@croz.net>  master  2020-12-11 19:00:57 +0100 CET  filter out projects and branches that have no commits 
        -   fhancic <fhancic@croz.net>  master  2020-12-11 18:50:03 +0100 CET  multiple writers, terminal width detection, group by branch, sort commits by date 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 19:14:47 +0100 CET  lowercase name 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 19:13:30 +0100 CET  log errors 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 19:08:22 +0100 CET  readme 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 19:00:39 +0100 CET  Fixed handling authors, author autocompletion 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 18:50:38 +0100 CET  Scan branches, color output 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 18:17:30 +0100 CET  README 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 18:11:30 +0100 CET  Flag fixes, improved handling authors, since flag accepts a number of days 
        -   fhancic <fhancic@croz.net>  master  2020-12-07 18:02:15 +0100 CET  Initial commit 
 
Failed: 
        * schach: exit status 128: fatal: not a git repository (or any parent up to mount point /); Stopping at filesystem boundary (GIT_DISCOVERY_ACROSS_FILESYSTEM not set).; 
```

### Yaml 
```yaml
failed:
  schach: {}
projects:
- name: taskio
  path: /home/user/example/taskio
  branches:
  - name: master
    commits:
    - name: removed unused files
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-03T17:42:49+01:00
      hash: 2d5ee47990361ebd594d7ebf29723b6c0e91ff09
    - name: updated README
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-03T17:38:01+01:00
      hash: 8bce5cbb9612899ecbdcdcfa2272593c43346901
    - name: README
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-03T17:32:14+01:00
      hash: 797698e960100136f21f94858f6db92e14e53551
    - name: Initial commit
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-03T17:17:30+01:00
      hash: 90b6959b8967111bd4f03c48664ec1f8da1657b8
- name: what-have-i-been-working-on
  path: /home/user/example/what-have-i-been-working-on
  branches:
  - name: master
    commits:
    - name: fixed terminal width, improved readme
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-11T19:09:57+01:00
      hash: a2ee081fb4cba33aaa3c8d263ba707942737c495
    - name: filter out projects and branches that have no commits
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-11T19:00:57+01:00
      hash: 2ddddbbeeb5d0b5ad8922117b2520384272c7b66
    - name: multiple writers, terminal width detection, group by branch, sort commits
        by date
      author: fhancic <fhancic@croz.net>
      date: !!timestamp 2020-12-11T18:50:03+01:00
      hash: bca4fe324f61a0588440091dda3225c64759de92
...
```

### JSON
```json
{
  "failed": {
    "schach": {}
  },
  "projects": [
    {
      "name": "taskio",
      "path": "/home/user/example/taskio",
      "branches": [
        {
          "name": "master",
          "commits": [
            {
              "name": "removed unused files",
              "author": "fhancic \u003cfhancic@croz.net\u003e",
              "date": "2020-12-03T17:42:49+01:00",
              "hash": "2d5ee47990361ebd594d7ebf29723b6c0e91ff09"
            },
            {
              "name": "updated README",
              "author": "fhancic \u003cfhancic@croz.net\u003e",
              "date": "2020-12-03T17:38:01+01:00",
              "hash": "8bce5cbb9612899ecbdcdcfa2272593c43346901"
            },
            {
              "name": "README",
              "author": "fhancic \u003cfhancic@croz.net\u003e",
              "date": "2020-12-03T17:32:14+01:00",
              "hash": "797698e960100136f21f94858f6db92e14e53551"
            },
            {
              "name": "Initial commit",
              "author": "fhancic \u003cfhancic@croz.net\u003e",
              "date": "2020-12-03T17:17:30+01:00",
              "hash": "90b6959b8967111bd4f03c48664ec1f8da1657b8"
            }
          ]
        }
      ]
    },
    {
      "name": "what-have-i-been-working-on",
      "path": "/home/user/example/what-have-i-been-working-on",
      "branches": [
        {
          "name": "master",
...
```
