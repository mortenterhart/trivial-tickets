---
title: Server Usage
layout: wiki
permalink: wiki/Server-Usage

---

# Server Usage

The ticket system main executable
[`ticketsystem.go`](https://github.com/mortenterhart/trivial-tickets/blob/master/cmd/ticketsystem/ticketsystem.go)
offers the following configuration flags. They can be used to alter
the default configuration of the server and the logger. No option
is required since every option has a meaningful default value.

The usage of the executable looks like following:

```bash
./ticketsystem [options]
```

## Server options

The following server options can change important connection settings
such as the port and file or directory paths to important resources.

### `-port <PORT>`

Specify the port on which the web server will be launched. On startup
a message with the server's URL will be shown. The `PORT` is a 16
bit unsigned integer (0 < `PORT` &le; 65535) and should not be used
by another process. If the specified port is blocked the server
throws an error at startup.

**Default**: `8443`

### `-tickets <DIR>`

Configure the directory where the tickets will be stored. The `DIR`
argument can be an existing directory with active write privileges,
otherwise it is created on startup automatically. Note that the
path has to be relative to the current working directory. The
`*.json` files in this directory have to contain valid JSON content.

**Default**: `./files/tickets`

### `-users <FILE>`

Change the default path to the users file. The `FILE` argument has
to be an existing JSON file with user definitions. Note that the
path has to be relative to the current working directory.

**Default**: `./files/users/users.json`

### `-mails <DIR>`

Specify the directory where new mails created by the server will
be saved. The `DIR` argument can be an existing directory with
active write privileges, otherwise it is created on startup
automatically.  Note that the path has to be relative to the current
working directory. If `DIR` already contains `*.json` files they
have to contain valid JSON format.

**Default**: `./files/mails`

### `-cert <FILE>`

Set the path to the SSL server certificate. The `FILE` argument has
to be an existing file with a valid SSL certificate. Note that the
path has to be relative to the current working directory.

**Default**: `./ssl/server.cert`

### `-key <FILE>`

Set the path to the SSL server key. The `FILE` argument has to be
an existing file with a valid SSL public key. Note that the path
has to be relative to the current working directory.

**Default**: `./ssl/server.key`

### `-web <DIR>`

Change the root directory of the web server. The `DIR` argument has
to be an existing directory and should match the templates path and
static paths pointing to the server resources (such as
`/static/js/ticketsystem.js`).

**Default**: `./www`

## Logging options

The logging options alter the way messages are logged to the console.

### `-log-level <LEVEL>`

Specify the level of logging. The given `LEVEL` can be one of the
following:

* `info`: display all log messages (recommended)
* `warning`: display only warnings, errors and fatal errors
* `error`: display only errors and fatal errors
* `fatal`: display only fatal errors (not recommended)

Fatal errors are always logged and cause the server to shutdown
directly. Therefore it is not advised to set the log level to
`fatal`.

**Default**: `info`

### `-verbose`

Enable verbose logging output (includes information about package
path and function name, filenames and line numbers on every log
message). For better readability, the package and file paths are
trimmed to the last component, so leading directories will be
stripped.

**Default**: `false`

### `-full-paths`

Disable the abbreviation of package and file paths. Paths are written
as is to the log. This option is compatible to `-verbose`.

Warning: This will extend log messages a lot making it harder to
read. Depending on your screen size the messages will probably not
fit on the screen.

**Default**: `false`

## Help options

The help options provide information about the usage of the ticket
system and its flags.

### `-h`, `-help`

Print a help text with information about the flags and exit.
