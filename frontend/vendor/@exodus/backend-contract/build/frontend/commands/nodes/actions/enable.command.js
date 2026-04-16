"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EnableNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var EnableNodeCommand;
(function (EnableNodeCommand) {
    EnableNodeCommand.url = api_1.REST_API.NODES.ACTIONS.ENABLE;
    EnableNodeCommand.TSQ_url = EnableNodeCommand.url(':uuid');
    EnableNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.ACTIONS.ENABLE(':uuid'), 'post', 'Enable a node');
    EnableNodeCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    EnableNodeCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodesSchema,
    });
})(EnableNodeCommand || (exports.EnableNodeCommand = EnableNodeCommand = {}));
