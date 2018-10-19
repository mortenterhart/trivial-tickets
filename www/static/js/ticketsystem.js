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
const CREATE_TICKET = ocument.querySelector("#create_ticket");
const ALL_TICKETS   = document.querySelector("#all_tickets");

/**
 * toggle sets the desired html div visible, while disabling the visibility of the others
 * @param {} e The given a element from the navigation
 */
function toggle(e){

    switch(document.querySelector(e.href.substring(e.href.indexOf('#')))){

        case CREATE_TICKET: 
            ALL_TICKETS.style.display   = "none";
            CREATE_TICKET.style.display = "";
            e.style.color               = "#ffffff";
        break;
        case ALL_TICKETS: 
            ALL_TICKETS.style.display   = "";
            CREATE_TICKET.style.display = "none";
        break;
    }
}