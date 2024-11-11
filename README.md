# go-tools

✌ Some go tools ive made for my convenience with my homelab and selfhosted setups. All can be used with, at most, a mediocre amount of Go experience and a lil adjustment, feel free to copy, reuse, and redistribute in any capacity. All tools are built for debian and debian-based machines, all CLI based.

Raw source, Go will need to be installed on your machine to build into executables. Please visit Go's documentation for installation instructions, here: https://go.dev/doc/install

---

### Set up and use one of these scripts

Git clone (https) and cd into designated tool dir
```shell
·> git clone https://github.com/adrian-griffin/go-tools.git && cd go-tools/ipx
```

Go build into executable
```shell
·> go build ipx
```

Run tool
```shell
·> ipx
--- Default route & Network Information ---
-------------------------------------------
Default Gateway  : 10.115.128.1
LAN IP           : 10.115.128.88
LAN Interface    : eth2
DNS Servers      : 10.115.128.8
NAT/Public IP    : <ip>
-------------------------------------------
```

Optionally, set rcfile alias to allow calling the command remotely from anywhere on the machine
```shell
alias ipx="$HOME/go-tools/ipx/ipx"
```

### baxup
Performs reliable backups on a target directory or Docker compose container. Specifically, data is compressed and copied to a remote machine and/or stored in another directory on the local machine. Reliability in data consistency and security while in transport via rsync over SSH (with data validation checksums) using the built in `remote-send` function. Optionally, with "Docker mode" also collects docker image digests and version information and stores this information alongside the `docker-compose.yml` file before compression. This allows better recovery reliability as specific docker sha256 digest hashes can be used to ensure recovery efforts of data are spun up with the exact docker image as stored with the previous backup.

Better yet, so long as SSH keys are set up between the local and remote machine(s), `baxup` can be used in conjunction with cronjobs to allow periodic, scheduled backups of directories or docker containers that ensures data consistency with docker services, and secure and reliable remote transfer of compressed backup data with checksum validation and SSH encryptions enforced. 

### ipx 
Like `ip a` but better (?)

A quick debian-based CLI IP & network information collection tool. It quickly collects the primary network interface name and address, primary route public/NAT address, DNS servers, and quickly checks for DNS leaks. IPv4 only.

### procli
Stupid name based on Provision CLI.

Debian-based CLI provisioning tool that auto-installs some of my most basic and needed linux packages, such as OhMyZSH, nvim, git, etc., as well as some OhMyZsh customizations/plugins. A GitLab or GitHub repository dotfile URL can be used supplied in order to sync a `~/.zshrc` file for use with OMZ. A little more rough around the edges and may need some extra elbow grease to adjust for anyone else's needs, but it can be done. 

### qlip
Named partially after the astral-human border region in the world of Berserk, qliphoth (yeah, yeah) 

Literally just a basic go clipboard program, but I've found some use for making it a "cli clipboard" of sorts, certainly not secure so don't throw any secrets in there, but useful for commonly copied and pasted CLI commands, for example.
