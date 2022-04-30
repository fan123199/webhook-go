# webhook-go

a simple webhook implementation  by golang, now support gitee and github.

## Status

The project is deprecated, it is hard to use.

## Usage

1. create or locate a deploy script for your project
1. `go get` this project to your web server, and config database.toml. (config example is in the project) and run it.
1. in your github/gitee project setting page, enable webhook. setup url + port + password.
1. all it is done.
1. every time you change database.toml(except port), it take effect next time , no necessary to restart.


## Idea

The deploy script should be anywhere, and this project should be able to
config the path of script.
