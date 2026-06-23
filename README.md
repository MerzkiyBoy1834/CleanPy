# CleanPy

CleanPy is a simple project manager for Python that runs static analysis with pylint before running source files.

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
| `run` | Check and run the project |
| `clean` | Remove `.build` artifacts (legacy) |
| `version` | Show version information |

## Project layout

Put all Python code in `src/`. You can use multiple files and subdirectories:

```
src/
  main.py          # entry point
  helpers.py
  utils/
    __init__.py
    math_ops.py
data.txt           # paths in code are relative to project root
```

`clean-py run` starts Python with the project root as the working directory and adds `src/` to `PYTHONPATH`, so `open("data.txt")` resolves from the project root and `from utils.math_ops import add` works across modules.

## Cross-Platform

CleanPy works on both Linux and Windows.

- **Linux**: Compiled as a native binary using Go
- **Windows**: Build with `go build` and run

The tool does not use any platform-specific features, so it works identically on both operating systems. Projects run on any platform where Python 3 is installed.

## Requirements

- Python 3.x
- pylint: `pip install pylint`
- Go 1.21+
