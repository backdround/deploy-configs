[![Go Reference](https://img.shields.io/badge/go-reference-%2300ADD8?style=flat-square)](https://pkg.go.dev/github.com/backdround/deploy-configs)
[![Tests](https://img.shields.io/github/actions/workflow/status/backdround/deploy-configs/tests.yml?branch=main&label=tests&style=flat-square)](https://github.com/backdround/deploy-configs/actions)
[![Codecov](https://img.shields.io/codecov/c/github/backdround/deploy-configs?style=flat-square)](https://app.codecov.io/gh/backdround/deploy-configs/)
[![Go Report](https://goreportcard.com/badge/github.com/backdround/deploy-configs?style=flat-square)](https://goreportcard.com/report/github.com/backdround/deploy-configs)

# Deploy-configs

It serves to deploy config files to your pc instance by `yaml` description.

It can:
- crate symlinks
- expand templates
- execute commnad

It shows execution log gracefully:
- what was changed
- what wasn't changed
- what wasn't able to succeed

### Example

Lets create sample dot repository to deploy:
```bash
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

`deploy-configs.yaml` contains yaml that drives deploying process:
```yaml
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

To deploy home instance we execute application:

```bash
deploy-configs home
```

Result tree with deployed configs:
```bash
/home/user/
├── .config
│   ├── flameshot
│   │   └── flameshot.ini
│   └── i3
│       └── config
├── configs
├── .gitconfig -> /home/user/configs/git/gitconfig
├── .tmux.conf -> /home/user/configs/terminal/tmux
└── configs
    └── ... # output truncated
```

### Deploy yaml format

#### Main `deploy-configs.yaml` structure

Schematic example:
```yaml
# Arbitrary data that you need in your instances
<any-shared-data>:
  any:
    - shared
    - data

# Instances contains a dictionary with all possible instances.
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

  laptop:
    commands:
      ...
```

#### Links field
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

#### Templates field
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

#### Commands field
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

#### Path replacement
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
