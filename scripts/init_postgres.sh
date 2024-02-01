#!/bin/bash
set -e

echo "Trying to stop postgresql server..."
pg_ctl -D ./data/toc-machine-trading -l ./data/toc-machine-trading/logfile stop >/dev/null 2>&1 || true

if [ -d ./data/toc-machine-trading ]; then
    echo "Database already initialized. Remove ./data/toc-machine-trading if you want to reinitialize it."
    read -p "Do you want to reinitialize the database? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Exiting..."
        exit 0
    fi
    rm -rf ./data/toc-machine-trading
    mkdir -p ./data/toc-machine-trading
fi

echo "Initializing database..."

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

pg_ctl -D ./data/toc-machine-trading -l ./data/toc-machine-trading/logfile stop
