# mailbox-duplicate-search [![Build Status](https://travis-ci.org/alexhokl/mailbox-duplicate-search.svg?branch=master)](https://travis-ci.org/alexhokl/mailbox-duplicate-search)

CLI tool to find out duplicate mails in a folder

### Usage

To dump filenames of duplicated mail messages and statistics of the search in the same directory.

```sh
export MAILBOX_SEARCH_IS_DRY_RUN=true
mailbox-duplicate-search
```

To delete all the duplicated mail messages in the same directory.

```sh
export MAILBOX_SEARCH_IS_DRY_RUN=false
rm $(mailbox-duplicate-search)
```

### Installation

##### Option 1

If you have Go installed, all you need is `go get github.com/alexhokl/mailbox-duplicate-search`.

##### Option 2

Download binary from release page and put the binary in one of the directories
specified in `PATH` enviornment variable.

