/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

/**
 * Holders for divs
 */
const DASHBOARD = document.querySelector("#dashboard");
const CREATE_TICKET = document.querySelector("#create_ticket");
const ALL_TICKETS = document.querySelector("#all_tickets");

/**
 * toggle sets the desired html div visible, while disabling the visibility of the others
 * @param {a} ALL_TICKETS The given a element from the navigation
 */
function toggle(a) {

    let ticket = document.querySelector("#ticket");

    if (ticket) {
        ticket.style.display = "none";
        window.history.replaceState({}, document.title, "/" + "");
    }


    switch (document.querySelector(a.href.substring(a.href.indexOf('#')))) {

        case DASHBOARD:
            DASHBOARD.style.display = "";
            CREATE_TICKET.style.display = "none";
            ALL_TICKETS.style.display = "none";
            break;

        case CREATE_TICKET:
            DASHBOARD.style.display = "none";
            CREATE_TICKET.style.display = "";
            ALL_TICKETS.style.display = "none";
            break;

        case ALL_TICKETS:
            DASHBOARD.style.display = "none";
            CREATE_TICKET.style.display = "none";
            ALL_TICKETS.style.display = "";
            break;
    }
}

/**
 * unassignTicket removes the ticket from a user
 * @param {*} btn Specific button tied to a ticket
 */
function unassignTicket(btn) {

    let id = btn.replace("btn_", "");

    let req = ajaxObject();

    let url = "/unassignTicket?id=" + id;

    req.open("GET", encodeURI(url), true);
    req.onreadystatechange = () => {
        if (req.readyState === 4 && req.status === 200) {
            document.querySelector("#" + btn.replace("btn_", "ticket_")).innerHTML = req.responseText;
        }
    };
    req.send(null);
}

/**
 * assignTicket assigns the ticket in the UI and blocks the ticket it from further manipulation
 * by disabling the button
 * @param {*} btn Given button to specific ticket
 */
function assignTicket(btn) {

    let id = btn.replace("btn_", "");
    let user = document.querySelector("#select_" + id).value;

    let req = ajaxObject();

    let url = "/assignTicket?id=" + id + "&user=" + user;

    req.open("GET", encodeURI(url), true);
    req.onreadystatechange = () => {
        if (req.readyState === 4 && req.status === 200) {
            document.querySelector("#" + btn.replace("btn_", "td_")).innerHTML = req.responseText;
            document.querySelector("#" + btn).disabled = true;
            document.querySelector("#" + btn).style.opacity = 0.25;
            document.querySelector("#" + btn.replace("btn_", "td_status_")).innerHTML = "In Bearbeitung";
        }
    };
    req.send(null);
}

/**
 * Create ajax object, supporting IE as well
 */
function ajaxObject() {

    let activexmodes = ["Msxml2.XMLHTTP", "Microsoft.XMLHTTP"];

    if (window.ActiveXObject) {

        for (var i = 0; i < activexmodes.length; i++) {

            try {
                return new ActiveXObjext(activexmodes[i]);
            }
            catch (e) {
            }
        }
    }
    else if (window.XMLHttpRequest) {
        return new XMLHttpRequest();
    }
    else {
        return false;
    }
}
