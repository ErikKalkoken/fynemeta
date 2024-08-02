# fynemeta

A tool to help use Fyne metadata in the build process.

[![GitHub Release](https://img.shields.io/github/v/release/ErikKalkoken/fynemeta)](https://github.com/ErikKalkoken/fynemeta)
[![CI/CD](https://github.com/ErikKalkoken/fynemeta/actions/workflows/go.yml/badge.svg)](https://github.com/ErikKalkoken/fynemeta/actions/workflows/go.yml)
[![GitHub License](https://img.shields.io/github/license/ErikKalkoken/fynemeta)](https://github.com/ErikKalkoken/fynemeta)

## Description

fynemeta is a small command line tool that helps to use Fyne metadata in the build process.

- Print values from a Fyne metadata file to stdout
- Generate an AppStream metadata file from Fyne metadata

## Print metadata values

The `lookup` command allows us to print a value from the metadata file to stdout, so we can use it in a shell script.

Let's say we have the following definition for name in the `FyneApp.toml` file:

```toml
# ...
[Details]
name="Janice"
# ...
```

Then we can use fynemeta to read and assign that value to a variable in a shell script:

```bash
name=$(fynemeta lookup -k Details.name)
```

or

```bourne
name=`fynemeta lookup -k name`
```

## Generate AppStream metadata

> [!NOTE]
> We currently only support a subset of the AppStream specification, which is needed to pass all checks when creating an AppImage. If you need additional parameters please open a feature request or submit a PR.

The `generate` command generates an AppStream metadata file from a `FyneApp.toml` file. For that to work it needs to have some additional parameters.

Here is an example of a complete file:

```toml
Website = "https://github.com/ErikKalkoken/janice"

[Details]
  Icon = "icon.png"
  Name = "Janice"
  ID = "io.github.erikkalkoken.janice"
  Version = "0.2.3"
  Build = 2

[Release]
  BuildName = "janice"
  Description = "<p>A desktop app for viewing large JSON files.</p>"  # note that some HTML is allowed here
  License = "MIT"
  Screenshots = "https://cdn.imgpile.com/f/0IrYBjJ_xl.png" # optional, use comma as delimiter to define multiple urls
  ContentRating = "oars-1.1"  # optional

[LinuxAndBSD]
  GenericName = "JSON viewer"
  Categories = ["Utility"]
  Comment = "A desktop app for viewing large JSON files"
  Keywords = ["json", "viewer"]  # optional
```

> [!TIP]
> For a more detailed explanation on which values are allows for each parameter, please see the [AppStream Specification - Desktop Application](https://www.freedesktop.org/software/appstream/docs/sect-Metadata-Application.html) at freedesktop.org.
