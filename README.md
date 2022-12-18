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

### Install via Snap

*Ignoreinit* is available as a [snap](https://snapcraft.io/ignoreinit) for *amd64* and *arm64* based systems.

```bash
sudo snap install ignoreinit
```

### Install via git or go

You may either install the executable directly into `$GOPATH` or download it from the [release page](https://github.com/Loupeznik/ignoreinit/releases).

```bash
git clone https://github.com/Loupeznik/ignoreinit
cd ignoreinit
go install github.com/loupeznik/ignoreinit
```

Or simply install the latest version with Go without needing to clone the repo:

```bash
go install github.com/loupeznik/ignoreinit@latest
```

Docker and general package manager support are coming at a later date.
