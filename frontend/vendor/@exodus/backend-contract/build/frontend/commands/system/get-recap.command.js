"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetRecapCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetRecapCommand;
(function (GetRecapCommand) {
    GetRecapCommand.url = api_1.REST_API.SYSTEM.STATS.RECAP;
    GetRecapCommand.TSQ_url = GetRecapCommand.url;
    GetRecapCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.RECAP, 'get', 'Get Recap');
    GetRecapCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            thisMonth: zod_1.z.object({
                users: zod_1.z.number(),
                traffic: zod_1.z.string(),
            }),
            total: zod_1.z.object({
                users: zod_1.z.number(),
                nodes: zod_1.z.number(),
                traffic: zod_1.z.string(),
                nodesRam: zod_1.z.string(),
                nodesCpuCores: zod_1.z.number(),
                distinctCountries: zod_1.z.number(),
            }),
            version: zod_1.z.string(),
            initDate: zod_1.z
                .string()
                .datetime({ local: true, offset: true })
                .transform((str) => new Date(str)),
        }),
    });
})(GetRecapCommand || (exports.GetRecapCommand = GetRecapCommand = {}));
