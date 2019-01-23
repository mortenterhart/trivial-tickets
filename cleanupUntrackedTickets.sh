#!/usr/bin/env bash
##
## Trivial Tickets Ticketsystem
## Copyright (C) 2019 The Contributors
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published by
## the Free Software Foundation, either version 3 of the License, or
## (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.
##
##
## Ticketsystem Trivial Tickets
##
## Matriculation numbers: 3040018, 6694964, 3478222
## Lecture:               Programmieren II, INF16B
## Lecturer:              Herr Prof. Dr. Helmut Neemann
## Institute:             Duale Hochschule Baden-Württemberg Mosbach
##
## ---------------
## Script to delete untracked tickets and mails
## to revert the server into its initial state
##
## Get more information about this script here:
## https://github.com/mortenterhart/trivial-tickets/wiki/Cleaning-untracked-Tickets-and-Mails
##

# Script constants
program_name="trivial-tickets"
script_name="${0##*/}"
root_dir="${0%/*}"

# Paths to the ticket and mail directories
ticket_dir="${root_dir}/files/tickets"
mail_dir="${root_dir}/files/mails"

# Script options
clean_options=()
force_removal="false"
dry_run="false"
interactive="false"
exclude_pattern=""

# info prints an info message consisting
# of all arguments concatenated to stdout.
function info() {
    printf "[%s] [INFO] %s\n" "${program_name}" "$*"
}

# error prints an error message consisting
# of all arguments concatenated to stderr.
function error() {
    printf "%s: %s\n" "${script_name}" "$*" >&2
    printf "Try '%s --help' for more information.\n" "${script_name}" >&2
}

# helpText prints an informative help text
# with options and usage instructions.
function helpText() {
    printf "Usage: %s [option(s)]\n" "${script_name}"
    printf "Delete untracked tickets and mails from cache\n"
    printf "in 'files/tickets' and 'files/mails' to revert\n"
    printf "the server to its original state.\n\n"
    printf "%s offers the following options:\n" "${script_name}"
    printf "  -n, --dry-run     Show which files would be deleted,\n"
    printf "                    but don't delete actually any.\n"
    printf "  -e, --exclude=<PATTERN>\n"
    printf "                    Exclude files matching the pattern\n"
    printf "                    from deletion.\n"
    printf "  -f, --force       Force deletion without prompting for.\n"
    printf "  -h, --help        Display this help text and exit.\n"
    printf "  -i, --interactive Delete files using an interactive menu.\n"
    printf "  -m, --mail-dir=<DIR>\n"
    printf "                    Specify another directory for mails.\n"
    printf "                    Defaults to \"files/mails\".\n"
    printf "  -q, --quiet       Don't print the file paths as they\n"
    printf "                    are removed.\n"
    printf "  -t, --ticket-dir=<DIR>\n"
    printf "                    Specify another directory for tickets.\n"
    printf "                    Defaults to \"files/tickets\".\n"
    printf "Mandatory arguments to long options are also mandatory\n"
    printf "for short options. Long options may be abbreviated as\n"
    printf "long as they do not get ambiguous.\n"
}

function testDirectoryExists() {
    local directory_name="$1"

    if ! [ -d "${directory_name}" ]; then
        error "directory does not exist: '${directory_name}'"
        exit 2
    fi
}

# processOptions parses all command-line options and
# sets the appropriate options for this script. The
# script exits with an error in case an option is
# not defined, requires an argument or is ambiguous.
function processOptions() {
    local clean_mode_set="false"

    local option
    OPTIND=1
    while getopts ":efhimnqt:-:" option "$@"; do
        case "${option}" in
            e)
                # -e, --exclude=<pattern>
                clean_options+=("-e" "${OPTARG}")
                exclude_pattern="${OPTARG}"
                ;;
            f)
                # -f, --force
                force_removal="true"
                ;;
            h)
                # -h, --help
                helpText
                exit 0
                ;;
            i)
                # -i, --interactive
                clean_options+=("-i")
                interactive="true"
                clean_mode_set="true"
                ;;
            m)
                # -m, --mail-dir
                testDirectoryExists "${OPTARG}"
                mail_dir="${OPTARG}"
                ;;
            n)
                # -n, --dry-run
                clean_options+=("-n")
                dry_run="true"
                clean_mode_set="true"
                ;;
            q)
                # -q, --quiet
                clean_options+=("-q")
                ;;
            t)
                # -t, --ticket-dir
                testDirectoryExists "${OPTARG}"
                ticket_dir="${OPTARG}"
                ;;
            -)
                # Parsing of long options
                local option="${OPTARG%%=*}"
                local option_argument="${OPTARG#*=}"

                local long_options="dry-run exclude force help interactive mail-dir quiet ticket-dir"
                local long_options_arguments="exclude mail-dir ticket-dir"

                local -a option_matches=($(compgen -W "${long_options}" -- "${option}"))

                if [ "${#option_matches[@]}" -eq 0 ]; then
                    error "unrecognized option: --${option}"
                    exit 1
                elif [ "${#option_matches[@]}" -eq 1 ]; then
                    if compgen -W "${long_options_arguments}" -- "${option_matches[0]}" >/dev/null 2>&1 && [[ "${OPTARG}" != *=* ]]; then
                        error "option requires an argument: --${option_matches[0]}"
                        exit 2
                    elif ! compgen -W "${long_options_arguments}" -- "${option_matches[0]}" >/dev/null 2>&1 && [[ "${OPTARG}" == *=* ]]; then
                        error "option does not allow an argument: --${option_matches[0]}"
                        exit 2
                    fi

                    case "${option_matches[0]}" in
                        dry-run)
                            clean_options+=("--dry-run")
                            dry_run="true"
                            clean_mode_set="true"
                            ;;
                        exclude)
                            clean_options+=("--exclude" "${option_argument}")
                            exclude_pattern="${option_argument}"
                            ;;
                        force)
                            force_removal="true"
                            ;;
                        help)
                            helpText;
                            exit 0
                            ;;
                        interactive)
                            clean_options+=("--interactive")
                            interactive="true"
                            clean_mode_set="true"
                            ;;
                        mail-dir)
                            testDirectoryExists "${option_argument}"
                            mail_dir="${option_argument}"
                            ;;
                        quiet)
                            clean_options+=("--quiet")
                            ;;
                        ticket-dir)
                            testDirectoryExists "${option_argument}"
                            ticket_dir="${option_argument}"
                            ;;
                    esac
                else
                    error "ambiguous option: '--${option}'; could be one of:"
                    printf "  --%s\n" "${option_matches[@]}"
                    exit 3
                fi
                ;;
            \?)
                # invalid option or '?'
                if [ "${OPTARG}" == "?" ]; then
                    helpText
                    exit 0
                fi

                error "invalid option -- ${OPTARG}"
                exit 1
                ;;
            :)
                # option requires an argument
                error "option requires an argument -- ${OPTARG}"
                exit 2
                ;;
        esac
    done

    # If no clean mode was specified assume '-f' because
    # the 'git clean' step requires one of the mode options
    # to be set
    if ! "${clean_mode_set}"; then
        clean_options+=("-f")
    fi

    # Shift the processed options so that only
    # non-option parameters are left
    shift "$((OPTIND - 1))"
}

# promptMenu prints the available commands and their
# keys for the prompt menu.
function promptMenu() {
    info "Commands:  [y] Continue deletion   [n] Show which files would be removed"
    info "           [*] Abort deletion      [q] Quit"
}

# Parse and process the command-line options
processOptions "$@"

info "You attempted to remove all untracked tickets and mails from these directories:"
info "  ${ticket_dir}"
info "  ${mail_dir}"

# Collect the untracked ticket and mail files from
# 'files/tickets' and 'files/mails' and store them in an
# array. Exclude the files matching the exclude pattern
# from search if one was given.
untracked_files=()
mapfile -t untracked_files < <(git -C "${root_dir}" ls-files --ignored --exclude-standard --others \
    $(if [ -n "${exclude_pattern}" ]; then echo "--exclude" "${exclude_pattern}"; fi) \
    "${ticket_dir}/*.json" "${mail_dir}/*.json" 2>/dev/null)

# Only start deletion process if untracked files were found
if [ "${#untracked_files[@]}" -gt 0 ]; then

    info "In total, ${#untracked_files[@]} untracked file(s) found to be removed."

    continue_deletion="y"

    # Show the prompt as long as no terminating command
    # was entered
    prompt_finished="false"
    while ! "${prompt_finished}"; do

        # Show the command prompt if none of --force, --dry-run or
        # --interactive were specified. In those cases skip the prompt
        # and execute the clean command directly.
        if ! "${force_removal}" && ! "${dry_run}" && ! "${interactive}"; then
            promptMenu
            read -rep "$(info "Command: ")" continue_deletion
        fi

        case "${continue_deletion}" in
            n)
                # [n] (dry run): Show which files would be removed and
                # return to the prompt
                info "These untracked files would be deleted by cleanup:"
                printf "\n"
                printf "%s\n" "${untracked_files[@]}"
                printf "\n"
                ;;
            q)
                # [q] (quit): Quit the program without deleting anything
                info "Quitting, no cleanup done."
                prompt_finished="true"
                ;;
            y)
                # [y] (continue): Continue execution and cleanup all
                # untracked tickets and mails. If --dry-run or --interactive
                # were passed display the files or the interactive menu.
                if "${dry_run}" || "${interactive}"; then
                    info "These untracked files would be deleted by cleanup:"
                else
                    info "Removing all untracked tickets and mails:"
                fi
                printf "\n"

                # Execute 'git clean' to manage the cleanup process
                git -C "${root_dir}" clean -X "${clean_options[@]}" -- "${ticket_dir}/*.json" "${mail_dir}/*.json"

                prompt_finished="true"
                ;;
            *)
                # [*] (abort): Any other user input aborts the deletion
                info "Aborted deletion of tickets and mails."
                info "Type 'y' explicitly to confirm."
                prompt_finished="true"
                ;;
        esac
    done
else
    info "No untracked files found and thus no removal done."
fi

exit 0
