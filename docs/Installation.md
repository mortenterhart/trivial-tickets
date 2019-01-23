---
title: Installation
layout: wiki
permalink: wiki/Installation

---

# Installation

To install the Trivial Tickets ticket system, first ensure you have
a working Go environment installed. Go version 1.7 and above are
supported. Visit the [Golang Installation
Site](https://golang.org/doc/install) for installation instructions
and system requirements. Set the `GOPATH` environment variable to
your desired location and then execute

```bash
go get -u -v github.com/mortenterhart/trivial-tickets
```

This command requires [git](https://git-scm.com) to be installed.
It downloads and installs all project and test dependencies. For
more information on how to set the `GOPATH` environment variable
accordingly visit [Setting
GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) from the
Go Wiki.

The ticket system is now installed under
`$GOPATH/src/github.com/mortenterhart/trivial-tickets`. Head over
to :book: [Build and Execution](Build-and-Execution.md) for the next
steps.

:warning: Please do not use the `go install` command as it puts the
binary to `$GOPATH/bin` and the default file and directory paths
get incorrect. If you do this please be aware that you always have
to provide the relative paths to the resource folders via command-line
flags (see :book: [Server Usage](Server-Usage.md) for more information).
