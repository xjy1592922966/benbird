# 随记 - 前端

1、接口请求，需要配置

2、@相对路径配置

# 随机 - 后端

1、需要使用 go 框架出一套菜单栏

# 【大模块】

1.这个项目需要开始写脚手架

在使用 antd vue 开发 Vue 中后台项目时，一般分为以下板块和模块：

1. 登录和身份验证模块：包括用户登录、注册、找回密码等功能，以及身份验证和权限控制等。

2. 主界面和导航模块：包括主界面布局、导航菜单、面包屑导航等。

3. 数据展示模块：包括表格、图表、图形等展示数据的组件和页面。

4. 数据操作模块：包括表单、弹窗、对话框等操作数据的组件和页面。

5. 通知和消息模块：包括消息提醒、通知中心、消息管理等。

6. 用户管理模块：包括用户列表、用户详情、用户编辑等。

7. 角色管理模块：包括角色列表、角色详情、角色编辑等。

8. 权限管理模块：包括权限列表、权限详情、权限编辑等。

9. 系统设置模块：包括系统配置、日志管理、数据备份等。

10. 帮助和支持模块：包括帮助中心、联系我们、常见问题等。

以上是一些常见的板块和模块，根据具体项目需求，可以进行适当调整和拓展。

如果是前台项目，需要的模块可能会有所不同，一般来说，前台项目需要的模块包括：

1. 首页模块：包括首页布局、轮播图、推荐商品等。

2. 商品列表模块：包括商品分类、商品列表、筛选、排序等。

3. 商品详情模块：包括商品详情、商品评价、商品图片等。

4. 购物车模块：包括购物车列表、结算、优惠券、地址管理等。

5. 订单模块：包括订单列表、订单详情、订单状态等。

6. 支付模块：包括支付方式、支付流程、支付状态等。

7. 用户中心模块：包括个人资料、订单列表、收货地址等。

8. 消息中心模块：包括消息提醒、站内信、通知等。

9. 帮助中心模块：包括常见问题、在线客服、联系我们等。

10. 其他模块：如搜索、推荐、热门活动等。

以上是一些常见的前台模块，具体根据项目需求和业务场景进行适当调整和拓展。

# 面包屑优化的菜单数据结构

[
{
id: 'dashboard',
name: 'Dashboard',
icon: 'dashboard',
redirect: '/dashboard/analysis',
children: [
{
id: 'analysis',
name: 'Analysis',
path: '/dashboard/analysis',
component: 'dashboard/Analysis',
},
{
id: 'monitor',
name: 'Monitor',
path: '/dashboard/monitor',
component: 'dashboard/Monitor',
},
],
},
{
id: 'system',
name: 'System',
icon: 'setting',
redirect: '/system/user',
children: [
{
id: 'user',
name: 'User Management',
path: '/system/user',
component: 'system/User',
roles: ['admin'],
},
{
id: 'role',
name: 'Role Management',
path: '/system/role',
component: 'system/Role',
roles: ['admin'],
},
],
},
]
