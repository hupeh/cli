# CLI

[![Go Reference](https://pkg.go.dev/badge/github.com/hupeh/cli.svg)](https://pkg.go.dev/github.com/hupeh/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/hupeh/cli)](https://goreportcard.com/report/github.com/hupeh/cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个简洁、轻量级的 Go CLI 应用框架，基于标准库 `flag` 包构建。

## 特性

- 🚀 **简单易用** - 基于标准库 `flag` 包，学习成本低
- 📦 **轻量级** - 无第三方依赖，代码简洁
- 🎯 **灵活可配置** - 支持自定义帮助和版本命令
- 🧩 **子命令支持** - 内置子命令路由和管理
- 🔧 **Context 支持** - 原生支持 `context.Context`，便于超时控制和取消操作
- 📝 **自动帮助生成** - 自动生成格式化的帮助信息

## 安装

```bash
go get github.com/hupeh/cli
```

## 快速开始

### 基础示例

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hupeh/cli"
)

func main() {
	// 创建 CLI 应用
	app := cli.NewProgram("myapp", "1.0.0")
	app.Usage = "A simple CLI application"

	// 创建命令
	initCmd := cli.NewCommand("init", "Initialize a new project")
	initCmd.Flags.String("path", ".", "Project path")
	initCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		path := cmd.Flags.Lookup("path").Value.String()
		fmt.Printf("Initializing project at %s\n", path)
		return nil
	}

	// 注册命令
	app.Commands = []*cli.Command{initCmd}

	// 运行应用
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

运行：

```bash
$ myapp init --path ./myproject
Initializing project at ./myproject
```

### 完整示例

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hupeh/cli"
)

func main() {
	app := cli.NewProgram("myapp", "1.0.0")
	app.Usage = "A feature-rich CLI application"
	app.Banner = `
 __  __                            
|  \/  |_   _  __ _ _ __  _ __    
| |\/| | | | |/ _' | '_ \| '_ \   
| |  | | |_| | (_| | |_) | |_) |  
|_|  |_|\__, |\__,_| .__/| .__/   
        |___/      |_|   |_|      
`

	// init 命令
	initCmd := cli.NewCommand("init", "Initialize a new project")
	initCmd.Description = "Create a new project with the specified configuration"
	var initPath string
	var verbose bool
	initCmd.Flags.StringVar(&initPath, "path", ".", "Project path")
	initCmd.Flags.BoolVar(&verbose, "verbose", false, "Verbose output")
	initCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Printf("Initializing project at %s (verbose mode)\n", initPath)
		} else {
			fmt.Printf("Initializing project at %s\n", initPath)
		}
		return nil
	}

	// build 命令
	buildCmd := cli.NewCommand("build", "Build the project")
	buildCmd.Description = "Compile the project with specified options"
	var buildOutput string
	var optimize bool
	buildCmd.Flags.StringVar(&buildOutput, "output", "bin/app", "Output path")
	buildCmd.Flags.BoolVar(&optimize, "optimize", false, "Enable optimization")
	buildCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		fmt.Printf("Building project to %s (optimize=%v)\n", buildOutput, optimize)
		return nil
	}

	// deploy 命令
	deployCmd := cli.NewCommand("deploy", "Deploy the application")
	var env string
	deployCmd.Flags.StringVar(&env, "env", "production", "Environment (development/staging/production)")
	deployCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		fmt.Printf("Deploying to %s environment\n", env)
		return nil
	}

	app.Commands = []*cli.Command{initCmd, buildCmd, deployCmd}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

## 使用指南

### 创建应用

```go
app := cli.NewProgram("myapp", "1.0.0")
app.Usage = "Application description"
app.Banner = "ASCII art banner (optional)"
```

### 创建命令

```go
cmd := cli.NewCommand("commandname", "Short description")
cmd.Description = "Long description (optional)"
```

### 添加标志

使用标准库 `flag` 包的方式添加标志：

```go
var name string
var age int
var verbose bool

cmd.Flags.StringVar(&name, "name", "default", "User name")
cmd.Flags.IntVar(&age, "age", 0, "User age")
cmd.Flags.BoolVar(&verbose, "verbose", false, "Enable verbose output")
```

### 定义命令行为

```go
cmd.Action = func(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Name: %s, Age: %d\n", name, age)
	return nil
}
```

### 注册命令

```go
app.Commands = []*cli.Command{cmd1, cmd2, cmd3}
```

### 运行应用

```go
if err := app.Run(os.Args); err != nil {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
```

### Context 支持

使用 `RunContext` 支持超时和取消：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := app.RunContext(ctx, os.Args); err != nil {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
```

在命令中检查 context：

```go
cmd.Action = func(ctx context.Context, cmd *cli.Command) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// 执行命令逻辑
	}
	return nil
}
```

## 内置功能

### 自动帮助

框架自动提供帮助功能：

```bash
# 显示应用帮助
$ myapp -h
$ myapp --help
$ myapp help

# 显示命令帮助
$ myapp init -h
$ myapp help init
```

### 自动版本

框架自动提供版本信息：

```bash
$ myapp -v
$ myapp --version
$ myapp version
```

### 自定义帮助和版本

```go
// 隐藏内置的帮助/版本
app.HideHelpCommand = true
app.HideVersionCommand = true
app.HideHelpFlag = true
app.HideVersionFlag = true

// 自定义帮助命令
customHelp := cli.NewCommand("help", "Custom help")
customHelp.Action = func(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("My custom help message")
	return nil
}
app.HelpCommand = customHelp

// 自定义版本命令
customVersion := cli.NewCommand("version", "Custom version")
customVersion.Action = func(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("My custom version info")
	return nil
}
app.VersionCommand = customVersion
```

### 默认命令

设置默认命令，当用户不提供命令时自动执行：

```go
app := cli.NewProgram("myapp", "1.0.0")
app.DefaultCommand = "serve" // 设置默认命令

serveCmd := cli.NewCommand("serve", "Start the server")
serveCmd.Flags.Int("port", 8080, "Port to listen on")
serveCmd.Action = func(ctx context.Context, cmd *cli.Command) error {
	port := cmd.Flags.Lookup("port").Value.String()
	fmt.Printf("Server listening on port %s\n", port)
	return nil
}

app.Commands = []*cli.Command{serveCmd}
```

使用示例：

```bash
# 不提供命令时，执行默认命令 serve
$ myapp
Server listening on port 8080

# 可以传递 flag 给默认命令
$ myapp --port 3000
Server listening on port 3000

# 显式指定其他命令
$ myapp help
# 显示帮助信息
```

## API 文档

### Program

```go
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
}

func NewProgram(appName, version string) *Program
func (p *Program) Run(args []string) error
func (p *Program) RunContext(ctx context.Context, args []string) error
func (p *Program) Get(name string) *Command
func (p *Program) SetOutput(w io.Writer)
func (p *Program) Output() io.Writer
func (p *Program) PrintUsage() error
```

### Command

```go
type Command struct {
	Name         string        // 命令名称（如 "init", "migrate"）
	Usage        string        // 命令用途简短描述（一行）
	Description  string        // 命令详细描述（多行）
	Flags        *flag.FlagSet // 命令标志集（用于定义和解析命令行参数）
	Action       ActionFunc    // 命令执行函数
	HideHelpFlag bool          // 是否隐藏 -h 帮助标志
}

func NewCommand(name, usage string) *Command
func DefaultHelpCommand() *Command
func DefaultVersionCommand() *Command
func (c *Command) Run(args []string) error
func (c *Command) RunContext(ctx context.Context, args []string) error
func (c *Command) SetOutput(w io.Writer)
func (c *Command) Output() io.Writer
func (c *Command) SetAppName(name string)
func (c *Command) PrintUsage() error
```

### ActionFunc

```go
type ActionFunc func(ctx context.Context, cmd *Command) error
```

## 示例项目

查看 [examples](./examples) 目录获取更多示例：

- [basic](./examples/basic) - 基础用法
- [flags](./examples/flags) - 标志使用
- [context](./examples/context) - Context 使用
- [custom](./examples/custom) - 自定义帮助和版本

## 测试

运行测试：

```bash
go test ./...
```

运行带覆盖率的测试：

```bash
go test -cover ./...
```

## 对比其他框架

| 特性 | cli | cobra | urfave/cli |
|------|-----|-------|-----------|
| 依赖 | 0 | 多个 | 0 |
| 基于标准库 | ✅ | ❌ | ✅ |
| 学习曲线 | 低 | 中 | 低 |
| 功能丰富度 | 中 | 高 | 高 |
| 代码量 | 极小 | 大 | 中 |

## 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

[hupeh](https://github.com/hupeh)
