## 目录结构

在使用 Gin 框架构建 API 应用时，一般可以按照以下目录结构进行组织：

```
├── config
│   ├── config.go
│   └── settings.go
├── controller
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── ...
├── middleware
│   ├── auth_middleware.go
│   ├── logger_middleware.go
│   └── ...
├── model
│   ├── auth_model.go
│   ├── user_model.go
│   └── ...
├── router
│   ├── api_router.go
│   └── router.go
├── service
│   ├── auth_service.go
│   ├── user_service.go
│   └── ...
├── utils
│   ├── jwt.go
│   ├── response.go
│   └── ...
├── main.go
```

其中，各个目录的作用如下：

- `config`：存放应用程序的配置文件和相关设置，比如数据库连接信息、端口号、密钥等。
- `controller`：存放控制器，处理 API 请求，调用相应的服务层方法，对数据进行处理和返回响应结果。
- `middleware`：存放中间件，用于增强应用程序的功能，比如日志中间件、认证中间件、缓存中间件等。
- `model`：存放模型，主要处理数据的存储、访问和处理，包括数据库相关的操作。
- `router`：存放路由相关的代码，包括 API 路由和 Web 路由等。
- `service`：存放服务层代码，主要负责处理业务逻辑，调用 Model 层的方法进行数据的处理和返回结果。
- `utils`：存放一些公共的工具类，比如 JWT 的生成和解析，响应数据的封装等。
- `main.go`：应用程序入口文件，主要用于启动应用程序。

以上是一种常见的目录结构，但具体的目录结构还要根据实际情况进行调整和优化。同时，还需要注意代码的规范和文档的编写，以便于代码的维护和协作。

这个目录结构并没有一个官方的参考文档，但是这个目录结构是根据 Gin 框架和 Go 语言的开发规范、经验和实际项目经验整合而来的。在实际项目中，这个目录结构已经被广泛采用，也被认为是一种比较成熟的目录结构设计。

此外，针对这个目录结构，也有一些开源项目和文档提供了一些参考和实践经验，可以帮助你更好地理解和应用这个目录结构，例如：

- [Gin Web App Example](https://github.com/go-programming-tour-book/blog-service)
- [RESTful API With Gin](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gin)
- [Golang Project Layout](https://github.com/golang-standards/project-layout)

这些资源可以帮助你更好地理解和使用这个目录结构，同时也可以帮助你在实际项目中更好地进行代码组织和开发。

# 多级菜单的实现

为了实现多级菜单，可以设计以下表结构：

1. 表名：menu
   字段列表：

- id：菜单 ID，类型为 INT，主键，自增
- parent_id：父菜单 ID，类型为 INT，外键关联 menu 表的 id 字段
- name：菜单名称，类型为 VARCHAR
- icon：菜单图标，类型为 VARCHAR
- path：菜单路径，类型为 VARCHAR
- component：菜单组件，类型为 VARCHAR
- redirect：菜单重定向路径，类型为 VARCHAR
- meta_title：菜单元数据标题，类型为 VARCHAR
- meta_roles：菜单元数据角色，类型为 VARCHAR
- version：菜单版本，类型为 INT

2. 表名：role
   字段列表：

- id：角色 ID，类型为 INT，主键，自增
- name：角色名称，类型为 VARCHAR

3. 表名：menu_role
   字段列表：

- menu_id：菜单 ID，类型为 INT，外键关联 menu 表的 id 字段
- role_id：角色 ID，类型为 INT，外键关联 role 表的 id 字段

通过以上表结构，可以实现多级菜单的存储，并支持角色控制和版本控制。

针对 go 语言 gin 框架的增删改查接口，可以设计如下：

1. GET /menu
   查询所有菜单。

2. GET /menu/:id
   查询指定菜单。

3. POST /menu
   创建菜单。

4. PUT /menu/:id
   更新菜单。

5. DELETE /menu/:id
   删除菜单。

6. GET /role
   查询所有角色。

7. GET /role/:id
   查询指定角色。

8. POST /role
   创建角色。

9. PUT /role/:id
   更新角色。

10. DELETE /role/:id
    删除角色。

11. GET /menu_role/:menu_id
    查询指定菜单的所有角色。

12. POST /menu_role
    创建菜单角色关联。

13. DELETE /menu_role/:menu_id/:role_id
    删除菜单角色关联。

以上接口可以满足多级菜单的增删改查需求，并支持角色控制和版本控制。

# 请按照以上方案，给出 go gin 框架代码。
