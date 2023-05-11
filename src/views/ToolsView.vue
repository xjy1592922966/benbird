<template>
	<a-menu :mode="mode" :theme="theme" v-if="menuData.length > 0">
		<template v-for="submenu in menuData" :key="submenu.id">
			<a-sub-menu v-if="checkPermission(submenu)">
				<template #title>
					<a-icon :type="submenu.icon" />
					<span>{{ submenu.name }}</span>
				</template>
				<template v-for="menu in submenu.children" :key="menu.id">
					<a-menu-item v-if="checkPermission(menu)">
						<router-link :to="menu.path">
							<a-icon :type="menu.icon" />
							<span>{{ menu.name }}</span>
						</router-link>
					</a-menu-item>
				</template>
			</a-sub-menu>
		</template>
	</a-menu>
</template>

<script>
import { computed, defineComponent, onMounted, ref } from 'vue';
import { useStore } from 'vuex';
// import { getMenuList } from '@/api/menu';
import { filterAsyncRouter } from '@/utils/permission';

export default defineComponent({
	components: {},
	setup() {
		const store = useStore();
		const mode = 'inline';
		const theme = 'dark';
		const menuData = ref([]);

		const roles = computed(() => store.getters.roles);

		const initMenu = async () => {
			// const response = await getMenuList();

			let response = {
				code: 200,
				data: [
					{
						id: 1,
						name: 'Dashboard',
						icon: 'dashboard',
						path: '/dashboard',
						component: 'Layout',
						redirect: '/dashboard/analysis',
						children: [
							{
								id: 2,
								name: 'Analysis',
								icon: '',
								path: 'analysis',
								component: 'dashboard/Analysis',
								meta: {
									title: 'Analysis',
									roles: ['admin', 'editor'],
								},
							},
							{
								id: 3,
								name: 'Monitor',
								icon: '',
								path: 'monitor',
								component: 'dashboard/Monitor',
								meta: {
									title: 'Monitor',
									roles: ['admin', 'editor'],
								},
							},
						],
					},
					{
						id: 4,
						name: 'System',
						icon: 'setting',
						path: '/system',
						component: 'Layout',
						redirect: '/system/user',
						children: [
							{
								id: 5,
								name: 'User',
								icon: '',
								path: 'user',
								component: 'system/User',
								meta: {
									title: 'User Management',
									roles: ['admin'],
								},
							},
							{
								id: 6,
								name: 'Role',
								icon: '',
								path: 'role',
								component: 'system/Role',
								meta: {
									title: 'Role Management',
									roles: ['admin'],
								},
							},
						],
					},
				],
			};

			const { data } = response;

			menuData.value = filterAsyncRouter(data);
			console.log('菜单信息menuData', menuData.value);
		};

		const checkPermission = (menu) => {
			console.log('menu', menu);
			if (menu.meta && menu.meta.roles) {
				return menu.meta.roles.includes(roles.value);
			} else {
				return true;
			}
		};

		onMounted(() => {
			initMenu();
		});

		return {
			mode,
			theme,
			menuData,
			checkPermission,
		};
	},
});
</script>
