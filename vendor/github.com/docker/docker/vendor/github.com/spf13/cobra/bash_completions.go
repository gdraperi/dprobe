package cobra

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
)

const (
	BashCompFilenameExt     = "cobra_annotation_bash_completion_filename_extentions"
	BashCompCustom          = "cobra_annotation_bash_completion_custom"
	BashCompOneRequiredFlag = "cobra_annotation_bash_completion_one_required_flag"
	BashCompSubdirsInDir    = "cobra_annotation_bash_completion_subdirs_in_dir"
)

func preamble(out io.Writer, name string) error ***REMOVED***
	_, err := fmt.Fprintf(out, "# bash completion for %-36s -*- shell-script -*-\n", name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = fmt.Fprint(out, `
__debug()
***REMOVED***
    if [[ -n $***REMOVED***BASH_COMP_DEBUG_FILE***REMOVED*** ]]; then
        echo "$*" >> "$***REMOVED***BASH_COMP_DEBUG_FILE***REMOVED***"
    fi
***REMOVED***

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__my_init_completion()
***REMOVED***
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
***REMOVED***

__index_of_word()
***REMOVED***
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
***REMOVED***

__contains_word()
***REMOVED***
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
***REMOVED***

__handle_reply()
***REMOVED***
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ $***REMOVED***#must_have_one_flag[@]***REMOVED*** -ne 0 ]; then
                allflags=("$***REMOVED***must_have_one_flag[@]***REMOVED***")
            else
                allflags=("$***REMOVED***flags[*]***REMOVED*** $***REMOVED***two_word_flags[*]***REMOVED***")
            fi
            COMPREPLY=( $(compgen -W "$***REMOVED***allflags[*]***REMOVED***" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "$***REMOVED***COMPREPLY[0]***REMOVED***" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="$***REMOVED***cur%%=****REMOVED***"
                __index_of_word "$***REMOVED***flag***REMOVED***" "$***REMOVED***flags_with_completion[@]***REMOVED***"
                if [[ $***REMOVED***index***REMOVED*** -ge 0 ]]; then
                    COMPREPLY=()
                    PREFIX=""
                    cur="$***REMOVED***cur#*=***REMOVED***"
                    $***REMOVED***flags_completion[$***REMOVED***index***REMOVED***]***REMOVED***
                    if [ -n "$***REMOVED***ZSH_VERSION***REMOVED***" ]; then
                        # zfs completion needs --flag= prefix
                        eval "COMPREPLY=( \"\$***REMOVED***COMPREPLY[@]/#/$***REMOVED***flag***REMOVED***=***REMOVED***\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __index_of_word "$***REMOVED***prev***REMOVED***" "$***REMOVED***flags_with_completion[@]***REMOVED***"
    if [[ $***REMOVED***index***REMOVED*** -ge 0 ]]; then
        $***REMOVED***flags_completion[$***REMOVED***index***REMOVED***]***REMOVED***
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ $***REMOVED***cur***REMOVED*** != "$***REMOVED***words[cword]***REMOVED***" ]]; then
        return
    fi

    local completions
    completions=("$***REMOVED***commands[@]***REMOVED***")
    if [[ $***REMOVED***#must_have_one_noun[@]***REMOVED*** -ne 0 ]]; then
        completions=("$***REMOVED***must_have_one_noun[@]***REMOVED***")
    fi
    if [[ $***REMOVED***#must_have_one_flag[@]***REMOVED*** -ne 0 ]]; then
        completions+=("$***REMOVED***must_have_one_flag[@]***REMOVED***")
    fi
    COMPREPLY=( $(compgen -W "$***REMOVED***completions[*]***REMOVED***" -- "$cur") )

    if [[ $***REMOVED***#COMPREPLY[@]***REMOVED*** -eq 0 && $***REMOVED***#noun_aliases[@]***REMOVED*** -gt 0 && $***REMOVED***#must_have_one_noun[@]***REMOVED*** -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "$***REMOVED***noun_aliases[*]***REMOVED***" -- "$cur") )
    fi

    if [[ $***REMOVED***#COMPREPLY[@]***REMOVED*** -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    __ltrim_colon_completions "$cur"
***REMOVED***

# The arguments should be in the form "ext1|ext2|extn"
__handle_filename_extension_flag()
***REMOVED***
    local ext="$1"
    _filedir "@($***REMOVED***ext***REMOVED***)"
***REMOVED***

__handle_subdirs_in_dir_flag()
***REMOVED***
    local dir="$1"
    pushd "$***REMOVED***dir***REMOVED***" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
***REMOVED***

__handle_flag()
***REMOVED***
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=$***REMOVED***words[c]***REMOVED***
    local flagvalue
    # if the word contained an =
    if [[ $***REMOVED***words[c]***REMOVED*** == *"="* ]]; then
        flagvalue=$***REMOVED***flagname#*=***REMOVED*** # take in as flagvalue after the =
        flagname=$***REMOVED***flagname%%=****REMOVED*** # strip everything after the =
        flagname="$***REMOVED***flagname***REMOVED***=" # but put the = back
    fi
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: looking for $***REMOVED***flagname***REMOVED***"
    if __contains_word "$***REMOVED***flagname***REMOVED***" "$***REMOVED***must_have_one_flag[@]***REMOVED***"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __contains_word "$***REMOVED***flagname***REMOVED***" "$***REMOVED***local_nonpersistent_flags[@]***REMOVED***"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    if [ -n "$***REMOVED***flagvalue***REMOVED***" ] ; then
        flaghash[$***REMOVED***flagname***REMOVED***]=$***REMOVED***flagvalue***REMOVED***
    elif [ -n "$***REMOVED***words[ $((c+1)) ]***REMOVED***" ] ; then
        flaghash[$***REMOVED***flagname***REMOVED***]=$***REMOVED***words[ $((c+1)) ]***REMOVED***
    else
        flaghash[$***REMOVED***flagname***REMOVED***]="true" # pad "true" for bool flag
    fi

    # skip the argument to a two word flag
    if __contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***two_word_flags[@]***REMOVED***"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

***REMOVED***

__handle_noun()
***REMOVED***
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    if __contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***must_have_one_noun[@]***REMOVED***"; then
        must_have_one_noun=()
    elif __contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***noun_aliases[@]***REMOVED***"; then
        must_have_one_noun=()
    fi

    nouns+=("$***REMOVED***words[c]***REMOVED***")
    c=$((c+1))
***REMOVED***

__handle_command()
***REMOVED***
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"

    local next_command
    if [[ -n $***REMOVED***last_command***REMOVED*** ]]; then
        next_command="_$***REMOVED***last_command***REMOVED***_$***REMOVED***words[c]//:/__***REMOVED***"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_$(basename "$***REMOVED***words[c]//:/__***REMOVED***")"
        else
            next_command="_$***REMOVED***words[c]//:/__***REMOVED***"
        fi
    fi
    c=$((c+1))
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: looking for $***REMOVED***next_command***REMOVED***"
    declare -F $next_command >/dev/null && $next_command
***REMOVED***

__handle_word()
***REMOVED***
    if [[ $c -ge $cword ]]; then
        __handle_reply
        return
    fi
    __debug "$***REMOVED***FUNCNAME[0]***REMOVED***: c is $c words[c] is $***REMOVED***words[c]***REMOVED***"
    if [[ "$***REMOVED***words[c]***REMOVED***" == -* ]]; then
        __handle_flag
    elif __contains_word "$***REMOVED***words[c]***REMOVED***" "$***REMOVED***commands[@]***REMOVED***"; then
        __handle_command
    elif [[ $c -eq 0 ]] && __contains_word "$(basename "$***REMOVED***words[c]***REMOVED***")" "$***REMOVED***commands[@]***REMOVED***"; then
        __handle_command
    else
        __handle_noun
    fi
    __handle_word
***REMOVED***

`)
	return err
***REMOVED***

func postscript(w io.Writer, name string) error ***REMOVED***
	name = strings.Replace(name, ":", "__", -1)
	_, err := fmt.Fprintf(w, "__start_%s()\n", name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = fmt.Fprintf(w, `***REMOVED***
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __my_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("%s")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __handle_word
***REMOVED***

`, name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = fmt.Fprintf(w, `if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_%s %s
else
    complete -o default -o nospace -F __start_%s %s
fi

`, name, name, name, name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, err = fmt.Fprintf(w, "# ex: ts=4 sw=4 et filetype=sh\n")
	return err
***REMOVED***

func writeCommands(cmd *Command, w io.Writer) error ***REMOVED***
	if _, err := fmt.Fprintf(w, "    commands=()\n"); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() || c == cmd.helpCommand ***REMOVED***
			continue
		***REMOVED***
		if _, err := fmt.Fprintf(w, "    commands+=(%q)\n", c.Name()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := fmt.Fprintf(w, "\n")
	return err
***REMOVED***

func writeFlagHandler(name string, annotations map[string][]string, w io.Writer) error ***REMOVED***
	for key, value := range annotations ***REMOVED***
		switch key ***REMOVED***
		case BashCompFilenameExt:
			_, err := fmt.Fprintf(w, "    flags_with_completion+=(%q)\n", name)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			if len(value) > 0 ***REMOVED***
				ext := "__handle_filename_extension_flag " + strings.Join(value, "|")
				_, err = fmt.Fprintf(w, "    flags_completion+=(%q)\n", ext)
			***REMOVED*** else ***REMOVED***
				ext := "_filedir"
				_, err = fmt.Fprintf(w, "    flags_completion+=(%q)\n", ext)
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case BashCompCustom:
			_, err := fmt.Fprintf(w, "    flags_with_completion+=(%q)\n", name)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if len(value) > 0 ***REMOVED***
				handlers := strings.Join(value, "; ")
				_, err = fmt.Fprintf(w, "    flags_completion+=(%q)\n", handlers)
			***REMOVED*** else ***REMOVED***
				_, err = fmt.Fprintf(w, "    flags_completion+=(:)\n")
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		case BashCompSubdirsInDir:
			_, err := fmt.Fprintf(w, "    flags_with_completion+=(%q)\n", name)

			if len(value) == 1 ***REMOVED***
				ext := "__handle_subdirs_in_dir_flag " + value[0]
				_, err = fmt.Fprintf(w, "    flags_completion+=(%q)\n", ext)
			***REMOVED*** else ***REMOVED***
				ext := "_filedir -d"
				_, err = fmt.Fprintf(w, "    flags_completion+=(%q)\n", ext)
			***REMOVED***
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func writeShortFlag(flag *pflag.Flag, w io.Writer) error ***REMOVED***
	b := (len(flag.NoOptDefVal) > 0)
	name := flag.Shorthand
	format := "    "
	if !b ***REMOVED***
		format += "two_word_"
	***REMOVED***
	format += "flags+=(\"-%s\")\n"
	if _, err := fmt.Fprintf(w, format, name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return writeFlagHandler("-"+name, flag.Annotations, w)
***REMOVED***

func writeFlag(flag *pflag.Flag, w io.Writer) error ***REMOVED***
	b := (len(flag.NoOptDefVal) > 0)
	name := flag.Name
	format := "    flags+=(\"--%s"
	if !b ***REMOVED***
		format += "="
	***REMOVED***
	format += "\")\n"
	if _, err := fmt.Fprintf(w, format, name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return writeFlagHandler("--"+name, flag.Annotations, w)
***REMOVED***

func writeLocalNonPersistentFlag(flag *pflag.Flag, w io.Writer) error ***REMOVED***
	b := (len(flag.NoOptDefVal) > 0)
	name := flag.Name
	format := "    local_nonpersistent_flags+=(\"--%s"
	if !b ***REMOVED***
		format += "="
	***REMOVED***
	format += "\")\n"
	if _, err := fmt.Fprintf(w, format, name); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func writeFlags(cmd *Command, w io.Writer) error ***REMOVED***
	_, err := fmt.Fprintf(w, `    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

`)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	localNonPersistentFlags := cmd.LocalNonPersistentFlags()
	var visitErr error
	cmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) ***REMOVED***
		if err := writeFlag(flag, w); err != nil ***REMOVED***
			visitErr = err
			return
		***REMOVED***
		if len(flag.Shorthand) > 0 ***REMOVED***
			if err := writeShortFlag(flag, w); err != nil ***REMOVED***
				visitErr = err
				return
			***REMOVED***
		***REMOVED***
		if localNonPersistentFlags.Lookup(flag.Name) != nil ***REMOVED***
			if err := writeLocalNonPersistentFlag(flag, w); err != nil ***REMOVED***
				visitErr = err
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if visitErr != nil ***REMOVED***
		return visitErr
	***REMOVED***
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) ***REMOVED***
		if err := writeFlag(flag, w); err != nil ***REMOVED***
			visitErr = err
			return
		***REMOVED***
		if len(flag.Shorthand) > 0 ***REMOVED***
			if err := writeShortFlag(flag, w); err != nil ***REMOVED***
				visitErr = err
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	if visitErr != nil ***REMOVED***
		return visitErr
	***REMOVED***

	_, err = fmt.Fprintf(w, "\n")
	return err
***REMOVED***

func writeRequiredFlag(cmd *Command, w io.Writer) error ***REMOVED***
	if _, err := fmt.Fprintf(w, "    must_have_one_flag=()\n"); err != nil ***REMOVED***
		return err
	***REMOVED***
	flags := cmd.NonInheritedFlags()
	var visitErr error
	flags.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		for key := range flag.Annotations ***REMOVED***
			switch key ***REMOVED***
			case BashCompOneRequiredFlag:
				format := "    must_have_one_flag+=(\"--%s"
				b := (flag.Value.Type() == "bool")
				if !b ***REMOVED***
					format += "="
				***REMOVED***
				format += "\")\n"
				if _, err := fmt.Fprintf(w, format, flag.Name); err != nil ***REMOVED***
					visitErr = err
					return
				***REMOVED***

				if len(flag.Shorthand) > 0 ***REMOVED***
					if _, err := fmt.Fprintf(w, "    must_have_one_flag+=(\"-%s\")\n", flag.Shorthand); err != nil ***REMOVED***
						visitErr = err
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return visitErr
***REMOVED***

func writeRequiredNouns(cmd *Command, w io.Writer) error ***REMOVED***
	if _, err := fmt.Fprintf(w, "    must_have_one_noun=()\n"); err != nil ***REMOVED***
		return err
	***REMOVED***
	sort.Sort(sort.StringSlice(cmd.ValidArgs))
	for _, value := range cmd.ValidArgs ***REMOVED***
		if _, err := fmt.Fprintf(w, "    must_have_one_noun+=(%q)\n", value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func writeArgAliases(cmd *Command, w io.Writer) error ***REMOVED***
	if _, err := fmt.Fprintf(w, "    noun_aliases=()\n"); err != nil ***REMOVED***
		return err
	***REMOVED***
	sort.Sort(sort.StringSlice(cmd.ArgAliases))
	for _, value := range cmd.ArgAliases ***REMOVED***
		if _, err := fmt.Fprintf(w, "    noun_aliases+=(%q)\n", value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func gen(cmd *Command, w io.Writer) error ***REMOVED***
	for _, c := range cmd.Commands() ***REMOVED***
		if !c.IsAvailableCommand() || c == cmd.helpCommand ***REMOVED***
			continue
		***REMOVED***
		if err := gen(c, w); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	commandName := cmd.CommandPath()
	commandName = strings.Replace(commandName, " ", "_", -1)
	commandName = strings.Replace(commandName, ":", "__", -1)
	if _, err := fmt.Fprintf(w, "_%s()\n***REMOVED***\n", commandName); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := fmt.Fprintf(w, "    last_command=%q\n", commandName); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := writeCommands(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := writeFlags(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := writeRequiredFlag(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := writeRequiredNouns(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := writeArgAliases(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := fmt.Fprintf(w, "***REMOVED***\n\n"); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func (cmd *Command) GenBashCompletion(w io.Writer) error ***REMOVED***
	if err := preamble(w, cmd.Name()); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(cmd.BashCompletionFunction) > 0 ***REMOVED***
		if _, err := fmt.Fprintf(w, "%s\n", cmd.BashCompletionFunction); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := gen(cmd, w); err != nil ***REMOVED***
		return err
	***REMOVED***
	return postscript(w, cmd.Name())
***REMOVED***

func (cmd *Command) GenBashCompletionFile(filename string) error ***REMOVED***
	outFile, err := os.Create(filename)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer outFile.Close()

	return cmd.GenBashCompletion(outFile)
***REMOVED***

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag, if it exists.
func (cmd *Command) MarkFlagRequired(name string) error ***REMOVED***
	return MarkFlagRequired(cmd.Flags(), name)
***REMOVED***

// MarkPersistentFlagRequired adds the BashCompOneRequiredFlag annotation to the named persistent flag, if it exists.
func (cmd *Command) MarkPersistentFlagRequired(name string) error ***REMOVED***
	return MarkFlagRequired(cmd.PersistentFlags(), name)
***REMOVED***

// MarkFlagRequired adds the BashCompOneRequiredFlag annotation to the named flag in the flag set, if it exists.
func MarkFlagRequired(flags *pflag.FlagSet, name string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompOneRequiredFlag, []string***REMOVED***"true"***REMOVED***)
***REMOVED***

// MarkFlagFilename adds the BashCompFilenameExt annotation to the named flag, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func (cmd *Command) MarkFlagFilename(name string, extensions ...string) error ***REMOVED***
	return MarkFlagFilename(cmd.Flags(), name, extensions...)
***REMOVED***

// MarkFlagCustom adds the BashCompCustom annotation to the named flag, if it exists.
// Generated bash autocompletion will call the bash function f for the flag.
func (cmd *Command) MarkFlagCustom(name string, f string) error ***REMOVED***
	return MarkFlagCustom(cmd.Flags(), name, f)
***REMOVED***

// MarkPersistentFlagFilename adds the BashCompFilenameExt annotation to the named persistent flag, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func (cmd *Command) MarkPersistentFlagFilename(name string, extensions ...string) error ***REMOVED***
	return MarkFlagFilename(cmd.PersistentFlags(), name, extensions...)
***REMOVED***

// MarkFlagFilename adds the BashCompFilenameExt annotation to the named flag in the flag set, if it exists.
// Generated bash autocompletion will select filenames for the flag, limiting to named extensions if provided.
func MarkFlagFilename(flags *pflag.FlagSet, name string, extensions ...string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompFilenameExt, extensions)
***REMOVED***

// MarkFlagCustom adds the BashCompCustom annotation to the named flag in the flag set, if it exists.
// Generated bash autocompletion will call the bash function f for the flag.
func MarkFlagCustom(flags *pflag.FlagSet, name string, f string) error ***REMOVED***
	return flags.SetAnnotation(name, BashCompCustom, []string***REMOVED***f***REMOVED***)
***REMOVED***
