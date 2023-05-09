import { ajaxUrl } from '@/config/env';
/**
 * 工具js集合
 */

// 密码验证
export function valPassword(val) {
	// 密码至少包含大写字母、小写字母、数字、长度为8-20位、且不能包含空格，请修改密码
	let reg = /^(?=.*?[A-Z])(?=(.*[a-z]){1,})(?=(.*[\d]){1,}).{8,20}$/;
	//空格
	let reg2 = /^[\S]*$/;

	if (!reg.test(val)) {
		//不符合
		return true;
	} else if (!reg2.test(val)) {
		//不符合
		return true;
	}
	return false;
}

/**
 * 时间戳 转换时间
 * @param {number} date:13位unix
 * @param {string} fmt:转换的格式 ["Y","M"]
 * */
export function formatDate(val, fmt = ['Y', 'M', 'D', 'h', 'm', 's']) {
	if (!val) {
		return '';
	}
	let date = new Date(parseInt(val) * 1000);
	let obj = {};
	let back = '';
	obj.Y = date.getFullYear() + '-';
	obj.M = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) + '-' : date.getMonth() + 1 + '-';
	obj.D = date.getDate() < 10 ? '0' + date.getDate() + ' ' : date.getDate() + ' ';
	obj.h = date.getHours() < 10 ? '0' + date.getHours() : date.getHours();
	obj.m = date.getMinutes() < 10 ? ':0' + date.getMinutes() : ':' + date.getMinutes();
	obj.s = date.getSeconds() < 10 ? ':0' + date.getSeconds() : ':' + date.getSeconds();
	for (let item of fmt) {
		back = back + obj[item];
	}
	return back;
}
/**
 * 时间格式化 但该项目没用到
 * */
export function dateFormatter(formatter, date) {
	date = date ? new Date(date) : new Date();
	const Y = date.getFullYear() + '',
		M = date.getMonth() + 1,
		D = date.getDate(),
		H = date.getHours(),
		m = date.getMinutes(),
		s = date.getSeconds();
	return formatter
		.replace(/YYYY|yyyy/g, Y)
		.replace(/YY|yy/g, Y.substr(2, 2))
		.replace(/MM/g, (M < 10 ? '0' : '') + M)
		.replace(/DD/g, (D < 10 ? '0' : '') + D)
		.replace(/HH|hh/g, (H < 10 ? '0' : '') + H)
		.replace(/mm/g, (m < 10 ? '0' : '') + m)
		.replace(/ss/g, (s < 10 ? '0' : '') + s);
}

/**
 * getStorage:获取缓存
 * @param name：获取的名字
 * */
export function getStorage(name) {
	//根据名称获取缓存值 默认为空
	let [storageString = ''] = [localStorage.getItem(name)];
	//如果未设置该缓存或没有值，直接返回空串
	if (storageString == null || storageString.length <= 0) {
		return '';
	}
	try {
		// 把缓存字符串转为对象、获取当前时间
		let [storageObj, nowTime] = [JSON.parse(storageString), new Date().getTime()];
		// 如果缓存对象不存在或为空，或者对象设置了时间值并且时间值小于当前时间(即已过期)，此时清除缓存
		if (!storageObj || (storageObj.expires != undefined && storageObj.expires < nowTime)) {
			cleanStorage(name);
			return null;
		} else {
			// 如果有值且时间未过期，则返回值，如果没有value则没有设置过期时间，此时不是对象，直接返回缓存字符串
			return storageObj.value || storageString;
		}
	} catch (e) {
		return storageString;
	}
}

/**
 * setStorage:设置缓存
 * @param name：设置的名字
 * @param value：设置的值
 * @param time：设置的过期时间 默认30天
 * */
export function setStorage(name, value, time = 30) {
	//获取过期时间的时间戳值和当前的时间戳
	let [days, nowTime] = [parseFloat(time) * 84000000, new Date().getTime()];
	//expires是过期时间
	localStorage.setItem(
		name,
		JSON.stringify({
			value: value,
			expires: nowTime + days,
		}),
	);
}

/**
 * cleanStorage:清除缓存
 * @param name：清除的名字
 * */
export function cleanStorage(name) {
	localStorage.removeItem(name);
}
/** setSessionStorage:设置缓存
 * @param {String} name：设置的名字
 * @param {*} value：设置的值
 * */
export function setSessionStorage(name, value) {
	sessionStorage.setItem(name, value);
}
/** getSessionStorage:获取缓存
 * @param {String} name：设置的名字
 * @return 返回session值
 * */
export function getSessionStorage(name) {
	return sessionStorage.getItem(name);
}
/**
 * cleanStorage:清除缓存
 * @param name：清除的名字
 * */
export function cleanSessionStorage(name) {
	sessionStorage.removeItem(name);
}

/**
 * @name 获取cookies
 */
export function getCookie(name) {
	var arr,
		reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)');
	if ((arr = document.cookie.match(reg))) return unescape(arr[2]);
	else return null;
}
/**
 * @name 删除cookies
 */
export function cleanCookie(name) {
	let domain0 = document.domain.split('.').slice(-2).join('.');
	document.cookie = name + '=0;path=/;expires=' + new Date(0).toUTCString(); // 清除当前域名下的,例如：m.gzzyac.com
	document.cookie =
		name + '=0;path=/;domain=' + document.domain + ';expires=' + new Date(0).toUTCString(); // 清除当前域名下的，例如 .m.gzzyac.com
	document.cookie = name + '=0;path=/;domain=' + domain0 + ';expires=' + new Date(0).toUTCString(); // 清除一级域名下的或指定的，例如 .gzzyac.com
}
/**
 * 当前站点所有域cookie
 */
export function delCookie() {
	let domain0 = document.domain.split('.').slice(-2).join('.');
	var keys = document.cookie.match(/[^ =;]+(?==)/g);
	if (keys) {
		for (var i = keys.length; i--; ) {
			document.cookie = keys[i] + '=0;path=/;expires=' + new Date(0).toUTCString(); // 清除当前域名下的,例如：m.gzzyac.com
			document.cookie =
				keys[i] + '=0;path=/;domain=' + document.domain + ';expires=' + new Date(0).toUTCString(); // 清除当前域名下的，例如 .m.gzzyac.com
			document.cookie =
				keys[i] + '=0;path=/;domain=' + domain0 + ';expires=' + new Date(0).toUTCString(); // 清除一级域名下的或指定的，例如 .gzzyac.com
		}
	}
}

// 下载文件
// 模板下载
export function downTemplate(url, type = 0) {
	let a = document.createElement('a');
	a.href = type == 0 ? ajaxUrl + url : url;
	a.click();
}
// 下载
export function downloadIamge(imgsrc, name) {
	var image = new Image();
	// 解决跨域 Canvas 污染问题
	image.setAttribute('crossOrigin', 'anonymous');
	image.onload = function () {
		var canvas = document.createElement('canvas');
		canvas.width = image.width;
		canvas.height = image.height;
		var context = canvas.getContext('2d');
		context.drawImage(image, 0, 0, image.width, image.height);
		var url = canvas.toDataURL('image/png'); // 得到图片的base64编码数据
		var a = document.createElement('a'); // 生成一个a元素
		var event = new MouseEvent('click'); // 创建一个单击事件
		a.download = name || 'photo'; // 设置图片名称
		a.href = url; // 将生成的URL设置为a.href属性
		a.dispatchEvent(event); // 触发a的单击事件
	};
	image.src = imgsrc;
}
// 复制
export function copyText(txt, success = function () {}) {
	// if (document.getElementById("copy-text")) {
	//   return document.getElementById("copy-text");
	// }
	var idObject = document.getElementById('copy-text');
	if (idObject != null) idObject.parentNode.removeChild(idObject);
	var input = document.createElement('input');
	input.setAttribute('type', 'text');
	input.setAttribute('value', txt);
	input.setAttribute('id', 'copy-text');
	input.setAttribute('style', 'position:absolute;left:9000px;');
	document.body.appendChild(input);
	let copyHtml = document.getElementById('copy-text');
	copyHtml.select(); // 选择对象
	document.execCommand('Copy'); // 执行浏览器复制命令
	success();
	input.remove();
}

// 解决偶现时间传值传的时中国时间问题
export function checkTime(val) {
	if (!val) return;
	for (let i in val) {
		if (
			(i == 'start_time' || i == 'end_time') &&
			val[i] &&
			val[i].toString().indexOf('中国标准时间') != -1
		) {
			val[i] = formatDate(new Date(val[i]).getTime() / 1000);
		}
	}
	return val;
}

/**
 * @description: 在生成二维码前就要把字符串转换成UTF-8
 * @param {*}string
 * @return {*}UTF-8
 */
export const toUtf8 = (str) => {
	let out, i, len, c;
	out = '';
	len = str.length;
	for (i = 0; i < len; i++) {
		c = str.charCodeAt(i);
		if (c >= 0x0001 && c <= 0x007f) {
			out += str.charAt(i);
		} else if (c > 0x07ff) {
			out += String.fromCharCode(0xe0 | ((c >> 12) & 0x0f));
			out += String.fromCharCode(0x80 | ((c >> 6) & 0x3f));
			out += String.fromCharCode(0x80 | ((c >> 0) & 0x3f));
		} else {
			out += String.fromCharCode(0xc0 | ((c >> 6) & 0x1f));
			out += String.fromCharCode(0x80 | ((c >> 0) & 0x3f));
		}
	}
	return out;
};

// 替换url指定值
export function changeURLArg(url, arg, arg_val, isGetObj) {
	let newUrl = window.location.origin + window.location.pathname;
	if (url.split('?')[1]) {
		let temp1 = url.split('?');
		let pram = temp1[1];
		let keyValue = pram.split('&');
		let key = '',
			value = '',
			arr = [],
			obj = {}; //参数对象集合
		for (let i = 0; i < keyValue.length; i++) {
			let item = keyValue[i].split('=');
			key = item[0];
			value = item[1];
			obj[key] = value;
		}
		if (isGetObj) {
			return obj;
		}

		obj[arg] = arg_val;
		let i = 0;
		for (let o_key in obj) {
			let sym = i == 0 ? '?' : '&';
			i++;
			arr.push(sym + o_key + '=' + obj[o_key]);
		}
		return `${newUrl}${arr.join('')}`;
	} else {
		return `${newUrl}?${arg}=${arg_val}`;
	}
}

/**
 * 判断是否有http域名
 */
export function setFileUrl(url) {
	const isHasHttp = new RegExp('http').test(url);
	return isHasHttp ? url : ajaxUrl + url;
}

/**
 * @name: 获取文件后缀
 * @param {String} filePath 文件路径
 */
export const getFileType = (filePath) => {
	var startIndex = filePath.lastIndexOf('.');
	if (startIndex != -1) return filePath.substring(startIndex + 1, filePath.length).toLowerCase();
	else return '';
};

//图片压缩
export const compressUpload = (file, isShowLoading = true) => {
	let _loading = {};
	return new Promise((resolve) => {
		if (isShowLoading) {
			// _loading = ELEMENT.Loading.service({
			// 	text: '压缩中',
			// 	background: 'rgba(0, 0, 0, 0.7)',
			// });
		}
		const reader = new FileReader();
		const image = new Image();
		image.onload = () => {
			const canvas = document.createElement('canvas');
			const context = canvas.getContext('2d');
			const width = image.width * 0.5; //宽缩一半
			const height = image.height * 0.5; //高缩一半
			canvas.width = width;
			canvas.height = height;
			context.clearRect(0, 0, width, height);
			context.drawImage(image, 0, 0, width, height);
			const dataUrl = canvas.toDataURL(file.type);
			const blobData = dataURItoBlob(dataUrl, file.type);
			_loading.close();
			resolve(blobData);
		};
		reader.onload = (e) => {
			image.src = e.target.result;
		};
		reader.readAsDataURL(file);
	});
};
//图片格式转换——base64转Blob对象
export const dataURItoBlob = (dataURI, type) => {
	var binary = atob(dataURI.split(',')[1]);
	var array = [];
	for (var i = 0; i < binary.length; i++) {
		array.push(binary.charCodeAt(i));
	}
	return new Blob([new Uint8Array(array)], { type: type });
};

//计算两个时间字符串的差值
export const getDateDiff = (startTime, endTime) => {
	//将日期字符串转换为时间戳
	var sTime = new Date(startTime).getTime(); //开始时间
	var eTime = new Date(endTime).getTime(); //结束时间

	//作为除数的数字
	var divNumSecond = 1000;
	var divNumMinute = 1000 * 60;
	var divNumHour = 1000 * 3600;
	var divNumDay = 1000 * 3600 * 24;

	const day = parseInt((eTime - sTime) / parseInt(divNumDay));
	const hour = parseInt(((eTime - sTime) % parseInt(divNumDay)) / parseInt(divNumHour));
	const minute = parseInt(
		parseInt(((eTime - sTime) % parseInt(divNumDay)) % parseInt(divNumHour)) /
			parseInt(divNumMinute),
	);
	const second =
		(parseInt(((eTime - sTime) % parseInt(divNumDay)) % parseInt(divNumHour)) %
			parseInt(divNumMinute)) /
		parseInt(divNumSecond);
	const str = day + '天' + hour + '小时' + minute + '分' + second + '秒';
	return str;
};
