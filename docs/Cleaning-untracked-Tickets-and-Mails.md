---
title: Cleaning untracked Tickets and Mails
layout: wiki
permalink: /wiki/Cleaning-untracked-Tickets-and-Mails
---

# Cleaning untracked Tickets and Mails

Occasionally, the server will create new tickets and mails according to
requirements. If a ticket is created or updated it is written to its own file.
Likewise, e-mails are generated on creation or update of a ticket which are also
saved in their own files. These files are usually located under `files/tickets`
and `files/mails`, unless otherwise specified by command-line flags.

The repository contains some predefined tickets as examples that are loaded by
the web server on startup. These tickets are tracked by Git and all other ticket
and mail files are ignored. This way the repository can be kept clean.

If required, the tickets and mails not tracked by Git can be easily removed by
the script `cleanupUntrackedTickets.sh`. It should be called on the command-line
and prompts for confirmation of the removal. If the confirmation is successful
all untracked tickets and mails will be deleted. The script offers some
command-line options which change the way files are going to be deleted.

## Usage

`cleanupUntrackedTickets.sh [option(s)]`

Delete untracked tickets and mails from server cache in `files/tickets` and
`files/mails` to revert the server to its original state. Before deletion, the
user is prompted to confirm the removal. Note that only files with `.json`
extension (matching `*.json`) are considered to be removed.

The confirmation prompt offers some commands which decide whether to continue
the deletion or not:

* `y`: Confirm and continue the deletion (if `--interactive` was supplied the
  menu will be shown before deletion)
* `n`: Show which files would be removed and return to the prompt (same as
  `--dry-run`)
* `q`: Quit
* `*`: Any other input aborts the deletion

If no untracked files matching the patterns are found, the script exits.

## Options

The following options configure the process of removal and are optional.
Mandatory arguments to long options are also mandatory for short options.
Long options may be abbreviated as long as they do not get ambiguous.

### `-n`, `--dry-run`

Show which files would be deleted, but do not delete actually any.

### `-e`, `--exclude=<pattern>`

Exclude files matching the `pattern` from deletion.

### `-f`, `--force`

Force deletion without prompting for.

### `-h`, `--help`

Display a help text and exit.

### `-i`, `--interactive`

Delete files using an interactive menu. The menu lists all files ready to be
deleted and there are options on excluding files based on a filter pattern and
other options too. The menu is displayed after the initial confirmation succeded.

### `-m`, `--mail-dir=<DIR>`

Specify another directory to search for untracked mail files. The directory in
the `DIR` argument has to exist and needs to point to a directory inside the
Trivial Tickets Git repository. Otherwise it cannot be determined if files are
already tracked by Git. The directory should be the same than the one that the
server uses for its mails, otherwise the script may delete wrong files which
causes data loss. First, make sure you are removing the desired files by using
either the `--dry-run` option or the `n` command in the prompt.

### `-q`, `--quiet`

Do not print the file paths as they are removed.

### `-t`, `--ticket-dir=<DIR>`

Specify another directory to search for untracked ticket files. The directory in
the `DIR` argument has to exist and needs to point to a directory inside the
Trivial Tickets Git repository. Otherwise it cannot be determined if files are
already tracked by Git. The directory should be the same than the one that the
server uses for its tickets, otherwise the script may delete wrong files which
causes data loss. First, make sure you are removing the desired files by using
either the `--dry-run` option or the `n` command in the prompt.
