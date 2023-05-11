// 业绩模块
import _axios from './_axios';
import { axiosUrl } from '@/config/env';

// 获取菜单
export let getMenuList = (data) =>
	_axios({
		data: data,
		url: `${axiosUrl}getMenuList`,
		loading: true,
	});
