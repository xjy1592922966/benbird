const fs = require('fs');

class InsertTimestampPlugin {
	apply(compiler) {
		compiler.hooks.afterEmit.tap('InsertTimestampPlugin', () => {
			const timestamp = Date.now();
			const content = JSON.stringify(timestamp);
			const filename = `new_version.txt`;
			const filepath = `${compiler.options.output.path}/${filename}`;

			//判断是否存在目录
			if (!fs.existsSync(filepath)) {
				//如果不存在目录。就新建目录，再创建文件夹
				fs.mkdirSync(compiler.options.output.path, { recursive: true });
				fs.writeFileSync(filepath, content);
				console.log(`The timestamp file has been created at ${filepath}.`);
			} else {
				//如果存在文件，就直接更新
				fs.writeFile(filepath, content, (err) => {
					if (err) throw err;
					console.log(`The timestamp file has been updated at ${filepath}.`);
				});
			}
		});
	}
}

module.exports = InsertTimestampPlugin;
