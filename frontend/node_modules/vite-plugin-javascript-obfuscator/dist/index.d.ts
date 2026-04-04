import { ObfuscatorOptions } from "javascript-obfuscator";
import { Matcher } from "anymatch";
type Options = {
    /**
     * (Array|String|RegExp|Function) String to be directly matched, string with glob patterns, regular expression test, function that takes the testString as an argument and returns a truthy value if it should be matched. default: ```[/\.(jsx?|tsx?|cjs|mjs)$/]```
     * [See more](https://github.com/micromatch/anymatch)
     */
    include?: Matcher;
    /**
     *  (Array|String|RegExp|Function) String to be directly matched, string with glob patterns, regular expression test, function that takes the testString as an argument and returns a truthy value if it should be matched. default: ```[/node_modules/, /\.nuxt/]```
     * [See more](https://github.com/micromatch/anymatch)
     */
    exclude?: Matcher;
    /**
     * Your javascript-obfuscator options
     * [See more options](https://github.com/javascript-obfuscator/javascript-obfuscator)
     */
    options?: ObfuscatorOptions;
    /**
     * Used for debugging, Print out the path of matching or excluding files
     */
    debugger?: boolean;
    /**
     * By default plugins are invoked for both serve and build. In cases where a plugin needs to be conditionally applied only during serve or build
     * https://vitejs.dev/guide/api-plugin.html
     */
    apply?: "serve" | "build" | ((this: void, config: any, env: any) => boolean);
};
export default function obfuscatorPlugin(obOptions?: Options): {
    name: string;
    enforce: "post";
    apply: "serve" | "build" | ((this: void, config: any, env: any) => boolean);
    transform(src: string, id: string): {
        map: string;
        code: string;
    } | undefined;
};
export {};
