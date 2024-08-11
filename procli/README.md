# procli

Procli is a Debian/Ubuntu machine CLI provisioning tool that installs some of the most basic packages, QoL changes, and tools that I prefer to have on most of my machines. 

Procli is by default configured to download the following linux packages, OhMyZsh custom plugins, and extra repositories for tools. This can be adjusted with some mediocre knowledge of Go. 

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

Procli also, by default, performs a `git clone` or `git pull` on a remote GitHub repository and stores the repo at `$HOME/dotties` wherein several custom dotfiles can be cloned to the machine for easy dotfile standardization. By default Procli only moves a `.zshrc` file from `$HOME/dotties` > `$HOME/.zshrc`.
