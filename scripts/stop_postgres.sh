#!/bin/bash

pg_ctl \
    -D ./data/toc-machine-trading \
    -l ./data/toc-machine-trading/logfile \
    stop
