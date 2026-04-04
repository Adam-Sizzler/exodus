export type Options = {
	/**
	* The CWD for the operation.
	* @default "." (process.cwd)
	*/
	cwd?: string
	/**
	* The directory to stop at.
	*
	* > [NOTE]
	* > This directory WILL NOT be included in the results.
	*
	* @default <none> (continue to system root)
	*/
	stop?: string
};
/**
* Get all parent directories of {@link base}.
* Stops at {@link Options['stop']} else system root ("/").
*
* @returns An array of absolute paths of all parent directories.
*/
export declare function up(base: string, options?: Options): string[];
