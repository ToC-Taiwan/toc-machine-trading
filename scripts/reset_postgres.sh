#!/bin/bash

pg_ctl -D ./data/toc-machine-trading -l ./data/toc-machine-trading/logfile stop

rm -rf ./data/toc-machine-trading
mkdir -p ./data/toc-machine-trading

initdb ./data/toc-machine-trading

gsed -i "$ a host    all    all    0.0.0.0/0    trust" ./data/toc-machine-trading/pg_hba.conf
gsed -i "$ a listen_addresses = '*'" ./data/toc-machine-trading/postgresql.conf

pg_ctl -D ./data/toc-machine-trading -l ./data/toc-machine-trading/logfile start

echo "\du
CREATE ROLE postgres WITH LOGIN PASSWORD 'password';
ALTER USER postgres WITH SUPERUSER;
\du" >sql_script

psql postgres -f sql_script
rm -rf sql_script

# pg_ctl -D ./data/toc-machine-trading -l ./data/toc-machine-trading/logfile stop
