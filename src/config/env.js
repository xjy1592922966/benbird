/*
	声明变量区域
*/

let axiosUrl, ajaxUrl, imgUrl;

const PRODUCTION_APIURL = ''; //生产环境后台API
const DEVELOP_APIURL = ''; //测试环境后台API
const LOCALHOST_APIURL = 'http://localhost:8080'; //本地环境后台API

const PRODUCTION_WENURL = ''; //生产环境前台
const DEVELOP_WENURL = ''; //测试环境前台
const LOCALHOST_WENURL = 'http://localhost:8081/'; //本地环境前台

/*
基于域名判断当前走哪个接口
*/
if (document.baseURI.indexOf('benbird.hcolor.pro') > -1) {
	axiosUrl = PRODUCTION_APIURL + '/api/';
	ajaxUrl = PRODUCTION_APIURL;
	imgUrl = PRODUCTION_WENURL;
}
if (document.baseURI.indexOf('benbird-test.hcolor.pro') > -1) {
	axiosUrl = DEVELOP_APIURL + '/api/';
	ajaxUrl = DEVELOP_APIURL;
	imgUrl = DEVELOP_WENURL; //测试环境前台
}
if (document.baseURI.indexOf('localhost') > -1) {
	axiosUrl = LOCALHOST_APIURL;
	ajaxUrl = LOCALHOST_APIURL;
	imgUrl = LOCALHOST_WENURL;
}

export { axiosUrl, ajaxUrl, imgUrl };
