"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSshTicketCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var CreateSshTicketCommand;
(function (CreateSshTicketCommand) {
    CreateSshTicketCommand.url = api_1.REST_API.NODE_SSH.CREATE_TICKET;
    CreateSshTicketCommand.TSQ_url = CreateSshTicketCommand.url(':uuid');
    CreateSshTicketCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_SSH_ROUTES.CREATE_TICKET(':uuid'), 'post', 'Create a single-use ticket for opening an SSH terminal session', { scope: 'node-ssh', kind: 'write' });
    CreateSshTicketCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid().describe('Node UUID'),
    });
    CreateSshTicketCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            ticket: zod_1.z.string(),
            path: zod_1.z.string(),
            expiresInSeconds: zod_1.z.number(),
        }),
    });
})(CreateSshTicketCommand || (exports.CreateSshTicketCommand = CreateSshTicketCommand = {}));
