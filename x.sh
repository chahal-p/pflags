#!/bin/bash

/data-disk/pflags/pflags --name "$0" parse ---- -s a -l abc -t string -r --default d1 --regex '.*' --default d2 -h 'ABC\nxyz' -- -s x -t number -- -s p -t bool ---- "$@"
