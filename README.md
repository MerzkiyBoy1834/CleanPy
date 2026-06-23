# CleanPy

CleanPy is a simple project manager for Python that compiles code to .pyc and runs static analysis with pylint before building.

## Installation

```bash
git clone https://github.com/MerzkiyBoy1834/CleanPy.git
cd CleanPy
go build
sudo cp clean-py /usr/local/bin/
```

## Usage

```bash
clean-py new <project_name>
cd <project_name>
clean-py check
clean-py run
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize project in current directory |
| `new <name>` | Create new project with given name |
| `check` | Run pylint static analysis |
| `run` | Run the project |
| `clean` | Remove build artifacts |
| `version` | Show version information |

## Cross-Platform

CleanPy works on both Linux and Windows.

- **Linux**: Compiled as a native binary using Go
- **Windows**: Build with `go build` and run

The tool does not use any platform-specific features, so it works identically on both operating systems. Python code compiled with CleanPy runs on any platform where Python 3 is installed.

## Requirements

- Python 3.x
- pylint: `pip install pylint`
- Go 1.21+
