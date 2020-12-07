# what-have-i-been-working-on (WHIBO)

List all your recent work in git repos. Perfect for bad time-trackers and will-do-it-laters.

## Usage

* `go run main.go -since 14 -authors 'author1,author2' ~/work`
    * lists all your work in the last 14 days, targetting author1 or author2 git commit authors
* `whibo -since 14 -authors 'author1,author2' ~/work`

## Installation

* `go build -o whibo main.go`
* `sudo cp whibo /usr/local/bin/`
* `whibo -h`
