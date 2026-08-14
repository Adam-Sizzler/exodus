"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeMetadataCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetNodeMetadataCommand;
(function (GetNodeMetadataCommand) {
    GetNodeMetadataCommand.url = api_1.REST_API.METADATA.NODE.GET;
    GetNodeMetadataCommand.TSQ_url = GetNodeMetadataCommand.url(':uuid');
    GetNodeMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.METADATA_ROUTES.NODE.GET(':uuid'), 'get', 'Get node metadata', { scope: 'get-node', kind: 'read' });
    GetNodeMetadataCommand.RequestParamsSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetNodeMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            metadata: zod_1.z.looseObject({}),
        }),
    });
})(GetNodeMetadataCommand || (exports.GetNodeMetadataCommand = GetNodeMetadataCommand = {}));
