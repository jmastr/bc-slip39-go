#!/bin/bash

# Get BlockchainCommons source code
git submodule update --init --recursive

# We place the BlockchainCommons libraries into our working directory
export CFLAGS=-I/usr/local/include/
export LDFLAGS=-L/usr/local/lib/

# Compile dependencies
for directory in "bc-crypto-base" "bc-shamir" "bc-slip39"; do
  pushd blockchaincommons/${directory}/
  ./configure
  make clean
  make check
  sudo make uninstall
  sudo make install
  popd
done
