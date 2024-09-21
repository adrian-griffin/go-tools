# go-tools

A collection of various Go tools that I have created for my own functionality. All can be used with, at most, a mediocre amount of Go experience and adjustment.

### baxup
Baxup performs efficient backups on target docker containers or directories. Data is compressed and stored locally by default, but can also be copied to a remote machine securely via rsync using the built in `remote-send` function. 

### ipx 
ipx is a quick CLI IP & network information collection tool. It quickly collects primary interface addresses, public/NAT addresses, DNS servers, and quickly checks for DNS leaks from the local machine.

### procli
Procli is a CLI provisioning tool that auto-installs some of the most basic and needed linux packages and OhMyZsh customizations/plugins, as well as all of the `/adrian-griffin/go-tools/` Go tools. 

### qlip
Named partially after the astral-human border region in the world of Berserk, qliphoth, qlip is a CLI clipboard tool that stores whatever text you would like in a local text file, allowing you to list, add, or remove items from its array via a quick command. This can be useful for commands that are commonly needing to be ran that are easier pasted than typed. 
