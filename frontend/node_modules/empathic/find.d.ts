import type { Options } from "empathic/walk";
export type { Options };
/**
* Find an item by name, walking parent directories until found.
*
* @param name The item name to find.
* @returns The absolute path to the item, if found.
*/
export declare function up(name: string, options?: Options): string | undefined;
/**
* Get the first path that matches any of the names provided.
*
* > [NOTE]
* > The order of {@link names} is respected.
*
* @param names The item names to find.
* @returns The absolute path of the first item found, if any.
*/
export declare function any(names: string[], options?: Options): string | undefined;
