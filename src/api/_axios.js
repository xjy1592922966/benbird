import axios from 'axios';
import store from '@/store';
// 引入localStorage相关方法
import { getStorage, checkTime } from '@/utils/myUtils';

// 这里是转化跟踪
// import { _trackEventRequest } from '@/utils/statistics';

// 校验回调状态
function checkStatus(response) {
	// 如果http状态码正常，则直接返回数据
	if (response && (response.status === 200 || response.status === 304 || response.status === 400)) {
		// 如果不需要除了data之外的数据，可以直接 return response.data
		return response.data;
	}
	// 异常状态下，把错误信息返回去
	// ELEMENT.Message({
	// 	message: `请求异常-${response.status}`,
	// 	type: 'error',
	// 	duration: 4000,
	// });
	return false;
}

/**
 * axios 请求方法
 * */
let _axios = (options) => {
	/**
	 * @param {String} method 请求方式
	 * @param {String} url 请求地址
	 * @param {Object} data 请求参数
	 * @param {boolean} loading 是否开启loading
	 * @param {boolean} needToken 是否需要提交token
	 * @param {Object} checkTime 是否有中国时间，将其转换成年月日时分秒
	 */
	if (navigator && navigator.onLine === false) {
		if (!localStorage.getItem('netWorking')) {
			localStorage.setItem('netWorking', true);
			setTimeout(() => {
				localStorage.removeItem('netWorking');
			}, 1000);
			alert('当前无网络或网络不稳定，请更换网络后重新操作！');
		}
		return;
	}
	let { method = 'post', url, data = {}, loading = false, needToken = false } = options;

	data = checkTime(data);

	let _loading = {};
	// 加载loding
	if (loading !== false) {
		// 如果loading传递的是具体div的id，则显示在该id上
		// var queryId = typeof loading === 'string' ? loading : '#main-el-main';
		// _loading = ELEMENT.Loading.service({
		// 	target: document.querySelector(queryId),
		// 	lock: true,
		// });
	}
	let setHeaders = {};
	// 如果需要设置token，则设置初始值
	if (needToken === true) {
		setHeaders = { 'XX-token': '' };
	}
	// 发起请求
	return new Promise((resolve, reject) => {
		// _trackEventRequest('API请求', url);

		axios({
			method: method,
			url: url,
			data: data,
			headers: setHeaders,
		})
			.then((response) => {
				_loading.close ? _loading.close() : '';
				// debugger
				// 校验状态并获取实际数据
				let responseData = checkStatus(response);
				// code=0请求失败,提示错误信息
				if (responseData.code === 0) {
					// ELEMENT.Message({
					// 	message: responseData.msg,
					// 	type: 'error',
					// 	duration: 5000,
					// });
					if (options.callback) {
						options.callback(responseData);
					}
					reject();
				}
				// code=10001 登录失效，前往登录页面
				if (responseData.code === 10001 && !options.noramin) {
					//noramin 不提示弹框
					// ELEMENT.Message({
					// 	message: responseData.msg,
					// 	type: 'error',
					// });
					// const mobileList = localStorage.getItem('useLoginMobile');
					// setTimeout(function () {
					// 	localStorage.clear();
					// 	mobileList && localStorage.setItem('useLoginMobile', mobileList);
					// 	router.replace({ path: '/login' });
					// }, 2000);
					resolve(responseData);
				}
				// 如果有token设置token
				// if (responseData.data && responseData.data.token !== undefined) {
				//   store.commit("SET_TOKEN", responseData.data.token);
				// }
				resolve(responseData.data);
			})
			.catch((error) => {
				console.log(error);

				_loading.close ? _loading.close() : '';
			});
	});
};

// 请求拦截器
axios.interceptors.request.use(
	(config) => {
		// 如果有定义XX-token参数，则是需要上传token
		if (config.headers['XX-token'] !== undefined) {
			config.headers['XX-token'] = store.state.userModule.token;
		}
		// 如果有token值且没有过期，则每次请求更新过期时间
		if (getStorage('token') !== '') {
			store.commit('SET_TOKEN', getStorage('token'));
		}
		config.headers['XX-Device-Type'] = 'pc';
		return config;
	},
	(err) => {
		return Promise.reject(err);
	},
);

export default _axios;

/**
 * code
 * 0:失败
 * 1：成功
 * 10001：登录失效
 *
 */
