# asana
A CLI utility for creating tasks in a Asana.

# Build
```sh
go install github.com/michael-lorenzo/asana@latest
asana -generate
```

# Usage
Edit config file `asana.toml` use https://developers.asana.com/docs/quick-start if you need guidance
```sh
asana "task name"
```
or
```sh
asana "task name" "task notes"
```
Task should be under "My tasks"
