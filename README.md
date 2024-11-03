# go-tools

✌ Some go tools ive made for my convenience with my homelab and selfhosted setups. All can be used with, at most, a mediocre amount of Go experience and a lil adjustment, feel free to copy, reuse, and redistribute in any capacity. All tools are built for debian and debian-based machines, all CLI based.

### baxup
Performs backups on a target directory or Docker compose container. Specifically, data is tarballed and either copied to a remote machine, and stored in another directory on the local machine. Data can be reliably and securely copied to a remote machine via rsync over SSH (with data validation checksums) using the built in `remote-send` function. "Docker mode" also collects docker image digests and version information and stores this information alongside the `docker-compose.yml` file before compression at docker container termination time.
Better yet, so long as SSH keys are set up between the local and remote machine, `baxup` can be used in conjunction with cronjobs to allow periodic, scheduled backups of directories or docker containers that ensures data consistency with docker services, and secure and reliable remote transfer of compressed backup data with checksum validation and SSH encryptions enforced. 

### ipx 
Like `ip a` but better (?)
A quick debian-based CLI IP & network information collection tool. It quickly collects the primary network interface name and address, primary route public/NAT address, DNS servers, and quickly checks for DNS leaks. IPv4 only.

### procli
Stupid name based on Provision CLI.
Debian-based CLI provisioning tool that auto-installs some of my most basic and needed linux packages, such as OhMyZSH, nvim, git, etc., as well as some OhMyZsh customizations/plugins. A GitLab or GitHub repository dotfile URL can be used supplied in order to sync a `~/.zshrc` file for use with OMZ. A little more rough around the edges and may need some extra elbow grease to adjust for anyone else's needs, but it can be done. 

### qlip
Named partially after the astral-human border region in the world of Berserk, qliphoth (yeah, yeah) 
Literally a basic a** go clipboard program, but I've found some use for making it a "cli clipboard" of sorts, certainly not secure so don't throw any secrets in there, but useful for commonly copied and pasted CLI commands, for example.
