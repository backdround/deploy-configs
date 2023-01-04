# Deploy-configs

[![Go Reference](https://img.shields.io/badge/go-reference-%2300ADD8?style=flat-square)](https://pkg.go.dev/github.com/backdround/deploy-configs)
[![Tests](https://img.shields.io/github/actions/workflow/status/backdround/deploy-configs/tests.yml?branch=main&label=tests&style=flat-square)](https://github.com/backdround/deploy-configs/actions)
[![Codecov](https://img.shields.io/codecov/c/github/backdround/deploy-configs?style=flat-square)](https://app.codecov.io/gh/backdround/deploy-configs/)
[![Go Report](https://goreportcard.com/badge/github.com/backdround/deploy-configs?style=flat-square)](https://goreportcard.com/report/github.com/backdround/deploy-configs)

It serves to deploy config files to pc by `yaml` description.

It can:
- create symlinks
- expand templates
- execute commands

It shows execution log gracefully:
- what was changed
- what wasn't changed
- what wasn't able to succeed



---
## Installation

By `go` compiler tools:
```bash
go install github.com/backdround/deploy-configs@main
```



---
## Example

Lets create sample config repository to deploy:
```bash
./configs
├── .git
├── deploy-configs.yaml
└── terminal
    ├── tmux
    └── gitconfig
```

`deploy-configs.yaml` contains yaml that drives deploying process:
```yaml
instances:
  home:
    links:
      tmux:
        target: "{{.GitRoot}}/terminal/tmux"
        link: "{{.Home}}/.tmux.conf"
      git:
        target: "{{.GitRoot}}/terminal/gitconfig"
        link: "{{.Home}}/.gitconfig"
```

To deploy `home` instance we execute application:

```bash
deploy-configs home
```

Result user home tree with deployed configs:
```bash
/home/user/
├── .gitconfig -> /home/user/configs/terminal/gitconfig
├── .tmux.conf -> /home/user/configs/terminal/tmux
└── configs
    └── ... # output truncated
```

## Complex example
<details><br>


```bash
# Config repository to deploy
./configs
├── .git
├── deploy-configs.yaml
├── desktop
│   ├── flameshot.ini
│   └── i3_template
├── git
│   └── gitconfig
└── terminal
    └── tmux
```

```yaml
# deploy-configs.yaml
instances:
  home:

    links:
      tmux:
        target: "{{.GitRoot}}/terminal/tmux"
        link: "{{.Home}}/.tmux.conf"
      git:
        target: "{{.GitRoot}}/git/gitconfig"
        link: "{{.Home}}/.gitconfig"

    commands:
      flameshot:
        input: "{{.GitRoot}}/desktop/flameshot.ini"
        output: "{{.Home}}/.config/flameshot/flameshot.ini"
        command: "sed \"s~%HOMEDIR%~$HOME~g\" '{{.Input}}' > '{{.Output}}'"

    templates:
      i3:
        input: "{{.GitRoot}}/desktop/i3_template"
        output: "{{.Home}}/.config/i3/config"
        data:
          telegram:
            size: "525 700"
            position: "1348 96"
          monitors:
            left: "DP-2"
            right: "HDMI-3"
```

```bash
# Deploying `home` instance
deploy-configs home
```

```bash
# Result home tree with deployed configs
/home/user/
├── .config
│   ├── flameshot
│   │   └── flameshot.ini
│   └── i3
│       └── config
├── .gitconfig -> /home/user/configs/git/gitconfig
├── .tmux.conf -> /home/user/configs/terminal/tmux
└── configs
    └── ... # output truncated
```

</details>



---
## Yaml format

Application deploys config files in accordance with `deploy-configs.yaml`. It
searches file recursively from current directory to root.

<details>
<summary> Main structure </summary><br>

Schematic example:
```yaml
# Arbitrary data that you need in your instances
<any-shared-data>:
  any:
    - shared
    - data

# Field contains a dictionary with all possible instances.
instances:
  # Instance is a set of deploying operation for performing at once.
  <instance-one>:
    [links:]
    [templates:]
    [commands:]

  <instance-two>:
    [links:]
    [templates:]
    [commands:]

  ...
```
Real example:
```yaml
.dev-links: &dev-links
  tmux:
    target: "{{.GitRoot}}/terminal/tmux"
    link: "{{.Home}}/.tmux.conf"
  zsh:
    target: "{{.GitRoot}}/terminal/zshrc"
    link: "{{.Home}}/.zshrc"

instances:
  home:
    links:
      <<: *dev-links
  work:
    links:
      <<: *dev-links
    templates:
      ...
```

</details>

---

<details>
<summary> Links </summary><br>

Links field describes links that are needed to be created.

Ripped out example:
```yaml
links:
  # Name is used in logs.
  tmux:
    # Target is a destination for link.
    target: "{{.GitRoot}}/terminal/tmux"
    # Link is used as a path to link creation.
    link: "{{.Home}}/.tmux.conf"
  zsh:
    target: "{{.GitRoot}}/terminal/zshrc"
    link: "{{.Home}}/.zshrc"
```

</details>

---

<details>
<summary> Templates </summary><br>

Templates field describes templates that are needed to be expanded and deployed.

Ripped out example:
```yaml
templates:
  # Name is used in logs.
  i3:
    # Input is a path to a `go` template (text/template).
    input: "{{.GitRoot}}/desktop/i3_template"
    # Output is a path to an expanded template.
    output: "{{.Home}}/.config/i3/config"
    # Data is an arbitrary structured data for template expantion.
    data:
      telegram:
        size: "525 700"
        position: "1348 96"
      monitors:
        left: "DP-2"
        right: "HDMI-3"
```

</details>

---

<details>
<summary> Commands </summary><br>

Commands field describes commands that create `output` files after execution.

Ripped out example:
```yaml
commands:
  # Name is used in logs.
  flameshot:
    # Input is a path to a command source config.
    input: "{{.GitRoot}}/desktop/flameshot.ini"
    # Output is a path to a generated config.
    output: "{{.Home}}/.config/flameshot/flameshot.ini"
    # Command converts the `input` config to an `output` config.
    # It allows {{.Input}} and {{.Output}} substitutions accordingly.
    command: "sed \"s~%HOMEDIR%~$HOME~g\" '{{.Input}}' > '{{.Output}}'"
```

</details>





---
## Path replacement
There are some replacements to define paths:
- {{.GitRoot}} - expands into current git directory.
- {{.Home}} - expands into current user home directory.

It expands only in specific fields which are used for path holding.

Example:
```yaml
links:
  git:
    target: "{{.GitRoot}}/git/gitconfig"
    link: "{{.Home}}/.gitconfig"
```

It will expand to:

```yaml
links:
  git:
    target: "/home/user/configs/git/gitconfig"
    link: "/home/user/.gitconfig"
```
