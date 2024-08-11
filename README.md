# go-tools

A collection of Go tools that I have created for my own functionality. All can be used with, at most, a mediocre amount of Go experience and adjustment.

### Procli

Procli is a Debian/Ubuntu machine CLI provisioning tool that installs some of the most basic QoL changes, packages, and tools that I prefer to have on most of my machines. 

Linux packages installed:
- `Zsh`
- `Git`
- `Nvim`
- `fzf` (for fzf searching via OhMyZsh)

OhMyZsh Custom Plugins & Customizations:
- `zsh-autosuggestions`
- `zsh-syntax-highlighting`
- `fzf-zsh-plugin`

Extra Tools Downloaded

- `/adrian-griffin/go-tools/`

Procli also, by default, perfoms a `git clone` or `git pull` on a remote GitHub repository and stores the repo at `$HOME/dotties` wherein several custom dotfiles can be cloned to the machine for easy dotfile standardization. By default Procli only moves a `.zshrc` file from `$HOME/dotties` > `$HOME/.zshrc`, but with some customization of the Go source, can be adjusted to clone and distribute dotfiles as needed.
