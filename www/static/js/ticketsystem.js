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
 * DASHBOARD is the HTML element containing the
 * dashboard view.
 * @type {HTMLElement}
 */
const DASHBOARD = document.querySelector("#dashboard");

/**
 * CREATE_TICKET is the HTML element containing the
 * view with the create ticket form.
 * @type {HTMLElement}
 */
const CREATE_TICKET = document.querySelector("#create_ticket");

/**
 * ALL_TICKETS is the HTML element containing the
 * view with all tickets.
 * @type {HTMLElement}
 */
const ALL_TICKETS = document.querySelector("#all_tickets");

/**
 * toggleVisibility sets the desired html <div> visible, while disabling
 * the visibility of the others.
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
 * hideDashboard hides the dashboard <div> container if its
 * visibility is not already set to none. This preserves the
 * ability to show the dashboard component again instead of
 * wrapping it inside a static invisible <div> container.
 * The dashboard can be only be hidden if a user is logged in
 * because otherwise the dashboard is not loaded at runtime.
 * @param {Boolean} isLoggedIn whether a user is currently
 *                             logged in or not
 */
function hideDashboard(isLoggedIn) {
    if (isLoggedIn && DASHBOARD.style.display !== "none") {
        DASHBOARD.style.display = "none"
    }
}

/**
 * unassignTicket releases the specific ticket from a user.
 * @param {String} button The specific button id tied to a ticket
 */
function unassignTicket(button) {

    let id = button.replace("btn_", "");

    let request = createAJAXObject();

    let url = "/unassignTicket?id=" + id;

    request.open("GET", encodeURI(url), true);
    request.onreadystatechange = () => {
        if (request.readyState === 4 && request.status === 200) {
            document.querySelector("#" + button.replace("btn_", "ticket_")).innerHTML = request.responseText;
        }
    };

    request.send(null);
}

/**
 * assignTicket assigns the ticket in the UI and blocks the ticket
 * from further manipulation by disabling the button.
 * @param {String} button The button id of the specific ticket
 */
function assignTicket(button) {

    let id = button.replace("btn_", "");
    let user = document.querySelector("#select_" + id).value;

    let request = createAJAXObject();

    let url = "/assignTicket?id=" + id + "&user=" + user;

    request.open("GET", encodeURI(url), true);
    request.onreadystatechange = () => {
        if (request.readyState === 4 && request.status === 200) {
            document.querySelector("#" + button.replace("btn_", "td_")).innerHTML = request.responseText;
            document.querySelector("#" + button).disabled = true;
            document.querySelector("#" + button).style.opacity = "0.25";
            document.querySelector("#" + button.replace("btn_", "td_status_")).innerHTML = "In Progress";
        }
    };

    request.send(null);
}

/**
 * createAJAXObject creates an AJAX object, supporting Internet
 * Explorer as well.
 * @return {XMLHttpRequest | ActiveXObject} the created AJAX object
 */
function createAJAXObject() {

    let activeXModes = ["Msxml2.XMLHTTP", "Microsoft.XMLHTTP"];

    if (window.ActiveXObject) {

        for (let mode in activeXModes) {

            try {
                return new ActiveXObjext(mode);
            } catch (error) {
                console.error(error);
            }
        }
    } else if (window.XMLHttpRequest) {
        return new XMLHttpRequest();
    }

    return null;
}
