# gompare

A simple utility for validating CSV columns

# Building

In project directly, run

```shell
go build
```

If using **make**, the following targets are available

- ``make`` or ``make build`` - Build
- ``make install`` - Install binary to bin path
- ``make clean`` - Clean
- ``make tidy`` - Tidy up dependencies
- ``make good-test`` - Quick validity test. Results should be `Valid`
- ``make bad-test`` - Quick validity test. Results should be 'Invalid'

# Usage

```shell
./gompare --template-file=template.csv --input-file=test-bad.csv
```

Type ```./gompare --help``` for more information

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.3] - 2021-11-12

### Added

- Added Makefile
- Added ability to show non-matching columns present in target CSV files
- Release workflow

## [0.0.2] - 2021-11-11

### Fixed

- Code cleanup. Remove unused constant

## [0.0.1] - 2021-11-10

### Added

- Initial setup

[0.0.1]: https://github.com/SharkFourSix/gompare/tree/v0.0.1

[0.0.2]: https://github.com/SharkFourSix/gompare/tree/v0.0.2

[0.0.3]: https://github.com/SharkFourSix/gompare/tree/v0.0.3