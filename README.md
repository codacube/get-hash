# get-hash

A simple project to get the sha256 of a file. Uses bubbletea as a tui.

Dockerfile included to run delve for debugging (or drop into a shell prompt)

## launch.json - remote debugging using delve (for bubbletea)

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to server",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "cwd": "${workspaceFolder}",
      "port": 2345,
      "host": "127.0.0.1",
      "remotePath": "/app"
    }
  ]
}
```
