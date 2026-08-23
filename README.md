# Umafetch

> Uma Musume character outfits for [fastfetch](https://github.com/fastfetch-cli/fastfetch).

![go](https://img.shields.io/badge/language-go-blue?style=flat-square)
![fastfetch](https://img.shields.io/badge/requires-fastfetch-blue?style=flat-square)
![license](https://img.shields.io/badge/license-MIT-lightgrey?style=flat-square)

Umafetch is a terminal application written in Go that generates dynamic and custom fastfetch configurations themed around the Uma Musume franchise, extracting representative color palettes from their outfits to match your terminal aesthetics.

## Requirements

- [Go](https://go.dev/) (to build and install)
- [Fastfetch](https://github.com/fastfetch-cli/fastfetch) (to display system info and images)
- A terminal with image rendering support (such as Kitty, Konsole, WezTerm, Ghostty, etc.)

## Installation

You can install umafetch directly using Go's install command:

```bash
go install github.com/operaodev/umafetch@latest
```

Make sure you have your Go binary directory in your `PATH` environment variable. This is usually located at `~/go/bin`.

Once the binary is installed, you need to download the character assets (images and metadata) by running the following command:

```bash
umafetch install
```

## Usage

To use umafetch, run the `render` command to update the fastfetch configuration, and then run `fastfetch` normally:

```bash
umafetch render && fastfetch
```

You can also customize the selection when rendering by using flags before running fastfetch:

```bash
umafetch render --name "Special Week" --template small && fastfetch
```

## Available Commands

The umafetch CLI provides the following commands:

### umafetch install
Downloads all Uma Musume character images and their required configuration files, processes representative colors from their outfits, and stores them locally in your user directory.

```bash
umafetch install
```

### umafetch render
Generates the theme configuration based on a character (selected randomly or explicitly) and prints the path of the generated configuration file for fastfetch to consume.

```bash
umafetch render [flags]
```

#### Available Flags for render:

- `--name <string>`: Forces the selection of a specific Uma Musume by her name (e.g., `"Special Week"`).
- `--outfit <int>`: Forces the selection of a specific outfit/costume by its index number.
- `--template <string>`: Overrides the fastfetch layout template to use. Valid values: `large` (default, with a large avatar and bio) or `small` (compact).

### umafetch umas
Lists, searches, and selects available Umas in your local database either interactively or directly to set your theme.

```bash
# Show general count of installed Umas and outfits
umafetch umas

# List all available Umas and their outfits
umafetch umas --full

# Search for an Uma Musume by her name and view her outfits
umafetch umas --search "Tokai Teio"

# Search and also display representative color details of the outfits
umafetch umas --search "Tokai Teio" --full

# Permanently select a character in your configuration file
umafetch umas --target "Special Week"         # Random outfit
umafetch umas --target "Special Week, 1"      # Specific outfit (ID 1)
```

#### Available Flags for umas:

- `-f, --full`: Shows the complete, detailed list of all characters alongside their outfits and representative colors (as RGB terminal background formats).
- `--search <string>`: Filters and shows characters matching the query.
- `--target <string>`: Statically sets the default Uma in the configuration file. Format: `"Name"` or `"Name,Outfit"`.

### umafetch template
Manages custom fastfetch layout templates (`large` and `small`).

```bash
# Generate both default templates (large and small)
umafetch template

# Quickly switch the active template between large and small
umafetch template --switch

# Force the generation/restoration of a specific template
umafetch template --gen-template "large"
```

#### Available Flags for template:

- `--gen-template <string>`: Generates the selected template file (`large` or `small`) on disk.
- `--switch`: Switches the active default template in your local configuration.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
