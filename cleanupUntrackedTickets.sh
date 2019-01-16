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
    printf "for short options too.\n"
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
                ;;
            -)
                case "${OPTARG}" in
                    exclude)
                        error "option requires an argument: --${OPTARG}"
                        helpText >&2
                        exit 2
                        ;;
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
                    exclude=*)
                        clean_options+=("--${OPTARG}")
                        ;;
                    *)
                        error "unrecognized option: --${OPTARG}"
                        helpText >&2
                        exit 1
                        ;;
                esac
                ;;
            \?)
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

processArguments "$@"

shift "$((OPTIND - 1))"

info "You attempted to remove all untracked tickets and mails from these directories:"
info "  ${root_dir}/files/tickets"
info "  ${root_dir}/files/mails"

untracked_files=()
mapfile -t untracked_files < <(git -C "${root_dir}" ls-files --ignored --exclude-standard --others "files/tickets/*.json" "files/mails/*.json")

if [ "${#untracked_files[@]}" -gt 0 ]; then

    info "In total, ${#untracked_files[@]} untracked file(s) found to be removed."

    continue_deletion="y"

    prompt_finished="false"
    while ! "${prompt_finished}"; do

        if ! "${force_removal}" && ! "${dry_run}"; then
            info "Commands:  [y/N] Continue deletion   [n] Show which files would be removed"
            read -rep "$(info "Command: ")" continue_deletion
        fi

        case "${continue_deletion}" in
            n)
                info "These untracked files would be deleted by cleanup:"
                printf "%s\n" "${untracked_files[@]}"
                ;;
            [yY])
                if "${dry_run}"; then
                    info "These untracked files would be deleted by cleanup:"
                else
                    info "Removing all untracked tickets and mails:"
                fi

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
