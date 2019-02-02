#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
homedir="$workspace/src/github.com/OpenFinancing"

# Are we actually doing cleanup?
if [ "$1" = "clean" ]; then
    rm -rf $workspace
    exit
fi

# Create the homedir
if [ ! -L "$homedir/openfinancing" ]; then
    mkdir -p $homedir
    back=$(pwd)
    cd $homedir
    ln -s ../../../../../. OpenFinancing
    cd $back
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$homedir/openfinancing"
PWD="$homedir/openfinancing"

# Launch the arguments with the configured environment.
if [ ! -z "$1" ]; then
	exec "$@"
fi
