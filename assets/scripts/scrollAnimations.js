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
 * Display or hide certain HTML elements such as the
 * download buttons or the "Go back to Top" button
 * depending on scroll position.
 *
 * Also enable smooth scrolling to anchor positions.
 */

$(function() {
    /**
     * headerHeight is the height of the initial site header
     * in pixels.
     * @type {number}
     */
    const headerHeight = 200;

    /**
     * fadeDuration defines the time interval in seconds of
     * display and hide animations and the scroll animation
     * to fixed anchors.
     * @type {number}
     */
    const fadeDuration = 400;

    /**
     * root is the HTML root element cached in a variable.
     * @type {HTMLElement|jQuery}
     */
    const root = $('html, body');

    /**
     * backToTop is the HTML element containing the back
     * to top button that is displayed if the header with
     * the height defined in `headerHeight` is out of the
     * screen.
     * @type {HTMLElement|jQuery}
     */
    const backToTop = $('#back-to-top');

    /**
     * downloadButtons references the HTML element that
     * encloses the .zip and .tar.gz download buttons.
     * @type {HTMLElement|jQuery}
     */
    const downloadButtons = $('#download-buttons');

    /**
     * displayScrollElements shows or hides the "back to top"
     * button and the download buttons depending on the
     * scroll position. The "back to top" button is faded in
     * below the header, the download buttons are faded out
     * then.
     *
     * This function gets called whenever a scroll event
     * happens.
     */
    function displayScrollElements() {
        if (document.body.scrollTop > headerHeight || document.documentElement.scrollTop > headerHeight) {
            backToTop.fadeIn(fadeDuration);
            downloadButtons.fadeOut(fadeDuration)
        } else {
            backToTop.fadeOut(fadeDuration);
            downloadButtons.fadeIn(fadeDuration)
        }
    }

    // Observe the window scroll events and call the
    // displayScrollElements on each new event.
    $(window).scroll(function() {
        displayScrollElements()
    });

    // Register a click event on the "back to top" button
    // and slowly scroll back to the top of the page on click.
    backToTop.each(function() {
        $(this).click(function() {
            root.animate({scrollTop: 0}, 'slow');
            return false;
        });
    });

    // When a link to a heading anchor is clicked (e.g. from
    // the table of contents), then scroll smoothly to the
    // anchor instead of going to it directly. The click event
    // is registered to each <a>-tag which link refers to an
    // anchor, i.e. the link starts with a '#'.
    $('a[href^="#"]').click(function() {
        let href = $.attr(this, 'href');

        root.animate({
            scrollTop: $(href).offset().top
        }, fadeDuration, function() {
            window.location.hash = href;
        });

        return false;
    });
});
