#!/bin/bash

EXENAME=gcrond
LDATE=`date +%Y%m%d%H%M%S`
RUNDIR=/usr/local/share/gcrond

$RUNDIR/$EXENAME \
        1>/dev/null 2>$RUNDIR/log/paniclog.$LDATE  &
