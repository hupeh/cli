package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
)

// ActionFunc 命令执行函数签名
//
// 参数:
//   - ctx: context.Context，用于传递取消信号和超时控制
//   - cmd: *Command，包含命令的所有信息（标志集、元数据等）
//
// 返回:
//   - error: 执行错误，nil 表示成功
type ActionFunc func(ctx context.Context, cmd *Command) error

// Command 命令结构，代表一个 CLI 命令
//
// Command 封装了命令的元数据（名称、描述）、
// 标志定义和执行逻辑。
type Command struct {
	Name         string        // 命令名称（如 "init", "migrate"）
	Usage        string        // 命令用途简短描述（一行）
	Description  string        // 命令详细描述（多行）
	Flags        *flag.FlagSet // 命令标志集（用于定义和解析命令行参数）
	Action       ActionFunc    // 命令执行函数
	HideHelpFlag bool          // 是否隐藏 -h 帮助标志
	appName      string        // 应用名称（用于打印帮助时显示完整用法）
}

// NewCommand 创建新命令
func NewCommand(name, usage string) *Command {
	return &Command{
		Name:  name,
		Usage: usage,
		Flags: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

// DefaultHelpCommand 创建默认的 help 命令
func DefaultHelpCommand() *Command {
	return &Command{
		Name:        "help",
		Usage:       "Show help information",
		Description: "Display help information for commands",
		Flags:       flag.NewFlagSet("help", flag.ContinueOnError),
	}
}

// DefaultVersionCommand 创建默认的 version 命令
func DefaultVersionCommand() *Command {
	return &Command{
		Name:        "version",
		Usage:       "Show version information",
		Description: "Display the version of this program",
		Flags:       flag.NewFlagSet("version", flag.ContinueOnError),
	}
}

// SetOutput 设置输出目标
func (c *Command) SetOutput(w io.Writer) {
	c.Flags.SetOutput(w)
}

// Output 获取输出目标
func (c *Command) Output() io.Writer {
	return c.Flags.Output()
}

// SetAppName 设置应用名称（用于打印帮助）
func (c *Command) SetAppName(name string) {
	c.appName = name
}

// PrintUsage 打印命令使用帮助到默认输出
func (c *Command) PrintUsage() {
	c.PrintUsageTo(c.Output())
}

// PrintUsageTo 打印命令使用帮助到指定的 Writer
func (c *Command) PrintUsageTo(w io.Writer) {
	// 如果有应用名称，显示完整用法
	if c.appName != "" {
		_, _ = fmt.Fprintf(w, "Usage: %s %s [options]\n\n", c.appName, c.Name)
	} else {
		_, _ = fmt.Fprintf(w, "Usage: %s [options]\n\n", c.Name)
	}
	_, _ = fmt.Fprintf(w, "%s\n", c.Usage)

	if c.Description != "" {
		_, _ = fmt.Fprintf(w, "\n%s\n", c.Description)
	}

	// 检查是否有标志
	hasFlags := false
	c.Flags.VisitAll(func(f *flag.Flag) {
		hasFlags = true
	})

	if hasFlags {
		_, _ = fmt.Fprintln(w, "\nOptions:")
		// 临时设置 FlagSet 的输出以便 PrintDefaults 输出到指定的 Writer
		oldOutput := c.Flags.Output()
		c.Flags.SetOutput(w)
		c.Flags.PrintDefaults()
		c.Flags.SetOutput(oldOutput)
	}
}

// Run 执行命令（使用 context.Background()）
func (c *Command) Run(args []string) error {
	return c.RunContext(context.Background(), args)
}

// RunContext 使用指定的 context 执行命令
func (c *Command) RunContext(ctx context.Context, args []string) error {
	// 设置 Usage 函数
	if c.HideHelpFlag {
		// 隐藏帮助时，设置一个空函数来阻止默认 usage 输出
		c.Flags.Usage = func() {}
	} else {
		// 显示帮助时，使用自定义的 PrintUsage
		c.Flags.Usage = c.PrintUsage
	}

	// 解析参数
	if err := c.Flags.Parse(args); err != nil {
		// 如果隐藏帮助，直接返回错误（包括 ErrHelp）
		if c.HideHelpFlag {
			return err
		}
		// 显示帮助时，ErrHelp 不算错误（因为已经打印了帮助信息）
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	// 执行命令
	if c.Action != nil {
		return c.Action(ctx, c)
	}

	return nil
}
