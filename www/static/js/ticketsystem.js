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
const DASHBOARD     = document.querySelector("#dashboard");
const CREATE_TICKET = document.querySelector("#create_ticket");
const ALL_TICKETS   = document.querySelector("#all_tickets");

/**
 * toggle sets the desired html div visible, while disabling the visibility of the others
 * @param {} ALL_TICKETS The given a element from the navigation
 */
function toggle(a) {
    switch(document.querySelector(a.href.substring(a.href.indexOf('#')))){

        case DASHBOARD: 
            DASHBOARD.style.display     = "";
            CREATE_TICKET.style.display = "none";
            ALL_TICKETS.style.display   = "none";  
            break;

        case CREATE_TICKET: 
            DASHBOARD.style.display     = "none";
            CREATE_TICKET.style.display = "";
            ALL_TICKETS.style.display   = "none";    
            break;
            
        case ALL_TICKETS: 
            DASHBOARD.style.display     = "none";
            CREATE_TICKET.style.display = "none";
            ALL_TICKETS.style.display   = "";  
            break;
    }
}