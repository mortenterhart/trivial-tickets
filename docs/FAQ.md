---
title: FAQ
layout: wiki
permalink: /wiki/FAQ
---

# Frequently asked Questions

The questions are referring to the :book: [Requirements](Requirements.md)
of the lecture's task.

**Q: May external sources be used to create the HTML frontend (especially
     for the fast and nice creation of the frontend Bootstrap 4, Popper.js
     and jQuery)?**

> **A**: Yes, the use of JS/CSS libraries is permitted. However, the
>        quality of the user interface is not included in the evaluation.

**Q: Regarding [9.1](Requirements.md#9-storage): The tickets should be
     saved together with the complete history in a file in the file system.
     Does this mean that the user assignment must be recorded and saved?
     It is unclear to us whether only the ticket itself or a complete
     "chronicle" is desired.**

> **A**: A ticket should be saved together with all comments in one
>        file.

**Q: Regarding [6.3](Requirements.md#6-e-mail-recipience-over-a-rest-api):
     When calling, it should be checked whether the message refers to an
     already existing ticket. The described API of requirement 6 creates
     a new ticket.**

> **A**: If a message is fed into the system via API, it should be
>        appended to the correct ticket. To do this, the corresponding
>        ticket must be searched and found. Only if no ticket can be
>        found, a new ticket should be created. The identification is
>        usually done by an ID in the subject of the message.

**Q: Does [6.3](Requirements.md#6-e-mail-recipience-over-a-rest-api) request
     that, for example, the subjects and descriptions of all existing tickets
     be checked for similarity? To what extent should effort be invested in
     this?**

> **A**: See above! A suitable procedure must be implemented. I leave
>        it up to you to think about it.

**Q: Regarding [7](Requirements.md#7-e-mail-dispatch-over-a-rest-api): E-Mail
     Dispatch over a REST-API: Should this "email queue" be persisted?**

> **A**: No messages should be lost during a restart.

**Q: Regarding [7.1](Requirements.md#7-e-mail-dispatch-over-a-rest-api): There
     should be a function which can be used to retrieve all e-mails that are
     still to be sent. It is unclear to us what is meant by "e-mails to be sent".
     Should agents receive e-mails for managed tickets? Should customers receive
     an e-mail when they receive a new message on a ticket?**

> **A**: If an editor adds a comment to a ticket, it should be possible
>        for the creator of this ticket to receive this comment by e-mail.

**Q: Must the REST APIs be hedged? Is HTTP Basic Auth sufficient for authentication
     or does session management including cookies have to be implemented? Do I have
     to register a user for the ticket system?**

> **A**: HTTP Basic Auth is useless because the password would be sent in plain
>        text. Basic Auth only makes sense if it is used in combination with
>        HTTPS. HTTPS Basic Auth for authentication is sufficient.
