# intro

gitlab ci requires some env variables, for diffent projects these env vars may be same. so we have this cmd tool -- `gitlab-vars`

# install

simplely download the binary and placed it in `$PATH`. the file name, whatever you wish !

# config

> default config file: `$HOME/.gitlab-vars.yaml`

`gitlab.private_token` supports three level:

- project level
- group level
- personal level

> Note: access token is **`sensitive`**

> vars refers: [project variables](https://docs.gitlab.com/ee/api/project_level_variables.html#list-project-variables)

```yaml
gitlab:
  private_token: <access token>
  server: <https://gitlab.xxx.com>

vars:
  - key: 'test_key' # required
    value: 'test_value' # required
    variable_type: 'env_var' # optional
    protected: true # optional
    masked: true # optional
    environment_scope: '*' # optional
```

# usage

```shell
# create vars for project
# check vars here:
# settings -> ci_cd -> Variables

$ ./gitlab-vars -p <projectId> -a create
```

```shell
# as for help

$ ./gitlab-vars -h
```

# build

> `Makefile` tested with platform `darwin` `amd64`

```shell
# build for windows
$ make win

# build for linux
$ make linux

# build for mac
$ make mac

# build for current platform
$ make build
```

# devlop

> cmd framework [cobra](https://github.com/spf13/cobra)

> gitlab api refered [plouc/go-gitlab-client](https://github.com/plouc/go-gitlab-client.git)
