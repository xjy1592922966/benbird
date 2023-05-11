import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';

const routes = [
	{
		path: '/',
		name: '首页',
		component: HomeView,
		children: [
			// {
			// 	path: '/dashboard',
			// 	name: 'Dashboard',
			// 	component: Dashboard,
			// 	meta: {
			// 		title: '首页',
			// 		icon: 'dashboard',
			// 	},
			// },
			// {
			// 	path: '/user',
			// 	name: 'UserList',
			// 	component: UserList,
			// 	meta: {
			// 		title: '用户管理',
			// 		icon: 'user',
			// 		roles: ['admin'], // 只有admin角色可以访问
			// 	},
			// },
			// {
			// 	path: '/role',
			// 	name: 'RoleList',
			// 	component: RoleList,
			// 	meta: {
			// 		title: '角色管理',
			// 		icon: 'team',
			// 		roles: ['admin'], // 只有admin角色可以访问
			// 	},
			// },
			// {
			// 	path: '/permission',
			// 	name: 'PermissionList',
			// 	component: PermissionList,
			// 	meta: {
			// 		title: '权限管理',
			// 		icon: 'lock',
			// 		roles: ['admin'], // 只有admin角色可以访问
			// 	},
			// },
		],
	},
	{
		path: '/about',
		name: '关于',
		// route level code-splitting
		// this generates a separate chunk (about.[hash].js) for this route
		// which is lazy-loaded when the route is visited.
		component: () => import(/* webpackChunkName: "about" */ '../views/AboutView.vue'),
	},
	{
		path: '/tools',
		name: '工具箱',
		// route level code-splitting
		// this generates a separate chunk (about.[hash].js) for this route
		// which is lazy-loaded when the route is visited.
		component: () => import(/* webpackChunkName: "about" */ '../views/ToolsView.vue'),
	},
	{
		path: '/md5conversiontool',
		name: 'MD5转换工具',
		// route level code-splitting
		// this generates a separate chunk (about.[hash].js) for this route
		// which is lazy-loaded when the route is visited.
		component: () => import(/* webpackChunkName: "about" */ '../views/MD5ConversionToolView.vue'),
	},
];

const router = createRouter({
	history: createWebHistory(process.env.BASE_URL),
	routes,
});

export default router;
