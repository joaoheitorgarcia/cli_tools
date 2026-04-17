# 1) Available tools
***

## `image_resizer`

Resize images from the terminal.

Current behavior:

- Accepts PNG, JPEG, and GIF images as input.
- Resizes using nearest-neighbor scaling.
- Lets you set width, height, or both.
- Preserves the original aspect ratio when one dimension is omitted.
- Writes the result next to the source file as `<original_name>_resized.png`.

### Flags

- `-f`: path to the input image
- `-w`: output width in pixels
- `-h`: output height in pixels

At least one of `-w` or `-h` must be provided.

## Requirements
- Go `1.26.2` or compatible local setup
- Dependency: `golang.org/x/image`

***

# 2) Build all tools

Use the repo-level build script to compile every Go CLI tool in the repository and export the binaries to a directory of your choice.

```bash
./build_and_export_all.sh ~/bin
```

Current behavior:

- Scans first-level subdirectories for Go modules.
- Runs tests for each module.
- Builds modules that contain a `main.go`.
- Writes each binary to the target directory using the folder name as the binary name.

With the current repo layout, this exports:

```text
~/bin/image_resizer
```
