# webhook-go

simple implement webhook, now support gitee and github.

## next

implement github webhook

## idea

是将部署的脚本放在对应项目里，还是放在本项目中？

如果是对应项目，，。优点: 那么解耦性比较好，可以在对应项目执行部署。缺点是脚本权限问题，以及项目代码混入了
奇怪的东西。

如果是本项目，缺点: 每添加一个项目，需要在这里新增一个部署脚本。和增加中的配置。
优点是可以统一管理全部的项目部署。
