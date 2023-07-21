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

### Merge existing gitignore

Merges a gitignore for given language into existing .gitignore in defined location (either relative or absolute).

```bash
ignoreinit merge <language> <location>
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
docker run --rm -v ${PWD}:/work loupeznik/ignoreinit:latest init go .

# Create .gitignore in another directory directory
docker run --rm -v $HOME/projects:/work loupeznik/ignoreinit:latest init go .
```

## Install

### Install via Snap

[![ignoreinit](https://snapcraft.io/ignoreinit/badge.svg)](https://snapcraft.io/ignoreinit)

*Ignoreinit* is available as a [snap](https://snapcraft.io/ignoreinit) for *amd64* and *arm64* based systems.

```bash
sudo snap install ignoreinit
```

### Install via apt

On Debian-based distros, *ignoreinit* can be installed via *apt* using a custom *apt* repo. This option is currently supported for *amd64* systems.

```bash
sudo -s

apt install -y curl gpg

curl -fsSL https://apt.dzarsky.eu/apt-repo-dzarsky.gpg | gpg --dearmor > /etc/apt/trusted.gpg.d/apt-repo-dzarsky.gpg
echo "deb [arch=amd64 signed-by=/etc/apt/trusted.gpg.d/apt-repo-dzarsky.gpg] https://apt.dzarsky.eu /" > /etc/apt/sources.list.d/apt-repo-dzarsky.list

apt update
apt install ignoreinit
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
