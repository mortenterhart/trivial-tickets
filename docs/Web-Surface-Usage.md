---
title: Web Surface Usage
layout: wiki
permalink: /wiki/Web-Surface-Usage
---

# Web Surface Usage

The ticket system can be controlled via the Command Line Interface
(CLI) (also known as the command-line tool) or via the web surface.

Tickets can be created, viewed and edited via the web interface.
The unregistered users can only add comments to the tickets. After
logging in, administrators / employees can also change the status
of a ticket, merge two tickets and change the agent of the ticket.
The CLI communicates with the server via a Web API. From the ticket
system server's point of view, it simulates a mail application that
queries messages from the server via the API and sends them to the
respective e-mail addresses. In addition, new messages can also be
entered to the ticket system server via the CLI, which are then
used on the server side either to create a new ticket or to attach
a comment to an existing ticket.

The following operations can be performed using the web interface:

* [Log into the System](#log-into-the-system)
* [Logout of the System](#logout-of-the-system)
* [Create a new Ticket](#create-a-new-ticket)
* [Assign a User to a Ticket](#assign-a-user-to-a-ticket)
* [View Tickets](#view-tickets)
* [Change the Status of a Ticket](#change-the-status-of-a-ticket)
* [Merge Tickets](#merge-tickets)
* [Activate the Vacation Mode](#activate-the-vacation-mode)
* [Remark on the Field _Reply Type_](#remark-on-the-field-reply-type)

## Log into the System

The login mask in the header is used to log in. There you have to
enter a valid username and password and click on <kbd>Login</kbd>.
Some predefined users are registered for the login and must be used.
If the login is successful, you will be redirected to your dashboard
where you can view your settings and your assigned tickets. In
addition, the login form disappears, for which the username and a
logout button are displayed instead. However, if the login failed,
you remain on the start page and have no access to the tickets.

## Logout of the System

Once you are logged in, you can also log out again. To do this,
click on the <kbd>Logout</kbd> button in the header next to your
username. This will terminate the session and you are logged out.
Then you will be redirected to the start page.

## Create a new Ticket

If you are not logged in, fill in the fields on the start page and
click on <kbd>Create ticket</kbd>. In the following ticket view you
can add comments to the ticket. Since normal users do not log in,
a user cannot automatically display all his or her tickets when he
or she accesses the page again. However, it is possible to find the
created ticket using the Ticket ID `<ticket-id>` by entering the
URL

```text
https://localhost:<Port>/ticket?id=<ticket-id>
```

in the address bar of the browser. In addition, the users receive
an email with a link to their ticket via the mailing service connected
to the Mail API (see :book: [Mail API Reference](Mail-API-Reference.md)
for more information). Registered users can create tickets by
selecting <kbd>Create ticket</kbd> in the navigation.

## Assign a User to a Ticket

This action can only be performed by registered users in the web
interface. For this purpose, a user can be selected from a selection
field for each unassigned ticket under the navigation point <kbd>All
tickets</kbd>. A click on <kbd>Assign</kbd> confirms the selection
of the agent. Editors who have just activated the vacation mode
cannot be selected.

## View Tickets

Logged-in users get the tickets assigned to them displayed on the
dashboard. For each ticket there is also a button which can be used
to open it. Tickets can only be opened by the editor assigned to
them. Tickets that do not yet have an editor can be opened by any
registered user (via the navigation point <kbd>All Tickets</kbd>).

Another method to open a ticket is to type the URL pointing to the
specific ticket into the browser's address bar. To do this, you
need to know the exact ticket id to be inserted as a URL query
parameter. This is described in
:paperclip: [Create a new Ticket](#create-a-new-ticket) above.

## Change the Status of a Ticket

There are three different processing statuses for a ticket:

* **Open** (The ticket is not assigned yet)
* **In Progress** (The ticket is assigned and being processed)
* **Closed** (The ticket has been processed and is closed)

A user can change the status of the tickets assigned to him at any
time by opening the ticket and selecting the status from a selection
list and then confirming his changes with a click on <kbd>Save</kbd>.
In addition, the status is automatically changed for certain events.
If an editor is assigned to a ticket, its status automatically
changes to _In Progress_. If an editor chooses <kbd>Release
ticket</kbd> on his dashboard for one of his tickets, the status
automatically changes to _Open_. If two tickets are merged, the
ticket merged to the other one is closed. And finally, if a new
answer is created using the E-Mail Recipience API and the corresponding
ticket is already closed it is opened again.

## Merge Tickets

In the event that a customer accidentally creates a new ticket
rather than commenting on the existing one, two tickets that share
the same customer (same email) and processor can be merged. To do
this, the agent must open the ticket that was incorrectly created
and select the desired original ticket in the "Merge with" selection
field. If the selection field is empty, the processor has no other
ticket from this customer. After the editor confirms the step with
<kbd>Save</kbd>, the wrongly created ticket is closed and the message
(and all comments) is copied into the original ticket. The ticket
that was merged to the other ticket will be hidden and the comments
will be sorted by creation time.

## Activate the Vacation Mode

Registered users can switch the vacation mode on and off on their
dashboard by clicking on "Activate vacation mode" / "Deactivate
vacation mode". If a user is in vacation mode, new tickets cannot
be assigned to him.

## Remark on the Field _Reply Type_

Whenever you want to change or comment on an open ticket as an
editor, you have to select the field Reply Type. Internal comments
are only displayed to registered users and can be recognized by a
light green background color. External comments are shown with a
light blue background color. In addition, no emails will be sent
to the customer for any changes and comments stored under the reply
type _internal_. The customer will be informed via email (to be
accessed via the Mail API) about changes where the field is set to
_external_.
