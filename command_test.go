package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	if cmd.Name != "test" {
		t.Errorf("Expected Name to be 'test', got '%s'", cmd.Name)
	}
	if cmd.Usage != "Test command" {
		t.Errorf("Expected Usage to be 'Test command', got '%s'", cmd.Usage)
	}
	if cmd.Flags == nil {
		t.Error("Expected Flags to be initialized")
	}
}

func TestDefaultHelpCommand(t *testing.T) {
	cmd := DefaultHelpCommand()
	if cmd.Name != "help" {
		t.Errorf("Expected Name to be 'help', got '%s'", cmd.Name)
	}
	if cmd.Usage != "Show help information" {
		t.Errorf("Expected Usage to be 'Show help information', got '%s'", cmd.Usage)
	}
	if cmd.Flags == nil {
		t.Error("Expected Flags to be initialized")
	}
}

func TestDefaultVersionCommand(t *testing.T) {
	cmd := DefaultVersionCommand()
	if cmd.Name != "version" {
		t.Errorf("Expected Name to be 'version', got '%s'", cmd.Name)
	}
	if cmd.Usage != "Show version information" {
		t.Errorf("Expected Usage to be 'Show version information', got '%s'", cmd.Usage)
	}
	if cmd.Flags == nil {
		t.Error("Expected Flags to be initialized")
	}
}

func TestCommand_SetOutput(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)

	if cmd.Output() != buf {
		t.Error("Expected Output() to return the same buffer set by SetOutput()")
	}
}

func TestCommand_SetAppName(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.SetAppName("myapp")

	if cmd.appName != "myapp" {
		t.Errorf("Expected appName to be 'myapp', got '%s'", cmd.appName)
	}
}

func TestCommand_PrintUsage(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.Description = "This is a detailed description"
	cmd.Flags.String("config", "", "Config file path")

	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.PrintUsage()

	output := buf.String()
	if !strings.Contains(output, "test") {
		t.Error("Expected output to contain command name")
	}
	if !strings.Contains(output, "Test command") {
		t.Error("Expected output to contain usage")
	}
	if !strings.Contains(output, "This is a detailed description") {
		t.Error("Expected output to contain description")
	}
	if !strings.Contains(output, "-config") {
		t.Error("Expected output to contain flag")
	}
}

func TestCommand_PrintUsageWithAppName(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.SetAppName("myapp")

	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.PrintUsage()

	output := buf.String()
	if !strings.Contains(output, "myapp test") {
		t.Error("Expected output to contain app name and command name")
	}
}

func TestCommand_PrintUsageTo(t *testing.T) {
	cmd := NewCommand("test", "Test command")

	buf := &bytes.Buffer{}
	cmd.PrintUsageTo(buf)

	output := buf.String()
	if !strings.Contains(output, "test") {
		t.Error("Expected output to contain command name")
	}
}

func TestCommand_Run(t *testing.T) {
	executed := false
	cmd := NewCommand("test", "Test command")
	cmd.Action = func(ctx context.Context, c *Command) error {
		executed = true
		return nil
	}

	err := cmd.Run([]string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Error("Expected action to be executed")
	}
}

func TestCommand_RunContext(t *testing.T) {
	ctx := context.Background()
	executed := false
	var receivedCtx context.Context

	cmd := NewCommand("test", "Test command")
	cmd.Action = func(c context.Context, co *Command) error {
		executed = true
		receivedCtx = c
		return nil
	}

	err := cmd.RunContext(ctx, []string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Error("Expected action to be executed")
	}
	if receivedCtx != ctx {
		t.Error("Expected action to receive the same context")
	}
}

func TestCommand_RunWithError(t *testing.T) {
	expectedErr := errors.New("test error")
	cmd := NewCommand("test", "Test command")
	cmd.Action = func(ctx context.Context, c *Command) error {
		return expectedErr
	}

	err := cmd.Run([]string{})
	if err != expectedErr {
		t.Errorf("Expected error to be %v, got %v", expectedErr, err)
	}
}

func TestCommand_RunWithFlags(t *testing.T) {
	var configValue string
	cmd := NewCommand("test", "Test command")
	cmd.Flags.StringVar(&configValue, "config", "", "Config file path")
	cmd.Action = func(ctx context.Context, c *Command) error {
		return nil
	}

	err := cmd.Run([]string{"-config", "test.conf"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if configValue != "test.conf" {
		t.Errorf("Expected config to be 'test.conf', got '%s'", configValue)
	}
}

func TestCommand_RunWithHelp(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.Action = func(ctx context.Context, c *Command) error {
		return errors.New("should not be called")
	}

	err := cmd.Run([]string{"-h"})
	if err != nil {
		t.Errorf("Expected no error for help flag, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("Expected help output to contain 'Usage:'")
	}
}

func TestCommand_RunWithHelpHidden(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.HideHelpFlag = true
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)

	err := cmd.Run([]string{"-h"})
	if err != flag.ErrHelp {
		t.Errorf("Expected ErrHelp, got %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("Expected no output when help is hidden, got: %s", output)
	}
}

func TestCommand_RunWithInvalidFlag(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)

	err := cmd.Run([]string{"-invalid"})
	if err == nil {
		t.Error("Expected error for invalid flag")
	}
}

func TestCommand_RunWithoutAction(t *testing.T) {
	cmd := NewCommand("test", "Test command")

	err := cmd.Run([]string{})
	if err != nil {
		t.Errorf("Expected no error when action is nil, got %v", err)
	}
}

func TestCommand_FlagsAreParsed(t *testing.T) {
	var verbose bool
	var count int

	cmd := NewCommand("test", "Test command")
	cmd.Flags.BoolVar(&verbose, "verbose", false, "Verbose output")
	cmd.Flags.IntVar(&count, "count", 0, "Count value")
	cmd.Action = func(ctx context.Context, c *Command) error {
		if !verbose {
			return errors.New("expected verbose to be true")
		}
		if count != 5 {
			return errors.New("expected count to be 5")
		}
		return nil
	}

	err := cmd.Run([]string{"-verbose", "-count", "5"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCommand_ActionReceivesCommand(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	var receivedCmd *Command

	cmd.Action = func(ctx context.Context, c *Command) error {
		receivedCmd = c
		return nil
	}

	err := cmd.Run([]string{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if receivedCmd != cmd {
		t.Error("Expected action to receive the same command instance")
	}
}
