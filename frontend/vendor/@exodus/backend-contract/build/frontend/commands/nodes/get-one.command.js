"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetOneNodeCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var GetOneNodeCommand;
(function (GetOneNodeCommand) {
    GetOneNodeCommand.url = api_1.REST_API.NODES.GET_BY_UUID;
    GetOneNodeCommand.TSQ_url = GetOneNodeCommand.url(':uuid');
    GetOneNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get node by UUID');
    GetOneNodeCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetOneNodeCommand.ResponseSchema = zod_1.z.object({
        response: models_1.NodesSchema,
    });
})(GetOneNodeCommand || (exports.GetOneNodeCommand = GetOneNodeCommand = {}));
