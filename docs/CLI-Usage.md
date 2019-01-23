---
title: CLI Usage
layout: wiki
permalink: /wiki/CLI-Usage
---

# Command-line Tool (mailing service)

The command-line tool can be used to interact with the server's E-Mail
Recipience and Dispatch APIs. A user can create an e-mail which is sent to the
server and has to supply a valid e-mail address, a subject and the message. The
created mail is then delivered to the E-Mail-Receive API of the server where a
new ticket will be created out of this mail. Optionally, the user may supply an
existing ticket id to create a new answer instead of a new ticket.

Another option is to retrieve all created e-mails on server-side. The mails are
fetched from the server and output to the console. For each fetched mail the 
command-line tool confirms the successful sending of the e-mail by calling
another API to verify the sending. The tool cannot send genuine e-mails of
course, but the sending of the mails is simulated. After each mail is verified
to be sent correctly the server deletes these mails from cache and from file
system.

## Build and Execution

The command-line tool can be built in the same way than the server by executing

```bash
go build ./cmd/command_line_tool
```

from the project's root directory. A better option is to use the attached start
scripts for UNIX and Windows systems. To use these, run

```bash
./startCLI.sh [options]
```

or `startCLI.bat` on Windows, respectively. Options can be appended to the call
and they are delegated to the command-line tool. See the following section for a
list of available options.

## Usage

The command-line tool can be run in the following way:

```bash
./command_line_tool [options]
```

If the command-line tool is started without options, the user is guided by an
interactive menu. Otherwise he has to supply the required options by himself and
there is no prompting for user input.

In the following the available options are listed.

## Client options

The following options can be used to configure the client's connection settings.

### `-host <HOST>`

Use another host name or IP address for the server the client is connecting to.

**Default**: `localhost`

### `-port <PORT>`

Specify the port the server listens to. The `PORT` argument has to be a 16 bit
unsigned integer (0 < `PORT` &le; 65535) and has to match the server port.

**Default**: `8443`

### `-cert <FILE>`

Use another SSL certificate for the connection to the server. The `FILE`
argument has to be an existing file with a valid SSL certificate. Note that the
path has to be relative to the current working directory.

**Default**: `./ssl/server.cert`

## Fetch options

The fetch options configure the way e-mails are retrieved from the server.

### `-f` (fetch)

Fetch the server-side created e-mails from the server. If this option is set,
there will be no interactive prompt, but only the messages are retrieved.

**Default**: `false`

## Submit options

The submit options are used to specify e-mail properties for a direct submission
of an e-mail without user inputs.

### `-s` (submit)

Use this flag to directly submit a message to the server without prompting. The
mail properties are set by the `-email`, `-subject` and `-message` flags and are
required. If the `-tID` option is specified with a valid ticket id a new answer
will be created instead of a new ticket.

**Default**: `false`

### `-email <EMAIL>`

Specify the sender's e-mail address used for the creation of a new e-mail. The
`EMAIL` argument has to be a valid e-mail address. Only applicable in
conjunction to the `-s` flag.

**Default**: empty

### `-subject <SUBJECT>`

Specify the subject of the new created ticket or answer. The `SUBJECT` argument
should not be empty. Only applicable in conjunction to the `-s` flag.

**Default**: empty

### `-message <MSG>`

Specify the message for the new created ticket or answer. The `MSG` argument
should not be empty. Only applicable in conjunction to the `-s` flag.

**Default**: empty

### `-tID <ID>`

Specify the ticket id of an existing ticket to append the message as a new
answer to this ticket. This flag is optional and if it is not set or if the
ticket id is invalid, a new ticket will be created. Only applicable in
conjunction to the `-s` flag.

**Default**: empty

## Help options

The help options provide information about the command-line tool and its
provided flags.

### `-h`, `-help`

Print a help text with information about the flags and exit.
