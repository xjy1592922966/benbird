/**
 * 根据角色过滤菜单和页面元素
 * @param {Array} asyncRouterMap 动态路由
 * @param {String} roles 用户角色
 */
export function filterAsyncRouter(asyncRouterMap, roles) {
	const accessedRouters = asyncRouterMap.filter((route) => {
		if (hasPermission(roles, route)) {
			if (route.children && route.children.length) {
				route.children = filterAsyncRouter(route.children, roles);
			}
			return true;
		}
		return false;
	});
	return accessedRouters;
}

/**
 * 判断角色是否有权限访问该菜单或页面元素
 * @param {String} roles 用户角色
 * @param {Object} route 菜单或页面元素
 */
function hasPermission(roles, route) {
	if (route.meta && route.meta.roles) {
		return route.meta.roles.includes(roles);
	} else {
		return true;
	}
}
