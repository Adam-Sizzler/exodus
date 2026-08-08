"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpsertUserMetadataCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var UpsertUserMetadataCommand;
(function (UpsertUserMetadataCommand) {
    UpsertUserMetadataCommand.url = api_1.REST_API.METADATA.USER.UPSERT;
    UpsertUserMetadataCommand.TSQ_url = UpsertUserMetadataCommand.url(':userId');
    UpsertUserMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.METADATA_ROUTES.USER.UPSERT(':userId'), 'put', 'Update or create User Metadata', { scope: 'upsert-user', kind: 'write' });
    UpsertUserMetadataCommand.RequestParamsSchema = zod_1.z.object({
        userId: zod_1.z.coerce.number(),
    });
    UpsertUserMetadataCommand.RequestBodySchema = zod_1.z.object({
        metadata: zod_1.z.looseObject({}),
    });
    UpsertUserMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            metadata: zod_1.z.looseObject({}),
        }),
    });

})(UpsertUserMetadataCommand || (exports.UpsertUserMetadataCommand = UpsertUserMetadataCommand = {}));
