const { join } = require("node:path");
const { existsSync } = require("node:fs");
const walk = require("empathic/walk");
/**
* Find an item by name, walking parent directories until found.
*
* @param name The item name to find.
* @returns The absolute path to the item, if found.
*/
function up(name, options) {
	let dir, tmp;
	let start = options && options.cwd || "";
	for (dir of walk.up(start, options)) {
		tmp = join(dir, name);
		if (existsSync(tmp)) return tmp;
	}
}
/**
* Get the first path that matches any of the names provided.
*
* > [NOTE]
* > The order of {@link names} is respected.
*
* @param names The item names to find.
* @returns The absolute path of the first item found, if any.
*/
function any(names, options) {
	let dir, start = options && options.cwd || "";
	let j = 0, len = names.length, tmp;
	for (dir of walk.up(start, options)) {
		for (j = 0; j < len; j++) {
			tmp = join(dir, names[j]);
			if (existsSync(tmp)) return tmp;
		}
	}
}

exports.up = up;
exports.any = any;