# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

Given a version number MAJOR.MINOR.PATCH, increment the:

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

##### Types of changes

* [Added] for new features.
* [Changed] for changes in existing functionality.
* [Deprecated] for soon-to-be removed features.
* [Removed] for now removed features.
* [Fixed] for any bug fixes.
* [Security] in case of vulnerabilities.

## [Unreleased]
### Changed
- Updated `--dryrun` to `--dry-run`
- Refactored the cleaning process to be more easily testable
- AWS SDK updated from v1.23.5 to v1.25.2
- AWS SDK Config updated from v1.25.11 to v1.27.4
- AWS SDK Services updated from v1.47.2 to v1.51.1
- Rodaine Table updated from v1.1.0 to v1.1.1
- Viper updated from v1.17.0 to v1.18.2
- Refactor to compensate for AWS SDK updates


## [1.3.2] - 2023-12-06
### Changed
- AWS SDK updated from v1.17.8 to v1.23.5
- AWS SDK Config updated from v1.18.21 to v1.25.11
- AWS SDK Services updated from v1.32.0 to v1.47.2
- Fatih Color updated from v1.15.0 to v1.16.0
- Viper updated from v1.15.0 to v1.17.0
- Cobra updated from v1.7.0 to v1.8.0

### Fixed
- Resolved breaking change to S3 client


## [1.3.1] - 2023-04-21
### Changed
- Cobra updated from v1.6.1 to v1.7.0
- Fatih Color updated from v1.14.1 to v1.15.0
- AWS SDK updated from v1.17.5 to v1.17.8
- AWS SDK Config updated from v1.18.15 to v1.18.21 
- AWS SDK Services updated from v1.30.5 to v1.32.0
- Updated make linting method
- Tweet workflow updated to newer method


## [1.3.0] - 2023-03-01
### Added
- Updated to use Go 1.19 for future compatibility

### Changed
- Viper updated from v1.14.0 to v1.15.0
- Fatih Color updated from v1.13.0 to v1.14.1
- Rodaine Table updated from v1.0.1 to v1.1.0


## [1.2.0] - 2023-02-01
### Added
- Integrated Cobra to flesh out the cli options

### Changed
- Opened this upto MIT License
- Added version output to the make file for easier archiving of the compiled binaries
- Began setting up CI and Code Analysis
- AWS SDK updated from v1.16.15 to v1.17.5
- AWS SDK Config updated from v1.17.6 to v1.18.15
- AWS SDK Services updated from v1.27.10 to v1.30.5

## [1.1.0] - 2022-11-25
### Added
- Config file support instead of hard coding the options
- Verbose flag option to help debug what is happening
- Debugging messages through the code to help see what is happening

### Changed
- Split out some files for easier maintenance
- Result output easier to understand and read

## [1.0.0] - 2022-10-13
### Added
- Initial Launch
