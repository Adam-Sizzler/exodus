import { dirname } from "node:path";
import { absolute } from "empathic/resolve";
/**
* Get all parent directories of {@link base}.
* Stops at {@link Options['stop']} else system root ("/").
*
* @returns An array of absolute paths of all parent directories.
*/
export function up(base, options) {
	let { stop, cwd } = options || {};
	let tmp = absolute(base, cwd), root = !stop;
	let prev, arr = [];
	if (stop) stop = absolute(stop, cwd);
	while (root || tmp !== stop) {
		arr.push(tmp);
		tmp = dirname(prev = tmp);
		if (tmp === prev) break;
	}
	return arr;
}
