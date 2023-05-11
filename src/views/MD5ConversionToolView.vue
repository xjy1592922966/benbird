<template>
	<a-form
		:model="formState"
		v-bind="layout"
		name="nest-messages"
		:validate-messages="validateMessages"
	>
		<a-form-item :name="['md5ConversionTools', 'originalText']" label="原文">
			<a-textarea
				:rows="4"
				placeholder="输入需要加密的内容，一行一个例如：
13632448075
asdasdasd"
				@change="changeValue"
				v-model:value="formState.md5ConversionTools.originalText"
			/>
		</a-form-item>

		<a-divider></a-divider>

		<a-form-item :name="['md5ConversionTools', 'ciphertext']" label="密文">
			<a-textarea
				:rows="4"
				placeholder=""
				v-model:value="formState.md5ConversionTools.ciphertext"
			/>
		</a-form-item>
	</a-form>
</template>

<script>
import { defineComponent, reactive } from 'vue';
export default defineComponent({
	setup() {
		const layout = {
			labelCol: {
				span: 8,
			},
			wrapperCol: {
				span: 16,
			},
		};
		const validateMessages = {
			required: '${label} is required!',
			types: {
				email: '${label} is not a valid email!',
				number: '${label} is not a valid number!',
			},
			number: {
				range: '${label} must be between ${min} and ${max}',
			},
		};
		const formState = reactive({
			md5ConversionTools: {
				originalText: '',
				ciphertext: '',
			},
		});
		const changeValue = (values) => {
			console.log('changeValue:', values);
			console.log('formState.originalText:', formState.originalText);
		};
		return {
			formState,
			layout,
			changeValue,
			validateMessages,
		};
	},
});
</script>

<style lang="scss" scoped>
.home {
	text-align: center;
}
</style>
