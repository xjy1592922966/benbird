<template>
	<a-layout>
		<a-layout-header class="header">
			<div class="logo" />
			<a-menu
				v-model:selectedKeys="selectedKeys1"
				theme="dark"
				mode="horizontal"
				:style="{ lineHeight: '64px' }"
			>
				<a-menu-item key="1">
					<router-link to="/">首页</router-link>
				</a-menu-item>
				<a-menu-item key="2">
					<router-link to="/about">测试页面</router-link>
				</a-menu-item>
				<a-menu-item key="3">
					<router-link to="/tools">工具箱</router-link>
				</a-menu-item>
				<a-menu-item key="4">
					<router-link to="/md5conversiontool">MD5转换工具</router-link>
				</a-menu-item>
			</a-menu>
		</a-layout-header>
		<a-layout>
			<a-layout-sider width="200" style="background: #fff">
				<a-menu
					v-model:selectedKeys="selectedKeys2"
					v-model:openKeys="openKeys"
					mode="inline"
					:style="{ height: '100%', borderRight: 0 }"
				>
					<a-sub-menu key="sub1">
						<template #title>
							<span>
								<user-outlined />
								subnav 1
							</span>
						</template>
						<a-menu-item key="1">option1</a-menu-item>
						<a-menu-item key="2">option2</a-menu-item>
						<a-menu-item key="3">option3</a-menu-item>
						<a-menu-item key="4">option4</a-menu-item>
					</a-sub-menu>
					<a-sub-menu key="sub2">
						<template #title>
							<span>
								<laptop-outlined />
								subnav 2
							</span>
						</template>
						<a-menu-item key="5">option5</a-menu-item>
						<a-menu-item key="6">option6</a-menu-item>
						<a-menu-item key="7">option7</a-menu-item>
						<a-menu-item key="8">option8</a-menu-item>
					</a-sub-menu>
					<a-sub-menu key="sub3">
						<template #title>
							<span>
								<notification-outlined />
								subnav 3
							</span>
						</template>
						<a-menu-item key="9">option9</a-menu-item>
						<a-menu-item key="10">option10</a-menu-item>
						<a-menu-item key="11">option11</a-menu-item>
						<a-menu-item key="12">option12</a-menu-item>
					</a-sub-menu>
				</a-menu>
			</a-layout-sider>
			<a-layout style="padding: 0 24px 24px">
				<a-breadcrumb style="margin: 16px 0">
					<a-breadcrumb-item>Home</a-breadcrumb-item>
					<a-breadcrumb-item>List</a-breadcrumb-item>
					<a-breadcrumb-item>App</a-breadcrumb-item>
				</a-breadcrumb>
				<a-layout-content
					:style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }"
				>
					<router-view />
				</a-layout-content>
			</a-layout>
		</a-layout>
	</a-layout>
</template>

<script>
import { UserOutlined, LaptopOutlined, NotificationOutlined } from '@ant-design/icons-vue';
import { defineComponent, ref } from 'vue';
export default defineComponent({
	components: {
		UserOutlined,
		LaptopOutlined,
		NotificationOutlined,
	},
	setup() {
		return {
			selectedKeys1: ref(['2']),
			selectedKeys2: ref(['1']),
			collapsed: ref(false),
			openKeys: ref(['sub1']),
		};
	},
});
</script>
<style lang="scss">
#components-layout-demo-top-side-2 .logo {
	float: left;
	width: 120px;
	height: 31px;
	margin: 16px 24px 16px 0;
	background: rgba(255, 255, 255, 0.3);
}

.ant-row-rtl #components-layout-demo-top-side-2 .logo {
	float: right;
	margin: 16px 0 16px 24px;
}

.site-layout-background {
	background: #fff;
}
#app {
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	color: #2c3e50;
}

nav {
	padding: 30px;

	a {
		font-weight: bold;
		color: #2c3e50;

		&.router-link-exact-active {
			color: #42b983;
		}
	}
}
</style>
