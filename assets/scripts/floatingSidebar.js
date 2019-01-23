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
 * Keep the sidebar with the table of contents
 * floating next to the main content
 */

$(function() {
    // Get the window and sidebar elements
    let $window = $(window);
    let sidebar = $('#sidebar');
    let offset = sidebar.offset();

    // Define a top padding that is used as
    // distance to the top screen border.
    let topPadding = 17;

    // Register a scroll function to the window that
    // moves the sidebar up and down
    $window.scroll(function() {
        if ($window.scrollTop() > offset.top) {
            sidebar.stop().animate({
                marginTop: $window.scrollTop() - offset.top + topPadding
            });
        } else {
            sidebar.stop().animate({
                marginTop: 0
            });
        }
    });
});