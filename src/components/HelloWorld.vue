<template>
	<div class="hello">
		<h1>事儿是永远做不完的，一定要定好目标</h1>
		<p>
			这个项目的意义，就是熟悉vue3，并且能有用
			<br />
			vue3要学的东西，直接按照文档来就行。 可以写个小工具。
		</p>

		<h3>功能</h3>
		<ul>
			<li>MD5转换:熟悉UI界面和VUE3</li>
			<li>手机转sha256：熟悉vue前端架构</li>
			<li>淘客转唤起链接：熟悉vue前端架构</li>
			<li>注册&登录&鉴权&角色控制&菜单控制：熟悉go语言+前后端配合</li>
			<li>密码加密：熟悉go语言+前后端配合+安全</li>
			<li>换微信系统：熟悉前后端配合，后端+oss</li>
			<li>换微信系统+数据统计： 打通所有技术栈</li>
		</ul>

		<h3>技术栈</h3>
		<ul>
			<li>界面：antd vue +vuex</li>
			<li>后端：gin</li>
			<li>日志:zap库</li>
			<li>缓存：redis</li>
			<li>数据库：暂时sqlite</li>
		</ul>

		<h3>下一步计划</h3>
		<ul>
			<li>数据透视分析-搜索：elk方案</li>
			<li>数据透视分析-分布式计算：Hadoop</li>
		</ul>
	</div>
</template>

<script>
export default {
	name: 'HelloWorld',
	props: {
		msg: String,
	},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.hello {
	text-align: center;
}
h3 {
	margin: 40px 0 0;
}

ul {
	list-style-type: none;
	padding: 0;
}

li {
	display: inline-block;
	margin: 0 10px;
}

a {
	color: #42b983;
}
</style>
