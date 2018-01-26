package instructions

import (
	"strings"
	"testing"

	"github.com/docker/docker/builder/dockerfile/command"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandsExactlyOneArgument(t *testing.T) ***REMOVED***
	commands := []string***REMOVED***
		"MAINTAINER",
		"WORKDIR",
		"USER",
		"STOPSIGNAL",
	***REMOVED***

	for _, command := range commands ***REMOVED***
		ast, err := parser.Parse(strings.NewReader(command))
		require.NoError(t, err)
		_, err = ParseInstruction(ast.AST.Children[0])
		assert.EqualError(t, err, errExactlyOneArgument(command).Error())
	***REMOVED***
***REMOVED***

func TestCommandsAtLeastOneArgument(t *testing.T) ***REMOVED***
	commands := []string***REMOVED***
		"ENV",
		"LABEL",
		"ONBUILD",
		"HEALTHCHECK",
		"EXPOSE",
		"VOLUME",
	***REMOVED***

	for _, command := range commands ***REMOVED***
		ast, err := parser.Parse(strings.NewReader(command))
		require.NoError(t, err)
		_, err = ParseInstruction(ast.AST.Children[0])
		assert.EqualError(t, err, errAtLeastOneArgument(command).Error())
	***REMOVED***
***REMOVED***

func TestCommandsNoDestinationArgument(t *testing.T) ***REMOVED***
	commands := []string***REMOVED***
		"ADD",
		"COPY",
	***REMOVED***

	for _, command := range commands ***REMOVED***
		ast, err := parser.Parse(strings.NewReader(command + " arg1"))
		require.NoError(t, err)
		_, err = ParseInstruction(ast.AST.Children[0])
		assert.EqualError(t, err, errNoDestinationArgument(command).Error())
	***REMOVED***
***REMOVED***

func TestCommandsTooManyArguments(t *testing.T) ***REMOVED***
	commands := []string***REMOVED***
		"ENV",
		"LABEL",
	***REMOVED***

	for _, command := range commands ***REMOVED***
		node := &parser.Node***REMOVED***
			Original: command + "arg1 arg2 arg3",
			Value:    strings.ToLower(command),
			Next: &parser.Node***REMOVED***
				Value: "arg1",
				Next: &parser.Node***REMOVED***
					Value: "arg2",
					Next: &parser.Node***REMOVED***
						Value: "arg3",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		_, err := ParseInstruction(node)
		assert.EqualError(t, err, errTooManyArguments(command).Error())
	***REMOVED***
***REMOVED***

func TestCommandsBlankNames(t *testing.T) ***REMOVED***
	commands := []string***REMOVED***
		"ENV",
		"LABEL",
	***REMOVED***

	for _, command := range commands ***REMOVED***
		node := &parser.Node***REMOVED***
			Original: command + " =arg2",
			Value:    strings.ToLower(command),
			Next: &parser.Node***REMOVED***
				Value: "",
				Next: &parser.Node***REMOVED***
					Value: "arg2",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		_, err := ParseInstruction(node)
		assert.EqualError(t, err, errBlankCommandNames(command).Error())
	***REMOVED***
***REMOVED***

func TestHealthCheckCmd(t *testing.T) ***REMOVED***
	node := &parser.Node***REMOVED***
		Value: command.Healthcheck,
		Next: &parser.Node***REMOVED***
			Value: "CMD",
			Next: &parser.Node***REMOVED***
				Value: "hello",
				Next: &parser.Node***REMOVED***
					Value: "world",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	cmd, err := ParseInstruction(node)
	assert.NoError(t, err)
	hc, ok := cmd.(*HealthCheckCommand)
	assert.True(t, ok)
	expected := []string***REMOVED***"CMD-SHELL", "hello world"***REMOVED***
	assert.Equal(t, expected, hc.Health.Test)
***REMOVED***

func TestParseOptInterval(t *testing.T) ***REMOVED***
	flInterval := &Flag***REMOVED***
		name:     "interval",
		flagType: stringType,
		Value:    "50ns",
	***REMOVED***
	_, err := parseOptInterval(flInterval)
	testutil.ErrorContains(t, err, "cannot be less than 1ms")

	flInterval.Value = "1ms"
	_, err = parseOptInterval(flInterval)
	require.NoError(t, err)
***REMOVED***

func TestErrorCases(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		name          string
		dockerfile    string
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***
			name: "copyEmptyWhitespace",
			dockerfile: `COPY	
		quux \
      bar`,
			expectedError: "COPY requires at least two arguments",
		***REMOVED***,
		***REMOVED***
			name:          "ONBUILD forbidden FROM",
			dockerfile:    "ONBUILD FROM scratch",
			expectedError: "FROM isn't allowed as an ONBUILD trigger",
		***REMOVED***,
		***REMOVED***
			name:          "ONBUILD forbidden MAINTAINER",
			dockerfile:    "ONBUILD MAINTAINER docker.io",
			expectedError: "MAINTAINER isn't allowed as an ONBUILD trigger",
		***REMOVED***,
		***REMOVED***
			name:          "ARG two arguments",
			dockerfile:    "ARG foo bar",
			expectedError: "ARG requires exactly one argument",
		***REMOVED***,
		***REMOVED***
			name:          "MAINTAINER unknown flag",
			dockerfile:    "MAINTAINER --boo joe@example.com",
			expectedError: "Unknown flag: boo",
		***REMOVED***,
		***REMOVED***
			name:          "Chaining ONBUILD",
			dockerfile:    `ONBUILD ONBUILD RUN touch foobar`,
			expectedError: "Chaining ONBUILD via `ONBUILD ONBUILD` isn't allowed",
		***REMOVED***,
		***REMOVED***
			name:          "Invalid instruction",
			dockerfile:    `foo bar`,
			expectedError: "unknown instruction: FOO",
		***REMOVED***,
	***REMOVED***
	for _, c := range cases ***REMOVED***
		r := strings.NewReader(c.dockerfile)
		ast, err := parser.Parse(r)

		if err != nil ***REMOVED***
			t.Fatalf("Error when parsing Dockerfile: %s", err)
		***REMOVED***
		n := ast.AST.Children[0]
		_, err = ParseInstruction(n)
		testutil.ErrorContains(t, err, c.expectedError)
	***REMOVED***

***REMOVED***
