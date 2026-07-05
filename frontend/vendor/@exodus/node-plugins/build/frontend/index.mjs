import * as models from './models/index.js';

const merged = { ...models };
export const ConnectionDropPluginSchema = merged.ConnectionDropPluginSchema;
export const EgressFilterPluginSchema = merged.EgressFilterPluginSchema;
export const IngressFilterPluginSchema = merged.IngressFilterPluginSchema;
export const NodePluginSchema = merged.NodePluginSchema;
export const SharedListSchema = merged.SharedListSchema;
export const TorrentBlockerPluginSchema = merged.TorrentBlockerPluginSchema;

export default merged;
