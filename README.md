# Ignoreinit

![GitHub](https://img.shields.io/github/license/loupeznik/ignoreinit?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/loupeznik/ignoreinit?style=for-the-badge)

Ignoreinit is a tool for creating .gitignore files from the command line. Gitignore files are pulled from [github/gitignore](https://github.com/github/gitignore) repo.

## Usage

### Create new .gitignore

Creates new .gitignore based on given language in defined location (either relative or absolute).

```bash
ignoreinit init <language> <location>
```

### Replace existing gitignore

Replaces existing .gitignore based on given language in defined location (either relative or absolute).

```bash
ignoreinit replace <language> <location>
```

## Build from source

```bash
git clone https://github.com/Loupeznik/ignoreinit
cd ignoreinit
go build -o build/ignoreinit github.com/loupeznik/ignoreinit
```

## Install

You may either install the executable directly into `$GOPATH` or download it from the [release page](https://github.com/Loupeznik/ignoreinit/releases).

```bash
git clone https://github.com/Loupeznik/ignoreinit
cd ignoreinit
go install github.com/loupeznik/ignoreinit
```

Docker and package manager support is coming come later.
