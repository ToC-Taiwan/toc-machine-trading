# DEBUG

## Launch json and Tasks json

```bash
mkdir -p .vscode
echo '{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug",
            "type": "go",
            "request": "launch",
            "debugAdapter": "dlv-dap",
            "mode": "exec",
            "program": "${workspaceFolder}/toc-machine-trading",
            "preLaunchTask": "Build",
            "console": "integratedTerminal",
            "internalConsoleOptions": "neverOpen"
        }
    ]
}' > .vscode/launch.json

echo '{
    "version": "2.0.0",
    "cwd": "${workspaceFolder}",
    "type": "shell",
    "presentation": {
        "close": true
    },
    "tasks": [
        {
            "label": "go generate",
            "command": "go",
            "args": [
                "generate",
                "./..."
            ],
        },
        {
            "label": "swag",
            "command": "bash",
            "args": [
                "./scripts/generate_swagger.sh",
            ],
        },
        {
            "label": "cp config",
            "command": "cp",
            "args": [
                "./configs/default.config.yml",
                "./configs/config.yml"
            ],
        },
        {
            "label": "go build",
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
                "swag",
                "cp config",
                "go build"
            ]
        }
    ]
}' > .vscode/tasks.json
```

## Postman

- add set token in collection

```js
// jsonData will contain the valid JSON object from the response body
let jsonData = pm.response.json();

// Variables and JSON Keys are case sensitive!
pm.collectionVariables.set("apiKey", jsonData.token);
```

- set baseUrl `https://trader.tocraw.com/tmt`
