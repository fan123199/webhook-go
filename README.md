# webhook-go

a simple webhook implementation  by golang, now support gitee and github.


## Usage

1. create or locate a deploy script for your project
1. `go get` this project to your web server, and config database.toml. (config example is in the project) and run it.
1. in your github/gitee project setting page, enable webhook. setup url + port + password.
1. all it is done.
1. every time you change database.toml(except port), it take effect next time , no necessary to restart.

## Next

+ a web UI for project management

## idea

The deploy script should be anywhere, and this project should be able to
config the path of script.