import { Plugin, ResolvedConfig } from "vite";
import ejs from "ejs";
type EjsRenderOptions = ejs.Options & {
    async?: false;
};
type EjsRenderOptionsFn = (config: ResolvedConfig) => EjsRenderOptions;
type ViteEjsPluginDataType = Record<string, any> | ((config: ResolvedConfig) => Record<string, any>);
type ViteEjsPluginOptions = {
    ejs?: EjsRenderOptions | EjsRenderOptionsFn;
};
/**
 * Vite Ejs Plugin Function
 * See https://github.com/trapcodeio/vite-plugin-ejs for more information
 * @example
 * export default defineConfig({
 *  plugins: [
 *    vue(),
 *    ViteEjsPlugin({foo: 'bar'})
 *  ],
 * });
 */
declare function ViteEjsPlugin(data?: ViteEjsPluginDataType, options?: ViteEjsPluginOptions): Plugin;
export { ViteEjsPlugin, ejs };
