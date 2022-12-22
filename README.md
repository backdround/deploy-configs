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
        command: "sed \"s~%HOMEDIR%~$HOME~g\" {{.Input}} > {{.Output}}"

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
