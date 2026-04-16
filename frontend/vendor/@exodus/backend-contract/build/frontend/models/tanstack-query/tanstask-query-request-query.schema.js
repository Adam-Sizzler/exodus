"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TanstackQueryRequestQuerySchema = void 0;
const zod_1 = require("zod");
const tanstask_query_sorting_schema_1 = require("./tanstask-query-sorting.schema");
const tanstask_query_filter_schema_1 = require("./tanstask-query-filter.schema");
exports.TanstackQueryRequestQuerySchema = zod_1.z.object({
    start: zod_1.z.coerce
        .number()
        .default(0)
        .describe('Start index (offset) of the results to return, default is 0'),
    size: zod_1.z.coerce
        .number()
        .min(1, 'Size (limit) must be greater than 0')
        .max(1000, 'Size (limit) must be less than 1000')
        .describe('Number of results to return, no more than 1000')
        .default(25),
    filters: zod_1.z
        .preprocess((str) => (typeof str === 'string' ? JSON.parse(str) : str), zod_1.z.array(tanstask_query_filter_schema_1.TanstackQueryFilterSchema))
        .optional(),
    filterModes: zod_1.z
        .preprocess((str) => (typeof str === 'string' ? JSON.parse(str) : str), zod_1.z.record(zod_1.z.string(), zod_1.z.string()))
        .optional(),
    globalFilterMode: zod_1.z.string().optional(),
    sorting: zod_1.z
        .preprocess((str) => (typeof str === 'string' ? JSON.parse(str) : str), zod_1.z.array(tanstask_query_sorting_schema_1.TanstackQuerySortingSchema))
        .optional(),
});
