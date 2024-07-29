# tomlq

Print a value from a TOML file.

## Description

tomlq is a small command line tool that can print a value from a TOML file. It's main purpose is to enable shell scripts to use values from a TOML file.

Example:

Let's say we have a TOML file with the name `config.toml`:

```toml
name=Charlie
```

Then we can use tomlq to read and assign that value to a variable in a shell script:

```bash
name=$(tomlq -k name config.toml)
```

For all features and options please run `tomlq -h`.
