<img width="300" alt="image" src="https://github.com/user-attachments/assets/6458f5b3-835a-49e3-817e-bdab93b7141d" />

# Umafetch

> Uma Musume character outfits for [fastfetch](https://github.com/fastfetch-cli/fastfetch).

![go](https://img.shields.io/badge/language-go-blue?style=flat-square)
![fastfetch](https://img.shields.io/badge/requires-fastfetch-blue?style=flat-square)
![license](https://img.shields.io/badge/license-MIT-lightgrey?style=flat-square)

Umafetch is a terminal application written in Go that generates dynamic and custom fastfetch configurations themed around the Uma Musume franchise, extracting representative color palettes from their outfits to match your terminal aesthetics.

<img width="1099" alt="image" src="https://github.com/user-attachments/assets/c61cd40b-647e-422d-aec6-edb6fc5560c2" />


## Requirements

- [Go](https://go.dev/) (to build and install)
- [Fastfetch](https://github.com/fastfetch-cli/fastfetch) (to display system info and images)
- A terminal with image rendering support (such as Kitty, Konsole, WezTerm, Ghostty, etc.)

## Installation

You can install umafetch directly using Go's install command:

    go install github.com/operaodev/umafetch@latest

Make sure you have your Go binary directory in your `PATH` environment variable. This is usually located at `~/go/bin`.

If `umafetch` is not recognized, you can run it directly using its full path:

    ~/go/bin/umafetch

On Linux, you can also create a symbolic link to make `umafetch` available system-wide:

    sudo ln -s ~/go/bin/umafetch /usr/local/bin/umafetch

Once the binary is installed, you need to download the character assets (images and metadata) by running:

    umafetch install

## Usage

To use umafetch, run the `render` command to update the fastfetch configuration, and then run `fastfetch` normally:

```bash
umafetch render && fastfetch
```

You can also customize the selection when rendering by using flags before running fastfetch:

```bash
umafetch render --name "Special Week" --template small && fastfetch
```

## Configuration

Umafetch configurations and layout templates are stored in your home directory:
- Main configuration file: `~/.config/umafetch/config.json`
- Templates directory: `~/.config/umafetch/templates/`

### Global Configuration (`config.json`)

Here is the default structure of the `config.json` file:

```json
{
  "template": "large",
  "separator": {
    "width": 52,
    "decorator": "─"
  },
  "theme": {
    "name": null,
    "outfit": null
  }
}
```

#### Fields Description:

- `template`: Defines the default fastfetch layout template to use (`large` or `small`).
- `separator`: Customizes the visual divider line generated between section modules.
  - `width`: The width (in characters) of the divider line.
  - `decorator`: The character or ANSI symbol used to construct the line (defaults to `"─"`).
- `theme`: Configures character selection behavior.
  - `name`: The name of a specific Uma Musume to display (e.g. `"Special Week"`). Set to `null` to select a random character on every run.
  - `outfit`: The index of the specific outfit/costume to load (e.g. `0` for Tracen Academy uniform, `1` for first special outfit). Set to `null` to select a random outfit.

### Editing Templates

The fastfetch configuration templates (`config_large.jsonc` and `config_small.jsonc`) are fully editable. You can modify the modules, borders, or any other fastfetch settings as you like.

However, you must keep the following placeholders intact, as they are programmatically replaced by umafetch during rendering:

- `{Image}`: Path to the character's image.
- `{PrimaryColor}`: Hex code for the primary color.
- `{SecondaryColor}`: Hex code for the secondary color.
- `{Name}`: The character's name.
- `{Title}`: The outfit title.
- `{SloganLines}`: Character's slogan formatted for the layout.
- `{BioLines}`: Character's profile bio formatted for the layout.
- `{separator}`: Customized separator module based on the character's primary color.

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
