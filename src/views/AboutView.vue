<template>
	<div class="about">
		<h1>{{ h1Text }}</h1>
		<h2>{{ h2Text }}</h2>

		<a-button type="primary" @click="click">click1</a-button>
		<a-button type="primary" @click="click2">click2</a-button>
		<a-button type="primary" @click="click3">{{ testBlockName }}</a-button>
		<a-button type="primary" @click="click4">请求首页</a-button>

		<a-form
			:model="formState"
			name="basic"
			:label-col="{ span: 8 }"
			:wrapper-col="{ span: 16 }"
			autocomplete="off"
			@finish="onFinish"
			@finishFailed="onFinishFailed"
		>
			<a-form-item
				label="Username"
				name="username"
				:rules="[{ required: true, message: 'Please input your username!' }]"
			>
				<a-input v-model:value="formState.username" />
			</a-form-item>

			<a-form-item
				label="Password"
				name="password"
				:rules="[{ required: true, message: 'Please input your password!' }]"
			>
				<a-input-password v-model:value="formState.password" />
			</a-form-item>

			<a-form-item name="remember" :wrapper-col="{ offset: 8, span: 16 }">
				<a-checkbox v-model:checked="formState.remember">Remember me</a-checkbox>
			</a-form-item>

			<a-form-item :wrapper-col="{ offset: 8, span: 16 }">
				<a-button type="primary" @click="submitForm">Submit</a-button>
			</a-form-item>
		</a-form>
	</div>
</template>

<script>
import { defineComponent, ref, reactive } from 'vue';
import * as TESTAPI from '@/api/testapi';
export default defineComponent({
	setup() {
		const formState = reactive({
			username: '',
			password: '',
			remember: true,
		});
		const onFinish = (values) => {
			console.log('Success:', values);
		};
		const onFinishFailed = (errorInfo) => {
			console.log('Failed:', errorInfo);
		};

		const textArr = {
			h1Text: ref('这里是标题1'),
			h2Text: ref('这里是标题2'),
		};

		const functionArr = {
			click2: () => {
				console.log('点击了按钮2');
			},
			click4: () => {
				console.log('请求首页');
				TESTAPI.home().then((apiResult) => {
					console.log('apiResult请求成功', apiResult);
				});
			},
			submitForm: () => {
				console.log('注册');
				TESTAPI.register({
					...formState,
				}).then((apiResult) => {
					console.log('注册请求成功apiResult', apiResult);
				});
			},
		};

		const testBlock1 = {
			testBlockName: ref(1),
			click3: () => {
				testBlock1.testBlockName.value += 1;
				console.log('点击了按钮2');

				TESTAPI.users().then((apiResult) => {
					console.log('apiResult请求成功', apiResult);
				});
			},
		};

		return {
			...textArr,
			...functionArr,
			...testBlock1,
			formState,
			onFinish,
			onFinishFailed,
		};
	},
});
</script>
