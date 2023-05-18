const InsertTimestampPlugin = require('./plugins/InsertTimestampPlugin'); //打包版本插件
const os = require('os'); // 系统底层支持
const path = require('path'); // 目录路径支持
const HappyPack = require('happypack'); // 多线程打包

module.exports = {
	configureWebpack: (config) => {
		//打包加入版本号
		config.plugins.push(new InsertTimestampPlugin());
		//多线程打包
		config.plugins.push(
			new HappyPack({
				id: 'babel',
				threads: os.cpus().length - 1,
				loaders: [
					{
						loader: 'babel-loader',
						options: {
							cacheDirectory: true, // 开启 babel-loader 缓存
						},
					},
				],
			}),
		);

		//打包加入版本号
		config.plugins.push(new InsertTimestampPlugin());

		//打包解析别名
		config.resolve = {
			alias: {
				'@': path.resolve(__dirname, 'src'),
			},
		};
	},
};
