"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteNodeIntegrationCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteNodeIntegrationCommand;
(function (DeleteNodeIntegrationCommand) {
    DeleteNodeIntegrationCommand.url = api_1.REST_API.NODE_INTEGRATIONS.DELETE;
    DeleteNodeIntegrationCommand.TSQ_url = DeleteNodeIntegrationCommand.url(':uuid');
    DeleteNodeIntegrationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_INTEGRATIONS_ROUTES.DELETE(':uuid'), 'delete', 'Delete Node Integration', { scope: 'delete', kind: 'write' });
    DeleteNodeIntegrationCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
})(DeleteNodeIntegrationCommand || (exports.DeleteNodeIntegrationCommand = DeleteNodeIntegrationCommand = {}));
