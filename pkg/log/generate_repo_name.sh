git config --get remote.origin.url | sed 's/.*\///' | sed 's/.git$//' >repo_name
