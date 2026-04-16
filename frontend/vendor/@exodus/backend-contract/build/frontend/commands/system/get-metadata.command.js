"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetMetadataCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetMetadataCommand;
(function (GetMetadataCommand) {
    GetMetadataCommand.url = api_1.REST_API.SYSTEM.METADATA;
    GetMetadataCommand.TSQ_url = GetMetadataCommand.url;
    GetMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.METADATA, 'get', 'Get Exodus Information');
    GetMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            version: zod_1.z.string(),
            build: zod_1.z.object({
                time: zod_1.z.string(),
                number: zod_1.z.string(),
            }),
            git: zod_1.z.object({
                backend: zod_1.z.object({
                    commitSha: zod_1.z.string(),
                    branch: zod_1.z.string(),
                    commitUrl: zod_1.z.string(),
                }),
                frontend: zod_1.z.object({
                    commitSha: zod_1.z.string(),
                    commitUrl: zod_1.z.string(),
                }),
            }),
        }),
    });
})(GetMetadataCommand || (exports.GetMetadataCommand = GetMetadataCommand = {}));
