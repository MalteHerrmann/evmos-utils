# Upgrade a local Evmos node

This utility helps executing the necessary commands to prepare a
software upgrade proposal, submit it to a running local Evmos node,
and vote on the proposal.

Note, that this script is designed to work with a local node that was 
started by calling the `local_node.sh` script from the Evmos main repository.

## Installation

```bash
go install github.com/MalteHerrmann/upgrade-local-node-go@latest
```

## Usage

Start a local node by running `local_node.sh` from the Evmos repository.
In order to schedule an upgrade, run:

```bash
upgrade-local-node-go [TARGET_VERSION]
```

The target version must be specified in the format `vX.Y.Z(-rc*)`, e.g. `v13.0.0-rc2`.
