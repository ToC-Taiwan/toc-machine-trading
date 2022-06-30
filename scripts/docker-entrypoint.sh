#!/bin/bash

export $(cat /toc-machine-trading/.env | xargs)

/toc-machine-trading/toc-machine-trading
