# Go Format nginx Config

A powerful and flexible nginx configuration formatter that recursively processes nginx configuration files in a directory structure.
This tool helps maintain consistent formatting across your nginx configuration files.

## Features

- Recursively processes nginx configuration files in a directory
- Configurable indentation
- Optional comment removal
- Concurrent file processing for better performance
- Automatic backup creation before modifications
- Dry-run mode for previewing changes
- Customisable file extensions
- Detailed statistics and logging

## Installation

```bash
go install github.com/ChrisMcKee/gofmtnginx@latest
```

## Usage

```bash
gofmtnginx [flags] <directory>
```

### Flags

- `--removecomments`: Remove comments from the configuration file (default: false)
- `--indent`: Number of spaces for indentation (default: 2)
- `--dry-run`: Show what would be done without making changes (default: false)
- `--verbose`: Enable verbose logging (default: false)
- `--backup`: Create backup files before modifying (default: false)
- `--concurrent`: Process files concurrently (default: true)
- `--workers`: Number of concurrent workers (default: 4)
- `--extensions`: Comma-separated list of file extensions to process (default: ".conf,.proxy")
- `--preserve-newlines`: Preserve existing newlines between blocks (default: false)

### Examples

Format all nginx configuration files in a directory:
```bash
gofmtnginx /etc/nginx
```

Format only .conf files with 4-space indentation:
```bash
gofmtnginx --indent=4 --extensions=.conf /etc/nginx
```

Preview changes without modifying files:
```bash
gofmtnginx --dry-run /etc/nginx
```

Remove comments and process specific file types:
```bash
gofmtnginx --removecomments --extensions=.conf,.nginx /etc/nginx
```

## Output

The tool provides statistics about the formatting process:
- Number of files processed
- Number of files skipped
- Number of files failed
- Total processing time

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache License 2.0
