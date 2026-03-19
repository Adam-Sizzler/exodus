"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DisableNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var DisableNodeCommand;
(function (DisableNodeCommand) {
    DisableNodeCommand.url = api_1.REST_API.NODES.ACTIONS.DISABLE;
    DisableNodeCommand.TSQ_url = DisableNodeCommand.url(':uuid');
    DisableNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.DISABLE(':uuid'), 'post', 'Disable a node');
    DisableNodeCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DisableNodeCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodesSchema,
    });
})(DisableNodeCommand || (exports.DisableNodeCommand = DisableNodeCommand = {}));
