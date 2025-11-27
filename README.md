# precompress

[![tests](https://github.com/akornatskyy/precompress/actions/workflows/tests.yml/badge.svg)](https://github.com/akornatskyy/precompress/actions/workflows/tests.yml)

*Precompress* is a lightweight command-line tool that generates precompressed versions of your static assets using Brotli, Gzip, and Zstandard. By producing ready-to-serve compressed files during your build or deployment process, it helps optimize web performance, reduce server load, and simplify production pipelines.

Precompress is ideal for static websites, SPAs, CDNs, edge deployments, and any environment where serving precompressed assets improves efficiency. It integrates easily with build systems, Docker images, and CI pipelines.

## Features

* Multi-algorithm compression - Brotli (.br), Gzip (.gz), and Zstandard (.zst).
* Pre-generate assets once - eliminates runtime compression overhead.
* Optimized for production - reduces transfer size and speeds up delivery.
* Directory recursion - process entire folders of files automatically.
* Configurable limits - minimum file size and maximum recursion depth.
* Simple CLI - no config files needed.

## Installation

You can download prebuilt binaries or install from source.

Using Go:

```sh
go install github.com/akornatskyy/precompress@latest
```

## Usage

```text
Usage:
  precompress [options] <files,dirs,...>

Description:
  Generates precompressed versions of static files using Brotli, Gzip, and
  Zstandard algorithms for faster content delivery in production environments.

Options:
  --min-size <bytes>     Minimum content size (default 1024)
  --max-depth <levels>   Descend at most levels of directories
                         (default 0 means no limit)
  -h, --help             Show this help message
  -v, --version          Show version

Examples:
  precompress public
```

## How It Works

For each file encountered:

1. Skips files smaller than `--min-size`.
2. Generates:
   * `filename.ext.br` (Brotli).
   * `filename.ext.gz` (Gzip).
   * `filename.ext.zst` (Zstandard).
3. Leaves original files untouched.
4. Avoids recompressing if the output file already exists and is newer.

These precompressed files can then be served directly by web servers or CDN providers.
