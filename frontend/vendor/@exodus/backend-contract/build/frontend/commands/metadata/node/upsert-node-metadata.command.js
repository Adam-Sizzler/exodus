"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpsertNodeMetadataCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var UpsertNodeMetadataCommand;
(function (UpsertNodeMetadataCommand) {
    UpsertNodeMetadataCommand.url = api_1.REST_API.METADATA.NODE.UPSERT;
    UpsertNodeMetadataCommand.TSQ_url = UpsertNodeMetadataCommand.url(':uuid');
    UpsertNodeMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.METADATA_ROUTES.NODE.UPSERT(':uuid'), 'put', 'Update or create Node Metadata', { scope: 'upsert-node', kind: 'write' });
    UpsertNodeMetadataCommand.RequestParamsSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    UpsertNodeMetadataCommand.RequestBodySchema = zod_1.z.object({
        metadata: zod_1.z.looseObject({}),
    });
    UpsertNodeMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            metadata: zod_1.z.looseObject({}),
        }),
    });
})(UpsertNodeMetadataCommand || (exports.UpsertNodeMetadataCommand = UpsertNodeMetadataCommand = {}));
