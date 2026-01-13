package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNewProgram(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	if prog.Name != "testapp" {
		t.Errorf("Expected Name to be 'testapp', got '%s'", prog.Name)
	}
	if prog.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got '%s'", prog.Version)
	}
}

func TestProgram_SetOutput(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	if prog.Output() != buf {
		t.Error("Expected Output() to return the same buffer set by SetOutput()")
	}
}

func TestProgram_OutputDefault(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	// 未设置 output 时应返回 os.Stderr
	output := prog.Output()
	if output == nil {
		t.Error("Expected Output() to return a non-nil writer")
	}
}

func TestProgram_Get(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	testCmd := NewCommand("test", "Test command")
	prog.Commands = []*Command{testCmd}

	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	cmd := prog.Get("test")
	if cmd == nil {
		t.Fatal("Expected to get test command")
	}
	if cmd.Name != "test" {
		t.Errorf("Expected command name to be 'test', got '%s'", cmd.Name)
	}
	if cmd.appName != "testapp" {
		t.Errorf("Expected appName to be 'testapp', got '%s'", cmd.appName)
	}
	if cmd.Output() != buf {
		t.Error("Expected command output to be set to program output")
	}
}

func TestProgram_GetNonExistent(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	cmd := prog.Get("nonexistent")
	if cmd != nil {
		t.Error("Expected nil for non-existent command")
	}
}

func TestProgram_GetHelpCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	cmd := prog.Get("help")
	if cmd == nil {
		t.Fatal("Expected to get help command")
	}
	if cmd.Name != "help" {
		t.Errorf("Expected command name to be 'help', got '%s'", cmd.Name)
	}
}

func TestProgram_GetHelpCommandHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideHelpCommand = true

	cmd := prog.Get("help")
	if cmd != nil {
		t.Error("Expected nil when help command is hidden")
	}
}

func TestProgram_GetCustomHelpCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	customHelp := NewCommand("help", "Custom help")
	prog.HelpCommand = customHelp

	cmd := prog.Get("help")
	if cmd == nil {
		t.Fatal("Expected to get help command")
	}
	if cmd.Usage != "Custom help" {
		t.Errorf("Expected custom help usage, got '%s'", cmd.Usage)
	}
}

func TestProgram_GetVersionCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	cmd := prog.Get("version")
	if cmd == nil {
		t.Fatal("Expected to get version command")
	}
	if cmd.Name != "version" {
		t.Errorf("Expected command name to be 'version', got '%s'", cmd.Name)
	}
}

func TestProgram_GetVersionCommandHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideVersionCommand = true

	cmd := prog.Get("version")
	if cmd != nil {
		t.Error("Expected nil when version command is hidden")
	}
}

func TestProgram_GetCustomVersionCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	customVersion := NewCommand("version", "Custom version")
	prog.VersionCommand = customVersion

	cmd := prog.Get("version")
	if cmd == nil {
		t.Fatal("Expected to get version command")
	}
	if cmd.Usage != "Custom version" {
		t.Errorf("Expected custom version usage, got '%s'", cmd.Usage)
	}
}

func TestProgram_PrintUsage(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.Usage = "A test application"
	prog.Commands = []*Command{
		NewCommand("init", "Initialize project"),
		NewCommand("build", "Build project"),
	}

	buf := &bytes.Buffer{}
	prog.SetOutput(buf)
	prog.PrintUsage()

	output := buf.String()
	if !strings.Contains(output, "testapp") {
		t.Error("Expected output to contain app name")
	}
	if !strings.Contains(output, "1.0.0") {
		t.Error("Expected output to contain version")
	}
	if !strings.Contains(output, "A test application") {
		t.Error("Expected output to contain usage")
	}
	if !strings.Contains(output, "init") {
		t.Error("Expected output to contain init command")
	}
	if !strings.Contains(output, "build") {
		t.Error("Expected output to contain build command")
	}
}

func TestProgram_PrintUsageWithBanner(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.Banner = "TEST APP BANNER"

	buf := &bytes.Buffer{}
	prog.SetOutput(buf)
	prog.PrintUsage()

	output := buf.String()
	if !strings.Contains(output, "TEST APP BANNER") {
		t.Error("Expected output to contain banner")
	}
}

func TestProgram_PrintUsageTo(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	buf := &bytes.Buffer{}
	prog.PrintUsageTo(buf)

	output := buf.String()
	if !strings.Contains(output, "testapp") {
		t.Error("Expected output to contain app name")
	}
}

func TestProgram_RunNoArgs(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "USAGE:") {
		t.Error("Expected usage to be printed when no command is provided")
	}
}

func TestProgram_RunVersionCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "version"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "testapp version 1.0.0") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestProgram_RunVersionFlag(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	tests := []string{"-v", "--version"}
	for _, flag := range tests {
		buf.Reset()
		err := prog.Run([]string{"testapp", flag})
		if err != nil {
			t.Errorf("Expected no error for %s, got %v", flag, err)
		}

		output := buf.String()
		if !strings.Contains(output, "testapp version 1.0.0") {
			t.Errorf("Expected version output for %s, got: %s", flag, output)
		}
	}
}

func TestProgram_RunVersionFlagHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideVersionFlag = true
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "-v"})
	if err == nil {
		t.Error("Expected error when version flag is hidden")
	}
}

func TestProgram_RunVersionCommandHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideVersionCommand = true
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "version"})
	if err == nil {
		t.Error("Expected error when version command is hidden")
	}
}

func TestProgram_RunHelpCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "help"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "USAGE:") {
		t.Error("Expected help output to contain USAGE")
	}
}

func TestProgram_RunHelpCommandWithSubcommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	testCmd := NewCommand("test", "Test command")
	testCmd.Description = "Detailed test description"
	prog.Commands = []*Command{testCmd}
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "help", "test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test command") {
		t.Error("Expected help output to contain test command usage")
	}
	if !strings.Contains(output, "Detailed test description") {
		t.Error("Expected help output to contain test command description")
	}
}

func TestProgram_RunHelpCommandWithInvalidSubcommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "help", "invalid"})
	if err == nil {
		t.Error("Expected error for invalid subcommand")
	}
}

func TestProgram_RunHelpFlag(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	tests := []string{"-h", "--help"}
	for _, flag := range tests {
		buf.Reset()
		err := prog.Run([]string{"testapp", flag})
		if err != nil {
			t.Errorf("Expected no error for %s, got %v", flag, err)
		}

		output := buf.String()
		if !strings.Contains(output, "USAGE:") {
			t.Errorf("Expected help output for %s", flag)
		}
	}
}

func TestProgram_RunHelpFlagHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideHelpFlag = true
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "-h"})
	if err == nil {
		t.Error("Expected error when help flag is hidden")
	}
}

func TestProgram_RunHelpCommandHidden(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	prog.HideHelpCommand = true
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "help"})
	if err == nil {
		t.Error("Expected error when help command is hidden")
	}
}

func TestProgram_RunCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	executed := false
	testCmd := NewCommand("test", "Test command")
	testCmd.Action = func(ctx context.Context, cmd *Command) error {
		executed = true
		return nil
	}
	prog.Commands = []*Command{testCmd}

	err := prog.Run([]string{"testapp", "test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Error("Expected command to be executed")
	}
}

func TestProgram_RunCommandWithArgs(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	var configValue string
	testCmd := NewCommand("test", "Test command")
	testCmd.Flags.StringVar(&configValue, "config", "", "Config file")
	testCmd.Action = func(ctx context.Context, cmd *Command) error {
		return nil
	}
	prog.Commands = []*Command{testCmd}

	err := prog.Run([]string{"testapp", "test", "-config", "test.conf"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if configValue != "test.conf" {
		t.Errorf("Expected config to be 'test.conf', got '%s'", configValue)
	}
}

func TestProgram_RunUnknownCommand(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	buf := &bytes.Buffer{}
	prog.SetOutput(buf)

	err := prog.Run([]string{"testapp", "unknown"})
	if err == nil {
		t.Error("Expected error for unknown command")
	}

	output := buf.String()
	if !strings.Contains(output, "Unknown command") {
		t.Error("Expected output to contain 'Unknown command'")
	}
}

func TestProgram_RunContext(t *testing.T) {
	ctx := context.Background()
	prog := NewProgram("testapp", "1.0.0")
	executed := false
	var receivedCtx context.Context

	testCmd := NewCommand("test", "Test command")
	testCmd.Action = func(c context.Context, cmd *Command) error {
		executed = true
		receivedCtx = c
		return nil
	}
	prog.Commands = []*Command{testCmd}

	err := prog.RunContext(ctx, []string{"testapp", "test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Error("Expected command to be executed")
	}
	if receivedCtx != ctx {
		t.Error("Expected command to receive the same context")
	}
}

func TestProgram_RunCommandWithError(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")
	expectedErr := errors.New("test error")

	testCmd := NewCommand("test", "Test command")
	testCmd.Action = func(ctx context.Context, cmd *Command) error {
		return expectedErr
	}
	prog.Commands = []*Command{testCmd}

	err := prog.Run([]string{"testapp", "test"})
	if err != expectedErr {
		t.Errorf("Expected error to be %v, got %v", expectedErr, err)
	}
}

func TestProgram_MultipleCommands(t *testing.T) {
	prog := NewProgram("testapp", "1.0.0")

	initExecuted := false
	initCmd := NewCommand("init", "Initialize")
	initCmd.Action = func(ctx context.Context, cmd *Command) error {
		initExecuted = true
		return nil
	}

	buildExecuted := false
	buildCmd := NewCommand("build", "Build")
	buildCmd.Action = func(ctx context.Context, cmd *Command) error {
		buildExecuted = true
		return nil
	}

	prog.Commands = []*Command{initCmd, buildCmd}

	// 测试第一个命令
	err := prog.Run([]string{"testapp", "init"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !initExecuted {
		t.Error("Expected init command to be executed")
	}
	if buildExecuted {
		t.Error("Did not expect build command to be executed")
	}

	// 测试第二个命令
	initExecuted = false
	err = prog.Run([]string{"testapp", "build"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if initExecuted {
		t.Error("Did not expect init command to be executed")
	}
	if !buildExecuted {
		t.Error("Expected build command to be executed")
	}
}

func TestProgram_CommandPriority(t *testing.T) {
	// 用户注册的命令应该优先于内置命令
	prog := NewProgram("testapp", "1.0.0")
	customHelp := NewCommand("help", "Custom help command")
	customHelp.Usage = "My custom help"
	prog.Commands = []*Command{customHelp}

	cmd := prog.Get("help")
	if cmd == nil {
		t.Fatal("Expected to get help command")
	}
	if cmd.Usage != "My custom help" {
		t.Error("Expected user-registered command to take priority over built-in command")
	}
}
