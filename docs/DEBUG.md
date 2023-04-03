# DEBUG

## Launch json and Tasks json

```bash
mkdir -p .vscode
echo '{
    "version": "0.2.0",
    "configurations": [
        {
            "preLaunchTask": "Build",
            "name": "Debug",
            "debugAdapter": "dlv-dap",
            "type": "go",
            "request": "launch",
            "mode": "exec",
            "program": "${workspaceFolder}/toc-machine-trading",
            "envFile": "${workspaceFolder}/.env",
        }
    ]
}' > .vscode/launch.json

echo '{
    "version": "2.0.0",
    "cwd": "${workspaceFolder}",
    "tasks": [
        {
            "label": "go generate",
            "type": "shell",
            "command": "go",
            "args": [
                "generate",
                "./..."
            ],
        },
        {
            "label": "cp config",
            "type": "shell",
            "command": "cp",
            "args": [
                "./configs/default.config.yml",
                "./configs/config.yml"
            ],
        },
        {
            "label": "go build",
            "type": "shell",
            "command": "go",
            "args": [
                "build",
                "-o",
                "toc-machine-trading",
                "-gcflags=all=\"-N -l\"",
                "./cmd/app",
            ],
        },
        {
            "label": "Build",
            "dependsOrder": "sequence",
            "dependsOn": [
                "go generate",
                "cp config",
                "go build"
            ]
        }
    ]
}' > .vscode/tasks.json
```
