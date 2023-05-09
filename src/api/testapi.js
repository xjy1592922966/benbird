// 业绩模块
import _axios from './_axios';
import { axiosUrl } from '@/config/env';

// 获取用户
export let users = (data) =>
	_axios({
		data: data,
		url: `${axiosUrl}users`,
		loading: true,
	});
// 业绩提成列表
export let home = (data) =>
	_axios({
		data: data,
		url: `${axiosUrl}`,
		loading: true,
	});

// 业绩提成列表
export let register = (data) =>
	_axios({
		data: data,
		url: `${axiosUrl}/register`,
		loading: true,
	});
