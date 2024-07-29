# tomlq

Print a value from a TOML file.

[![GitHub Release](https://img.shields.io/github/v/release/ErikKalkoken/tomlq)](https://github.com/ErikKalkoken/tomlq)
[![CI/CD](https://github.com/ErikKalkoken/tomlq/actions/workflows/go.yml/badge.svg)](https://github.com/ErikKalkoken/tomlq/actions/workflows/go.yml)
[![GitHub License](https://img.shields.io/github/license/ErikKalkoken/tomlq)](https://github.com/ErikKalkoken/tomlq)

## Description

tomlq is a small command line tool that can print a value from a TOML file. It's main purpose is to make it possible to use values from a TOML file in a shell script.

All output is printed to stdout. Date values are formatted as RFC3339.

tomlq will exit with a non-zero value if it encounters an error (e.g. value not found). We recommend enabling the "strict error" mode in shell scripts to ensure those errors are not omitted:

```sh
set -e
```

## Example

Let's say we have a TOML file with the name `config.toml`:

```toml
name="Charlie"
```

Then we can use tomlq to read and assign that value to a variable in a UNIX style shell script:

```bash
name=$(tomlq -k name config.toml)
```

## Usage

To get the complete usage documentation run the following:

```sh
tomlq -h
```
