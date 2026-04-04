"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ejs = exports.ViteEjsPlugin = void 0;
const ejs_1 = __importDefault(require("ejs"));
exports.ejs = ejs_1.default;
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
function ViteEjsPlugin(data = {}, options) {
    let config;
    return {
        name: "vite-plugin-ejs",
        // Get Resolved config
        configResolved(resolvedConfig) {
            config = resolvedConfig;
        },
        transformIndexHtml: {
            order: "pre",
            handler(html) {
                if (typeof data === "function")
                    data = data(config);
                let ejsOptions = options && options.ejs ? options.ejs : {};
                if (typeof ejsOptions === "function")
                    ejsOptions = ejsOptions(config);
                html = ejs_1.default.render(html, Object.assign({ NODE_ENV: config.mode, isDev: config.mode === "development" }, data), Object.assign(Object.assign({ 
                    // setting views enables includes support
                    views: [config.root] }, ejsOptions), { async: false // Force sync
                 }));
                return html;
            }
        }
    };
}
exports.ViteEjsPlugin = ViteEjsPlugin;
