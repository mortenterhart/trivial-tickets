#!/usr/bin/env bash
#
# Ticketsystem Trivial Tickets
#
# Matriculation numbers: 3040018, 6694964, 3478222
# Lecture:               Programmieren II, INF16B
# Lecturer:              Herr Prof. Dr. Helmut Neemann
# Institute:             Duale Hochschule Baden-Württemberg Mosbach
#
# ---------------
# Script to delete untracked tickets and
# mails to revert into the initial state

program_name="trivial-tickets"
script_name="${0##*/}"
root_dir="${0%/*}"

clean_options=()
force_removal="false"
dry_run="false"
exclude_pattern=""

function info() {
    printf "[%s] INFO %s\n" "${program_name}" "$*"
}

function error() {
    printf "[%s] ERROR %s\n" "${program_name}" "$*" >&2
}

function helpText() {
    printf "usage: %s [option(s)]\n" "${script_name}"
    printf "Delete untracked tickets and mails from cache\n"
    printf "in 'files/tickets' and 'files/mails' to revert\n"
    printf "the server to its original state.\n\n"
    printf "%s offers the following options:\n" "${script_name}"
    printf "  -n, --dry-run     Show which files would be deleted,\n"
    printf "                    but don't delete actually any.\n"
    printf "  -e, --exclude=<pattern>\n"
    printf "                    Exclude files matching the pattern\n"
    printf "                    from deletion.\n"
    printf "  -f, --force       Force deletion without prompting for.\n"
    printf "  -h, --help        Display this help text.\n"
    printf "  -i, --interactive Delete files using an interactive menu.\n"
    printf "  -q, --quiet       Don't print the file paths as they\n"
    printf "                    are removed.\n"
    printf "Mandatory arguments to long options are also mandatory\n"
    printf "for short options too. Long options may be abbreviated as\n"
    printf "long as they do not get ambiguous.\n"
}

function processArguments() {
    local clean_mode_set="false"

    local option
    OPTIND=1
    while getopts ":fhnqie:-:" option "$@"; do
        case "${option}" in
            f)
                force_removal="true"
                ;;
            h)
                helpText
                exit 0
                ;;
            n)
                clean_options+=("-n")
                dry_run="true"
                ;;
            q)
                clean_options+=("-q")
                ;;
            i)
                clean_options+=("-i")
                clean_mode_set="true"
                ;;
            e)
                clean_options+=("-e" "${OPTARG}")
                exclude_pattern="${OPTARG}"
                ;;
            -)
                local option="${OPTARG%%=*}"
                local option_argument="${OPTARG#*=}"

                local long_options="dry-run exclude force help interactive quiet"
                local long_options_arguments="exclude"

                local -a option_matches=($(compgen -W "${long_options}" -- "${option}"))

                if [ "${#option_matches[@]}" -eq 0 ]; then
                    error "unrecognized option: --${option}"
                    helpText >&2
                    exit 1
                elif [ "${#option_matches[@]}" -eq 1 ]; then
                    if compgen -W "${long_options_arguments}" -- "${option_matches[0]}" >/dev/null 2>&1 && [[ "${OPTARG}" != *=* ]]; then
                        error "option requires an argument: --${option_matches[0]}"
                        helpText >&2
                        exit 2
                    elif ! compgen -W "${long_options_arguments}" -- "${option_matches[0]}" >/dev/null 2>&1 && [[ "${OPTARG}" == *=* ]]; then
                        error "option does not allow an argument: --${option_matches[0]}"
                        helpText >&2
                        exit 2
                    fi

                    case "${option_matches[0]}" in
                        force)
                            force_removal="true"
                            clean_mode_set="true"
                            ;;
                        help)
                            helpText;
                            exit 0
                            ;;
                        dry-run)
                            clean_options+=("--dry-run")
                            dry_run="true"
                            clean_mode_set="true"
                            ;;
                        quiet)
                            clean_options+=("--quiet")
                            ;;
                        interactive)
                            clean_options+=("--interactive")
                            clean_mode_set="true"
                            ;;
                        exclude)
                            clean_options+=("--exclude" "${option_argument}")
                            exclude_pattern="${option_argument}"
                            ;;
                    esac
                else
                    error "ambiguous option: --${option}; could be one of:"
                    printf "  --%s\n" "${option_matches[@]}"
                    exit 3
                fi
                ;;
            \?)
                if [ "${OPTARG}" == "?" ]; then
                    helpText
                    exit 0
                fi

                error "invalid option -- ${OPTARG}"
                helpText >&2
                exit 1
                ;;
            :)
                error "option requires an argument -- ${OPTARG}"
                helpText >&2
                exit 2
                ;;
        esac
    done

    if ! "${clean_mode_set}"; then
        clean_options+=("-f")
    fi
}

function promptMenu() {
    info "Commands:  [y/N] Continue deletion   [n] Show which files would be removed"
    info "           [*]   Abort deletion      [q] Quit"
}

processArguments "$@"

shift "$((OPTIND - 1))"

info "You attempted to remove all untracked tickets and mails from these directories:"
info "  ${root_dir}/files/tickets"
info "  ${root_dir}/files/mails"

untracked_files=()
mapfile -t untracked_files < <(git -C "${root_dir}" ls-files --ignored --exclude-standard --others \
    $(if [ -n "${exclude_pattern}" ]; then echo "--exclude" "${exclude_pattern}"; fi) "files/tickets/*.json" "files/mails/*.json")

if [ "${#untracked_files[@]}" -gt 0 ]; then

    info "In total, ${#untracked_files[@]} untracked file(s) found to be removed."

    continue_deletion="y"

    prompt_finished="false"
    while ! "${prompt_finished}"; do

        if ! "${force_removal}" && ! "${dry_run}"; then
            promptMenu
            read -rep "$(info "Command: ")" continue_deletion
        fi

        case "${continue_deletion}" in
            n)
                info "These untracked files would be deleted by cleanup:"
                printf "\n"
                printf "%s\n" "${untracked_files[@]}"
                printf "\n"
                ;;
            q)
                info "Quitting, no cleanup done."
                prompt_finished="true"
                ;;
            [yY])
                if "${dry_run}"; then
                    info "These untracked files would be deleted by cleanup:"
                else
                    info "Removing all untracked tickets and mails:"
                fi
                printf "\n"

                git -C "${root_dir}" clean -X "${clean_options[@]}" -- "files/tickets/*.json" "files/mails/*.json"

                prompt_finished="true"
                ;;
            *)
                info "Aborted deletion of tickets and mails."
                info "Type explicitly 'y' or 'Y' to confirm."
                prompt_finished="true"
                ;;
        esac
    done
else
    info "No untracked files found and thus no removal done."
fi

exit 0
