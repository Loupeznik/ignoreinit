# Ignoreinit

![GitHub](https://img.shields.io/github/license/loupeznik/ignoreinit?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/loupeznik/ignoreinit?style=for-the-badge)

Ignoreinit is a tool for creating .gitignore files from the command line. Gitignore files are pulled from [github/gitignore](https://github.com/github/gitignore) repo.

## Usage

### List available templates

Prints the available templates from [github/gitignore](https://github.com/github/gitignore).

```bash
ignoreinit list
```

### Search available templates

Finds matching templates by exact name, partial name, or close spelling.

```bash
ignoreinit search <term>
```

### Create new .gitignore

Creates new .gitignore based on given language in defined location (either relative or absolute).

```bash
ignoreinit init <template...> <location>
```

If location is omitted, ignoreinit writes to the current directory. Pass multiple templates to combine them:

```bash
ignoreinit init go node terraform
ignoreinit init go node terraform ./my-project
```

### Replace existing gitignore

Replaces existing .gitignore based on given language in defined location (either relative or absolute).

```bash
ignoreinit replace <template...> <location>
```

### Merge existing gitignore

Merges a gitignore for given language into existing .gitignore in defined location (either relative or absolute).

```bash
ignoreinit merge <template...> <location>
```

## Build from source

```bash
git clone https://github.com/Loupeznik/ignoreinit
cd ignoreinit
go build -o build/ignoreinit github.com/loupeznik/ignoreinit
```

## Run with Docker

```bash
# Create .gitignore in the current directory
docker run --rm --user "$(id -u):$(id -g)" -v ${PWD}:/work loupeznik/ignoreinit:latest init go .

# Create .gitignore in another directory directory
docker run --rm --user "$(id -u):$(id -g)" -v $HOME/projects:/work loupeznik/ignoreinit:latest init go .
```

## Install

### Install via Snap

[![ignoreinit](https://snapcraft.io/ignoreinit/badge.svg)](https://snapcraft.io/ignoreinit)

*Ignoreinit* is available as a [snap](https://snapcraft.io/ignoreinit) for *amd64* and *arm64* based systems.

```bash
sudo snap install ignoreinit
```

### Install via Homebrew

```bash
brew install --cask loupeznik/tap/ignoreinit
```

### Install via AUR

```bash
yay -S ignoreinit-bin
```

### Install via go

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
