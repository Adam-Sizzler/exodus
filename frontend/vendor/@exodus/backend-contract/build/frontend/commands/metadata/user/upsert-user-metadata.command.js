"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpsertUserMetadataCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var UpsertUserMetadataCommand;
(function (UpsertUserMetadataCommand) {
    UpsertUserMetadataCommand.url = api_1.REST_API.METADATA.USER.UPSERT;
    UpsertUserMetadataCommand.TSQ_url = UpsertUserMetadataCommand.url(':uuid');
    UpsertUserMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.METADATA_ROUTES.USER.UPSERT(':uuid'), 'put', 'Update or create User Metadata', { scope: 'upsert-user', kind: 'write' });
    UpsertUserMetadataCommand.RequestParamsSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    UpsertUserMetadataCommand.RequestBodySchema = zod_1.z.object({
        metadata: zod_1.z.object({}).passthrough(),
    });
    UpsertUserMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            metadata: zod_1.z.object({}).passthrough(),
        }),
    });
})(UpsertUserMetadataCommand || (exports.UpsertUserMetadataCommand = UpsertUserMetadataCommand = {}));
