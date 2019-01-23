---
title: Mail API Reference
layout: wiki
permalink: wiki/Mail-API-Reference

---

# Mail API Reference

**Table of Contents**

* [E-Mail Recipience API: Creating tickets and answers](#e-mail-recipience-api-creating-tickets-and-answers)
* [E-Mail Dispatch API: Fetching Mails](#e-mail-dispatch-api-fetching-mails)
* [E-Mail Dispatch API: Verifying sent Mails](#e-mail-dispatch-api-verifying-sent-mails)

## E-Mail Recipience API: Creating tickets and answers

Messages can be fed into the system via a REST interface of the ticket system,
which can be addressed by an external service via an HTTPS request. The external
service, in this case our command-line tool, cannot receive genuine emails, but
it can create emails and send them as JSON formatted to the ticket creation API
(see [`api_in.ReceiveMail`](https://github.com/mortenterhart/trivial-tickets/blob/master/api/api_in/api_in.go#L98)).

The API expects the parameters `from` (the sender's email address), `subject`
(the subject of the ticket) and `message` (the actual message). These are
meticulously searched for inconsistencies (e.g. missing or too many parameters)
and trigger an error if necessary. Each erroneous query results in a return
status (response code) of `400 Bad Request`. In addition, the system checks
whether the sender e-mail address is valid.

If the request was made correctly, a new ticket with the respective parameters
is created and an OK status is returned to the caller as JSON. To confirm the
creation, a noreply e-mail is created for the sender so that he can be sure that
his ticket has been created. In this mail he will also find a link to the ticket,
because only logged in users can see a list of tickets. There is also a mailto
link in the mail that can be used to create an email that can create a new
response to the created ticket. For this, the subject of the email must have a
certain syntax:

```text
[Ticket "<Ticket id>"] <Ticket subject>
```

If this syntax is given in the subject and the ticket id exists, the message is
appended to this ticket. If the ticket had the status "Closed", it is reset to
"Open". Also in this case a notification email will be sent to the sender. A
further advantage of this method is that an email ping pong can be prevented,
since automatic reply messages (e.g. from absence assistants) are sent to the
noreply email address and thus no reply is sent.

### Request Information and Headers

The API is accessible under the following URL with the specified HTTP method and
content type.

| **Property** |     **Value**.     |
| :----------: | :----------------: |
| HTTP Method  |       `POST`       |
|     URL      |   `/api/receive`   |
| Content-Type | `application/json` |

**Resulting URL**: `https://localhost:<PORT>/api/receive` (Replace `<PORT>` with
the port the server is listening to)

### Sample HTTP request

```http
POST /api/receive HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 133
Content-Type: application/json
Host: localhost:8443

{
    "from": "customer@example.com",
    "subject": "Example API request",
    "message": "Showing an example call to the E-Mail Recipience API"
}
```

### Request Parameters

The parameters have to be formatted as JSON (see the
:paperclip: [Example Request](#example-request) below). The request is examined
in terms of checking for all required parameters, checking that there are no
additional parameters than the defined ones, checking the parameter data types
and checking the validity of the supplied email address. If any of these checks
is unsuccessful the request is considered to be invalid and is aborted with a
`400 Bad Request` response status.

| **Parameter** | **Type** | **Required** | **Description**                                                                                                                                                                                    |
| :-----------: | :------: | :----------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|    `from`     |  string  |   required   | The sender's email address (usually the customer's email address). The email address is checked to be valid.                                                                                       |
|   `subject`   |  string  |   required   | The subject attached to the new ticket. If the subject has the certain syntax described above and if the applied ticket id exists a new answer is created to this ticket rather than a new ticket. |
|   `message`   |  string  |   required   | The message for the ticket or answer.                                                                                                                                                              |

### Response Parameters

If the request was successful or if the used method is not permitted by the API
the response is in the JSON format. The `Content-Type` header is set to
`application/json` accordingly. If the request is invalid or if an internal
server error occurred (e.g. the server could not create the ticket file) the
response is given as plain text with the response status and an error message.

A JSON response consists of the following parameters:

| **Parameter** | **Type** | **Description**                           |
| :-----------: | :------: | :---------------------------------------- |
|   `status`    | integer  | The HTTP response status code.            |
|   `message`   |  string  | A message describing the response status. |

### Response Statuses

The following HTTP statuses can be returned by the API.

| **Status**                  | **Reason**                                                                                                                                                                                                                          |
| :-------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `200 OK`                    | The request was processed successfully. A new ticket or answer has been created alongside with an email informing about the change.                                                                                                 |
| `400 Bad Request`           | The request was invalid. Possible reasons: invalid JSON syntax, missing required parameters, additional parameters, mismatching parameter data types or invalid email address. Note that the response is `text/plain` in this case. |
| `405 Method Not Allowed`    | The used method other than `POST` is not permitted. Use `POST` instead.                                                                                                                                                             |
| `500 Internal Server Error` | There was a problem while reading the request body or while writing the ticket file. Note that the response is `text/plain` in this case.                                                                                           |

### Example Request

```bash
curl -v --insecure -X POST https://localhost:8443/api/receive -d '{
    "from": "customer@example.com",
    "subject": "Example API request",
    "message": "Showing an example call to the E-Mail Recipience API"
}'
```

```http
POST /api/receive HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 133
Content-Type: application/json
Host: localhost:8443

{
    "from": "customer@example.com",
    "subject": "Example API request",
    "message": "Showing an example call to the E-Mail Recipience API"
}

HTTP/1.1 200 OK
Content-Length: 43
Content-Type: application/json; charset=utf-8
Date: Wed, 23 Feb 2019 10:27:25 GMT

{
    "message": "OK",
    "status": 200
}
```

## E-Mail Dispatch API: Fetching Mails

Another API provided by the ticket system is the Email Dispatch API. The
external service actively requests emails created at the server and can then
send them. It verifies the successful sending of the email with a confirmation
to the server, which can finally remove the email. With the
[`api_out.SendMail`](https://github.com/mortenterhart/trivial-tickets/blob/master/api/api_out/api_out.go#L58)
function, the Trivial Tickets ticket system offers a function that is called at
various so-called mail events and then creates an email. Mail events include the
creation and update of tickets, the creation of new comments, the assignment of
an agent and the release of a ticket (see package
[`mail_events`](https://github.com/mortenterhart/trivial-tickets/blob/master/mail_events/mail_events.go)).
With the help of associated mail templates for all events, the notification
emails are enriched with the information from the respective ticket and the
email is saved in a file. Using the API function
[`api_out.FetchMails`](https://github.com/mortenterhart/trivial-tickets/blob/master/api/api_out/api_out.go#L93)
the generated emails can be retrieved by a GET request.

The E-Mail Dispatch API is primarily used to retrieve the emails generated by
the server. The server creates new emails on every action refering to a ticket.
Using the FetchMails API a client can collect these emails remaining to be sent
and send them to the respective recipient. The request is a simple `GET` request
without any parameters. Note that the request body is not validated against
emptiness, but is ignored completely. The response is returned as JSON formatted
emails with information about the sender, the mail id, the subject and the message.

### Request Information and Headers

The API uses the following URL and method whereas other methods are not permitted and are rejected with a response status of `405 Method Not Allowed`.

| **Property** |     **Value**     |
| :----------: | :---------------: |
|     URL      | `/api/fetchMails` |
|    Method    |       `GET`       |

Resource URL: `https://localhost:<PORT>/api/fetchMails` (Replace `<PORT>` with
the port the server is listening to)

### Request Parameters

No parameters accepted. Any request parameters are ignored. See the
:paperclip: [Example Request](#example-request-1) below.

### Response Parameters

The server response contains the current emails ready to be sent in JSON format.
The email contents are mapped to their unique id and consist of the sender's
email address `from`, the targeting email address `to`, the mail id `id`, the
subject `subject` and the message `message`. Multiple emails are concatenated in
the JSON object and separated by commas. Please note that the emails are not
returned in the order of creation, but sorted alphabetically after their mail
ids. If the server does not host recently created emails the response is an
empty JSON object `{}`. If the server faces an error converting the mails to
JSON the response is the HTTP status and the error message in plain text.

A response with two emails can look like the following:

```json
{
    "GdDEFr5Bg3": {
        "from": "no-reply@trivial-tickets.com",
        "to": "customer@example.com",
        "id": "GdDEFr5Bg3",
        "subject": "[trivial-tickets] Example API request",
        "message": "Sehr geehrter Kunde, sehr geehrte Kundin,\n\nIhr Ticket 'jEOxHA5dVP' ist erfolgreich erstellt worden.\nWenn Sie eine neuen Kommentar zu diesem Ticket schreiben wollen,\nnutzen Sie bitte den folgenden Link: mailto:support@trivial-tickets.com?subject=%5BTicket%20%22jEOxHA5dVP%22%5D%20Example%20API%20request\n-----------------------------\nKunde:      customer@example.com\nSchlüssel:  jEOxHA5dVP\nURL:        https://localhost:8443/ticket?id=jEOxHA5dVP\nBearbeiter: kein Bearbeiter zugewiesen\nStatus:     Offen\n\nBetreff: Example API request\n\nShowing an example call to the E-Mail Recipience API\n\n-----------------------------\n\nMit freundlichen Grüßen\nIhr trivial-tickets Team\n\nDiese Meldung wurde automatisch durch trivial-tickets.com generiert.\nBitte antworten Sie nicht auf diese E-Mail."
    },
    "Vs8KaELgC7": {
        "from": "no-reply@trivial-tickets.com",
        "to": "customer@example.com",
        "id": "Vs8KaELgC7",
        "subject": "[trivial-tickets] Example API request",
        "message": "Sehr geehrter Kunde, sehr geehrte Kundin,\n\nIhr Ticket 'plvGRdlTgj' ist erfolgreich erstellt worden.\nWenn Sie eine neuen Kommentar zu diesem Ticket schreiben wollen,\nnutzen Sie bitte den folgenden Link: mailto:support@trivial-tickets.com?subject=%5BTicket%20%22plvGRdlTgj%22%5D%20Example%20API%20request\n-----------------------------\nKunde:      customer@example.com\nSchlüssel:  plvGRdlTgj\nURL:        https://localhost:8443/ticket?id=plvGRdlTgj\nBearbeiter: kein Bearbeiter zugewiesen\nStatus:     Offen\n\nBetreff: Example API request\n\nShowing an example call to the E-Mail Recipience API\n\n-----------------------------\n\nMit freundlichen Grüßen\nIhr trivial-tickets Team\n\nDiese Meldung wurde automatisch durch trivial-tickets.com generiert.\nBitte antworten Sie nicht auf diese E-Mail."
    }
}
```

### Response Statuses

The following HTTP statuses can be returned by the API.

| **Status**                  | **Reason**                                                                                                          |
| :-------------------------- | :------------------------------------------------------------------------------------------------------------------ |
| `200 OK`                    | The request could be processed succesfully. The remaining emails to be sent are returned as JSON.                   |
| `405 Method Not Allowed`    | The used method other than `GET` is not permitted. Use `GET` instead.                                               |
| `500 Internal Server Error` | An unexpected error occurred while building the JSON response. Note that the response is `text/plain` in this case. |

### Example Request

```bash
curl -v --insecure -X GET https://localhost:8443/api/fetchMails
```

```http
GET /api/fetchMails HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8443

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 05 Feb 2019 11:22:38 GMT
Transfer-Encoding: chunked

{
    "CvjHp0bbVI": {
        "from": "no-reply@trivial-tickets.com",
        "id": "CvjHp0bbVI",
        "message": "Sehr geehrter Kunde, sehr geehrte Kundin,\n\nunser Mitarbeiter 'Boris Floricic' bearbeitet nun Ihr Ticket:\n\n-----------------------------\nKunde:      worried-person@example.com\nSchlüssel:  aS11SPZKRb\nURL:        https://localhost:8443/ticket?id=aS11SPZKRb\nBearbeiter: Boris Floricic (boris.floricic@trivial-tickets.com)\nStatus:     In Bearbeitung\n\nBetreff: Warning: Cross Site Scripting Attempt\n\nI saw a hacker on your site trying to exploit the forms on your website with Cross Site Scripting (XSS).\r\n\r\nPlease fix those security issues on your site.\n\n-----------------------------\n\nMit freundlichen Grüßen,\nIhr Trivial Tickets Team\n\nDiese Meldung wurde automatisch durch trivial-tickets.com generiert.\nBitte antworten Sie nicht auf diese E-Mail.",
        "subject": "[trivial-tickets] Warning: Cross Site Scripting Attempt",
        "to": "worried-person@example.com"
    },
    "HlJfbVWtu6": {
        "from": "no-reply@trivial-tickets.com",
        "id": "HlJfbVWtu6",
        "message": "Sehr geehrter Kunde, sehr geehrte Kundin,\n\nIhr Ticket 'aS11SPZKRb' ist erfolgreich erstellt worden.\nWenn Sie eine neuen Kommentar zu diesem Ticket schreiben wollen,\nnutzen Sie bitte den folgenden Link: mailto:support@trivial-tickets.com?subject=%5BTicket%20%22aS11SPZKRb%22%5D%20Warning:%20Cross%20Site%20Scripting%20Attempt\n-----------------------------\nKunde:      worried-person@example.com\nSchlüssel:  aS11SPZKRb\nURL:        https://localhost:8443/ticket?id=aS11SPZKRb\nBearbeiter: kein Bearbeiter zugewiesen\nStatus:     Offen\n\nBetreff: Warning: Cross Site Scripting Attempt\n\nI saw a hacker on your site trying to exploit the forms on your website with Cross Site Scripting (XSS).\n\nPlease fix those security issues on your site.\n\n-----------------------------\n\nMit freundlichen Grüßen,\nIhr Trivial Tickets Team\n\nDiese Meldung wurde automatisch durch trivial-tickets.com generiert.\nBitte antworten Sie nicht auf diese E-Mail.",
        "subject": "[trivial-tickets] Warning: Cross Site Scripting Attempt",
        "to": "worried-person@example.com"
    }
}
```

The resulting email message looks like this:

```text
Dear customer,

Your ticket 'aS11SPZKRb' was created successfully.
If you want to write a new comment to this ticket
please use the following link: mailto:support@trivial-tickets.com?subject=%5BTicket%20%22aS11SPZKRb%22%5D%20Warning:%20Cross%20Site%20Scripting%20Attempt
-----------------------------
Customer:  worried-person@example.com
Key:       aS11SPZKRb
URL:       https://localhost:8443/ticket?id=aS11SPZKRb
Assignee:  no user assigned
Status:    Open

Subject: Warning: Cross Site Scripting Attempt

I saw a hacker on your site trying to exploit the forms on your website with Cross Site Scripting (XSS).

Please fix those security issues on your site.

-----------------------------

Yours sincerely,
Your Trivial Tickets Team

This notification was generated automatically by trivial-tickets.com.
Please do not respond to this e-mail.
```

## E-Mail Dispatch API: Verifying sent Mails

The server saves the remaining emails in its cache. During this phase a client
can retrieve the generated emails as often as necessary. An external mailing
service (just like the :book: [Trivial Tickets Command-line Tool](CLI-Usage.md))
can load the server-created emails on this way and send them to their respective
recipients. After the sending process the service can confirm the successful
transmission by making a `POST` request to this mail verification API for each
mail by applying the mail-specific id. The server then checks if the requested
id belongs to an existing mail and if it does the corresponding mail is removed
from the server cache. If the mail id does not exist the verification fails and
an appropriate response with the result is returned.

### Request Information and Headers

This API uses the following URL and only permits `POST` requests. It expects
JSON input from the caller in the request.

| **Property** |     **Value**      |
| :----------: | :----------------: |
|     URL      | `/api/verifyMail`  |
|    Method    |       `POST`       |
| Content-Type | `application/json` |

Resource URL: `https://localhost:<PORT>/api/verifyMail` (Replace `<PORT>` with
the port the server is listening to)

### Request Parameters

The request is expected to be in valid JSON format. The mail `id` of a mail that
is expected here can be taken from the response message of the
:paperclip: [Fetch Mails API](#e-mail-dispatch-api-fetching-mails). They are
guaranteed to be valid. See the :paperclip: [Example Request](#example-request-2)
below.

| **Parameter** | **Type** | **Required** | **Description**                                      |
| :-----------: | :------: | :----------: | :--------------------------------------------------- |
|     `id`      |  string  |   required   | The mail id of the mail that is verified to be sent. |

### Response Parameters

A successful request always returns a JSON response with the following parameters.
An invalid request causing a `400 Bad Request` response or a
`500 Internal Server Error` response contains plain text with the error message.

| **Parameter** | **Type** | **Description**                                                                                                                                                                                                                    |
| :-----------: | :------: | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  `verified`   | boolean  | The status if the mail could be verified to exist and could be deleted accordingly.<br>`true`: The email was found on the server and could be deleted from cache.<br>`false`: The email did not exist or has already been deleted. |
|   `status`    | integer  | The HTTP response status code.                                                                                                                                                                                                     |
|   `message`   |  string  | A message describing the verification status.                                                                                                                                                                                      |

### Response Statuses

The following HTTP statuses can be returned by the API.

| **Status**                  | **Reason**                                                                                                                                                                                                                        |
| :-------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `200 OK`                    | The request was successful and a verification status is returned.                                                                                                                                                                 |
| `400 Bad Request`           | The request was invalid due to one of the following reasons: invalid JSON syntax, missing required parameters, any additional parameter or mismatching parameter data types. Note that the response is `text/plain` in this case. |
| `405 Method Not Allowed`    | The used method other than `POST` is not permitted. Use `POST` instead.                                                                                                                                                           |
| `500 Internal Server Error` | An unexpected problem occurred while trying to remove the mail file. Note that the response is `text/plain` in this case.                                                                                                         |

### Example Request

```bash
curl -v --insecure -X POST https://localhost:8443/api/verifyMail -d '{ "id": "tqlxiYe4lY" }'
```

```http
POST /api/verifyMail HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 20
Content-Type: application/json
Host: localhost:8443

{
    "id": "tqlxiYe4lY"
}

HTTP/1.1 200 OK
Content-Length: 113
Content-Type: application/json; charset=utf-8
Date: Tue, 05 Feb 2019 16:07:48 GMT

{
    "message": "mail 'tqlxiYe4lY' was successfully sent and deleted from server cache",
    "status": 200,
    "verified": true
}
```