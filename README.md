# vscode2kitty

I really like how my VS Code terminal looks with the amazing [SynthWave '84](https://github.com/robb0wen/synthwave-vscode/blob/master/themes/synthwave-color-theme.json) theme and I wanted it for my terminal of choice, [kitty](https://sw.kovidgoyal.net/kitty/). This is a pretty lame converter of VS Code themes to kitty themes. It does some stuff:

* Maps VS Code colors to kitty colors as best I can
   * Supports fallback colors and transforms
* Colors containing alpha are blended over the background color because VS Code colors can have an alpha component, but kitty colors can't
* VS Code terminal selection background and foreground are settable, but you can also just have an alpha-blended selection, as SynthWave '84 does. That means that there is no selection foreground. This code just makes that up by inverting the selection background.

## Installation

If you have the Go language installed, this tool can be installed the usual way:

```shell
go install github.com/dangermike/vscode2kitty@latest
```

If you do not, install the Go language and return to the top of this section.

## Usage

The converter can be called with either a local file or a URL. Results will be emitted to stdout. If, like me, you have `include theme.conf` in your `~/.config/kitty/kitty.conf` file, you can just run this:

```shell
vscode2kitty https://raw.githubusercontent.com/robb0wen/synthwave-vscode/master/themes/synthwave-color-theme.json | tee ~/.config/kitty/theme.conf
```
