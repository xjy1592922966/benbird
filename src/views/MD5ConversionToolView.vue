<template>
	<a-form
		:model="formState"
		v-bind="layout"
		name="nest-messages"
		:validate-messages="validateMessages"
	>
		<a-row :gutter="16">
			<a-col class="gutter-row" :span="11">
				<div class="gutter-box">
					<a-form-item :name="['md5ConversionTools', 'originalText']" label="原文">
						<a-textarea
							:rows="4"
							placeholder="输入需要加密的内容，一行一个例如：
13632448075
asdasdasd"
							@change="funChangeValue"
							v-model:value="formState.md5ConversionTools.originalText"
						/>
					</a-form-item>
				</div>
			</a-col>
			<a-col class="gutter-row" :span="2">
				<div class="gutter-box">
					<a-form-item>
						<a-button type="primary" shape="round" :size="large" @click="funEncryptionConversion">
							转换
						</a-button>
					</a-form-item>
				</div>
			</a-col>
			<a-col class="gutter-row" :span="11">
				<div class="gutter-box">
					<a-form-item :name="['md5ConversionTools', 'ciphertext']" label="密文">
						<a-textarea
							:rows="4"
							placeholder=""
							v-model:value="formState.md5ConversionTools.ciphertext"
						/>
					</a-form-item>
				</div>
			</a-col>
		</a-row>
	</a-form>
</template>

<script>
import { defineComponent, ref } from 'vue';
import CryptoJS from 'crypto-js';
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
		const formState = ref({
			md5ConversionTools: {
				originalText: '',
				ciphertext: '',
			},
		});

		const funEncryptionConversion = () => {
			console.log('转换按钮');

			console.log(
				'formState.originalText:',
				formState.value.md5ConversionTools.originalText.split('\n'),
			);

			let originalTextArr = formState.value.md5ConversionTools.originalText.split('\n');

			let ciphertextArr = originalTextArr.map((item) => {
				return CryptoJS.MD5(item).toString();
			});

			formState.value.md5ConversionTools.ciphertext = ciphertextArr.join('\n');
		};

		const funChangeValue = () => {
			// console.log('changeValue:', values);
		};
		return {
			formState,
			layout,
			funChangeValue,
			funEncryptionConversion,
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
