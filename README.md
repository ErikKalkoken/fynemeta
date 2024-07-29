# tomlq

Print a value from a TOML file.

[![GitHub Release](https://img.shields.io/github/v/release/ErikKalkoken/tomlq)](https://github.com/ErikKalkoken/tomlq)
[![build status](https://github.com/ErikKalkoken/tomlq/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/ErikKalkoken/tomlq/actions/workflows/ci-cd.yml)
[![GitHub License](https://img.shields.io/github/license/ErikKalkoken/tomlq)](https://github.com/ErikKalkoken/tomlq)

## Description

tomlq is a small command line tool that can print a value from a TOML file. It's main purpose is to make it possible to use values from a TOML file in a shell script.

All output is printed to stdout. Date values are formatted as RFC3339.

## Example

Let's say we have a TOML file with the name `config.toml`:

```toml
name=Charlie
```

Then we can use tomlq to read and assign that value to a variable in a UNIX style shell script:

```bash
name=$(tomlq -k name config.toml)
```
