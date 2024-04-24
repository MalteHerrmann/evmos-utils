<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## Unreleased

### Improvements

- [#32](https://github.com/MalteHerrmann/evmos-utils/pull/32) Minor refactor in CLI commands
- [#35](https://github.com/MalteHerrman/evmos-utils/pull/35) Update to Evmos v17.
- [#38](https://github.com/MalteHerrmann/evmos-utils/pull/38) Add flags to CLI commands to enable more configuration.

## [v0.4.0](https://github.com/MalteHerrmann/evmos-utils/releases/tag/v0.4.0) - 2023-12-18

### Features

- [#26](https://github.com/MalteHerrmann/evmos-utils/pull/26) Enable just depositing with the binary

### Improvements

- [#27](https://github.com/MalteHerrmann/evmos-utils/pull/27) Add MIT license
- [#28](https://github.com/MalteHerrmann/evmos-utils/pull/28) Use [Cobra CLI](https://github.com/spf13/cobra) package
- [#29](https://github.com/MalteHerrmann/evmos-utils/pull/29) Adjust repository name from `upgrade-local-node-go` to `evmos-utils`
- [#30](https://github.com/MalteHerrmann/evmos-utils/pull/30) Use [zerolog](https://github.com/rs/zerolog) for logging

## [v0.3.0](https://github.com/MalteHerrmann/evmos-utils/releases/tag/v0.3.0) - 2023-08-30

### Features

- [#14](https://github.com/MalteHerrmann/evmos-utils/pull/14) Enable just voting with the binary

### Improvements

- [#7](https://github.com/MalteHerrmann/evmos-utils/pull/7) Add linters plus corresponding refactors
- [#6](https://github.com/MalteHerrmann/evmos-utils/pull/6) Restructuring and refactoring
- [#4](https://github.com/MalteHerrmann/evmos-utils/pull/4) Add GH actions and Makefile for testing
- [#3](https://github.com/MalteHerrmann/evmos-utils/pull/3) Use broadcast mode `sync` instead of `block`
- [#2](https://github.com/MalteHerrmann/evmos-utils/pull/2) Only vote if account has delegations

## [v0.2.0](https://github.com/MalteHerrmann/evmos-utils/releases/tag/v0.2.0) - 2023-08-09

### Improvements

- [#1](https://github.com/MalteHerrmann/evmos-utils/pull/1) Adaptively gets keys and current proposal ID from the local node

## [v0.1.0](https://github.com/MalteHerrmann/evmos-utils/releases/tag/v0.1.0) - 2023-08-01

### Features

- Gets current block height of local node (at `http://localhost:26657`)
- Submit a software upgrade proposal to a running local Evmos node for the target version
- Vote on the software proposal
