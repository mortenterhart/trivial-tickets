/*
 * Trivial Tickets Ticketsystem
 * Copyright (C) 2019 The Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 6694964, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 * JavaScript resources
 */

/**
 * Holders for <div> elements
 */
const DASHBOARD = document.querySelector("#dashboard");
const CREATE_TICKET = document.querySelector("#create_ticket");
const ALL_TICKETS = document.querySelector("#all_tickets");

/**
 * toggleVisibility sets the desired html <div> visible, while disabling the visibility
 * of the others.
 * @param {HTMLLinkElement} a The given <a> element from the navigation
 */
function toggleVisibility(a) {

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
 * unassignTicket releases the specific ticket from a user.
 * @param {String} btn The specific button id tied to a ticket
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
 * assignTicket assigns the ticket in the UI and blocks the ticket from further manipulation
 * by disabling the button.
 * @param {String} btn The button id of the specific ticket
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
            document.querySelector("#" + btn).style.opacity = "0.25";
            document.querySelector("#" + btn.replace("btn_", "td_status_")).innerHTML = "In Progress";
        }
    };
    req.send(null);
}

/**
 * Create an Ajax object, supporting Internet Explorer as well.
 * @return {XMLHttpRequest|ActiveXObject} the created Ajax object
 */
function ajaxObject() {

    let activeXModes = ["Msxml2.XMLHTTP", "Microsoft.XMLHTTP"];

    if (window.ActiveXObject) {

        for (let mode in activeXModes) {

            try {
                return new ActiveXObjext(mode);
            }
            catch (error) {
                console.error(error);
            }
        }
    }
    else if (window.XMLHttpRequest) {
        return new XMLHttpRequest();
    }

    return null;
}
