package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Program CLI 应用程序
type Program struct {
	Commands           []*Command // 命令列表
	Name               string     // 应用名称
	Usage              string     // 应用描述
	Version            string     // 应用版本
	Banner             string     // 应用横幅（ASCII 艺术字等）
	DefaultCommand     string     // 默认命令名称（当未指定命令时使用）
	HideHelpCommand    bool       // 隐藏 help 命令
	HideVersionCommand bool       // 隐藏 version 命令
	HideHelpFlag       bool       // 隐藏 -h/--help 标志
	HideVersionFlag    bool       // 隐藏 -v/--version 标志
	HelpCommand        *Command   // help 命令（可自定义）
	VersionCommand     *Command   // version 命令（可自定义）
	output             io.Writer  // 输出目标（测试时可替换，默认 os.Stderr）
}

// NewProgram 创建 CLI 应用程序
func NewProgram(appName, version string) *Program {
	return &Program{
		Name:    appName,
		Version: version,
	}
}

// SetOutput 设置输出目标
func (p *Program) SetOutput(w io.Writer) {
	p.output = w
}

// Output 获取输出目标，如果未设置则返回 os.Stderr
func (p *Program) Output() io.Writer {
	if p.output == nil {
		return os.Stderr
	}
	return p.output
}

// Get 获取命令并配置其输出和应用名称
//
// 从已注册的命令和内置命令（help、version）中查找指定名称的命令。
// 找到命令后会自动设置命令的输出目标和应用名称。
func (p *Program) Get(name string) *Command {
	if cmd := p.get(name); cmd != nil {
		cmd.SetOutput(p.Output())
		cmd.SetAppName(p.Name)
		return cmd
	}

	return nil
}

func (p *Program) get(name string) *Command {
	// 首先查找用户注册的命令
	for _, cmd := range p.Commands {
		if cmd.Name == name {
			return cmd
		}
	}

	// 查找内置命令
	if !p.HideHelpCommand && name == "help" {
		if p.HelpCommand != nil {
			// 用户自定义的 help 命令
			return p.HelpCommand
		}
		// 临时创建默认命令
		return DefaultHelpCommand()
	}

	if !p.HideVersionCommand && name == "version" {
		if p.VersionCommand != nil {
			// 用户自定义的 version 命令
			return p.VersionCommand
		}
		// 临时创建默认命令
		return DefaultVersionCommand()
	}

	return nil
}

// PrintUsage 打印总体使用帮助到指定输出
func (p *Program) PrintUsage() error {
	w := p.Output()
	var b []byte

	// 如果有横幅，先打印横幅
	if p.Banner != "" {
		b = fmt.Appendln(b, p.Banner)
		b = fmt.Appendln(b) // 空行分隔
	}

	b = fmt.Appendf(b, "%s version %s\n", p.Name, p.Version)

	// 如果有应用描述，打印它
	if p.Usage != "" {
		b = fmt.Appendf(b, "%s\n", p.Usage)
	}

	b = fmt.Appendf(b, "\nUSAGE:\n")
	b = fmt.Appendf(b, "    %s [command] [options]\n\n", p.Name)
	b = fmt.Appendf(b, "COMMANDS:\n")

	// 计算最长命令名长度，用于对齐
	maxLen := 0
	for _, cmd := range p.Commands {
		if len(cmd.Name) > maxLen {
			maxLen = len(cmd.Name)
		}
	}

	// 按注册顺序打印命令
	for _, cmd := range p.Commands {
		b = fmt.Appendf(b, "    %-*s    %s\n", maxLen, cmd.Name, cmd.Usage)
	}

	b = fmt.Appendf(b, "\nRun '%s [command] -h' for more information on a command.\n", p.Name)

	// 一次性写入到 w
	_, err := w.Write(b)
	return err
}

// Run 运行命令（使用 context.Background()）
func (p *Program) Run(args []string) error {
	return p.RunContext(context.Background(), args)
}

// isFlag 判断参数是否是 flag（以 - 或 -- 开头）
func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

// RunContext 使用指定的 context 运行命令
func (p *Program) RunContext(ctx context.Context, args []string) error {
	// 解析命令名称和参数起始位置
	var cmdName string
	var cmdArgs []string
	var usingDefaultCommand bool

	if len(args) < 2 || isFlag(args[1]) {
		// 没有提供命令，或第一个参数是 flag
		if p.DefaultCommand != "" {
			cmdName = p.DefaultCommand
			usingDefaultCommand = true
			if len(args) >= 2 {
				cmdArgs = args[1:] // 将 flag 传递给默认命令
			} else {
				cmdArgs = []string{}
			}
		} else {
			// 没有默认命令，显示帮助
			return p.PrintUsage()
		}
	} else {
		// 显式指定了命令
		cmdName = args[1]
		cmdArgs = args[2:]
	}

	// 处理全局 flag（检查 cmdArgs 中是否包含全局 flag）
	for _, arg := range cmdArgs {
		if !p.HideVersionFlag && (arg == "-v" || arg == "--version") {
			if _, err := fmt.Fprintf(p.Output(), "%s version %s\n", p.Name, p.Version); err != nil {
				return err
			}
			return nil
		}
		if !p.HideHelpFlag && (arg == "-h" || arg == "--help") {
			return p.PrintUsage()
		}
	}

	// 处理特殊命令
	// 1. 处理 version 命令
	if !p.HideVersionCommand && cmdName == "version" {
		if _, err := fmt.Fprintf(p.Output(), "%s version %s\n", p.Name, p.Version); err != nil {
			return err
		}
		return nil
	}

	// 2. 处理 help 命令：help [command]
	if !p.HideHelpCommand && cmdName == "help" {
		if len(cmdArgs) > 0 {
			// help [command] - 显示特定命令的帮助
			subCmdName := cmdArgs[0]
			cmd := p.Get(subCmdName)
			if cmd == nil {
				if _, err := fmt.Fprintf(p.Output(), "help: unknown command: %s\n", subCmdName); err != nil {
					return err
				}
				return fmt.Errorf("unknown command: %s", subCmdName)
			}
			return cmd.PrintUsage()
		}
		// help - 显示总体帮助
		return p.PrintUsage()
	}

	// 查找并执行命令
	cmd := p.Get(cmdName)
	if cmd == nil {
		if usingDefaultCommand {
			if _, err := fmt.Fprintf(p.Output(), "Default command '%s' not found\n\n", cmdName); err != nil {
				return err
			}
			if err := p.PrintUsage(); err != nil {
				return err
			}
			return fmt.Errorf("default command not found: %s", cmdName)
		}
		if _, err := fmt.Fprintf(p.Output(), "Unknown command: %s\n\n", cmdName); err != nil {
			return err
		}
		if err := p.PrintUsage(); err != nil {
			return err
		}
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	return cmd.RunContext(ctx, cmdArgs)
}
