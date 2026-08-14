"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkNodesProfileModificationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkNodesProfileModificationCommand;
(function (BulkNodesProfileModificationCommand) {
    BulkNodesProfileModificationCommand.url = api_1.REST_API.NODES.BULK_ACTIONS.PROFILE_MODIFICATION;
    BulkNodesProfileModificationCommand.TSQ_url = BulkNodesProfileModificationCommand.url;
    BulkNodesProfileModificationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.BULK_ACTIONS.PROFILE_MODIFICATION, 'post', 'Modify Inbounds & Profile for many nodes', { scope: 'bulk-profile-modification', kind: 'write' });
    BulkNodesProfileModificationCommand.RequestBodySchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.uuid()).min(1),
        configProfile: zod_1.z.object({
            activeConfigProfileUuid: zod_1.z.uuid(),
            activeInbounds: zod_1.z.array(zod_1.z.uuid()).min(1),
        }),
    });
})(BulkNodesProfileModificationCommand || (exports.BulkNodesProfileModificationCommand = BulkNodesProfileModificationCommand = {}));
