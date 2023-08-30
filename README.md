# Evmos Dev Utils

This tool contains several utility functionalities that are useful during
development of the Evmos blockchain. 

At the core, all interactions go through the Evmos CLI interface, which is
called from within the Go code.

Note, that this script is designed to work with a local node that was
started by calling the `local_node.sh` script from the Evmos main repository.

## Installation

In order to install the tool, clone the source and install locally.
Note, that using `go install github.com/MalteHerrmann/upgrade-local-node-go@latest`
does not work because of the replace directives in `go.mod`,
which are necessary for the Evmos dependencies.

```bash
git clone https://github.com/MalteHerrmann/upgrade-local-node-go.git
cd upgrade-local-node-go
make install
```

## Features

### Upgrade Local Node

The tool creates and submits a software upgrade proposal to a running local Evmos node,
and votes on the proposal. To do so, run:

```bash
upgrade-local-node-go [TARGET_VERSION]
```

The target version must be specified in the format `vX.Y.Z(-rc*)`, e.g. `v13.0.0-rc2`.

### Vote on Proposal

The tool can vote with all keys from the configured keyring, that have delegations
to validators. This can either target the most recent proposal, or a specific one when
passing an ID to the command.

```bash
upgrade-local-node-go vote [PROPOSAL_ID]
```
