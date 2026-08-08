"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserMetadataCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetUserMetadataCommand;
(function (GetUserMetadataCommand) {
    GetUserMetadataCommand.url = api_1.REST_API.METADATA.USER.GET;
    GetUserMetadataCommand.TSQ_url = GetUserMetadataCommand.url(':userId');
    GetUserMetadataCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.METADATA_ROUTES.USER.GET(':userId'), 'get', 'Get user metadata', { scope: 'get-user', kind: 'read' });
    GetUserMetadataCommand.RequestParamsSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    GetUserMetadataCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            metadata: zod_1.z.looseObject({}),
        }),
    });

})(GetUserMetadataCommand || (exports.GetUserMetadataCommand = GetUserMetadataCommand = {}));
