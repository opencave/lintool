# lintool

![License](https://img.shields.io/badge/license-Apache2.0-green)
![Language](https://img.shields.io/badge/Language-Go-blue.svg)
[![version](https://img.shields.io/github/v/tag/openholes/bencode?label=release&color=blue)](https://github.com/openholes/bencode/releases)
[![Go report](https://goreportcard.com/badge/github.com/openholes/bencode)](https://goreportcard.com/report/github.com/openholes/bencode)
[![Go Reference](https://pkg.go.dev/badge/github.com/openholes/bencode.svg)](https://pkg.go.dev/github.com/openholes/bencode)

`lintool` is a lint tool for project source code.

## About openHoles

openHoles is an open source organization focusing on peer-to-peer solutions, find more information [here](https://github.com/openholes/openholes)

## Install

```bash
go install github.com/openholes/lintool@latest
```

## Usage

Use the `-h` parameter to obtain usage help.

```bash
> lintool -h
lintool is a lint tool

Usage:
  lintool [command]

Examples:
lintool blankline

Available Commands:
  blankline   check if file ends with a blank line
  license     check if the source code files have a license header
  help        Help about any command

Flags:
  -h, --help   help for lintool

Use "lintool [command] --help" for more information about a command.
```

Use the `blankline` command to detect blank lines in files.

```bash
> lintool blankline -h
check if file ends with a blank line

Usage:
  lintool blankline [flags]

Aliases:
  blankline, bl

Examples:
lintool bl -d . -e .idea,testdata

Flags:
  -d, --directory string   directory to check (default ".")
  -e, --exclude strings    directories or files to exclude (comma-separated)
  -h, --help               help for blankline

> # check /tmp/demo directory and exclude .idea
> lintool blankline -d /tmp/demo -e .idea
the following files are not end with a blank line:
testdata/blankline/noline.txt
2025/03/30 19:12:59 ERROR lintool execute failed reason="blank line issue found"
exit status 1
```


Use the `license` command to detect license header in files.

```bash
> lintool license -h
check if the source code files have a license header

Usage:
  lintool license [flags]

Aliases:
  license, lic

Examples:
lintool lic -d . -e .idea,testdata

Flags:
  -d, --directory string     directory to check (default ".")
  -e, --exclude strings      directories or files to exclude (comma-separated)
  -x, --extensions strings   file extension to check, if not set, all files will be checked (comma-separated)
  -h, --help                 help for license
  -l, --license string       license file to use, if no license file is found and this flag is not set, process will be skipped, support [Apache-2.0, MIT, GPL-2.0, GPL-3.0, LGPL, MPL, BSD]

> # check /tmp/demo directory and exclude .idea
> lintool license -d /tmp/demo -e .idea
the following files are not end with a blank line:
testdata/blankline/noline.txt
2025/03/30 19:12:59 ERROR lintool execute failed reason="blank line issue found"
exit status 1
```

## License

openHoles is licensed under the Apache License 2.0. Refer to [LICENSE](https://github.com/openholes/bencode/blob/main/LICENSE) for more details.
