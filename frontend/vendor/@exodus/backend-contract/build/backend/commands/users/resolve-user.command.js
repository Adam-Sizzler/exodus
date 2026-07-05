"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResolveUserCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var ResolveUserCommand;
(function (ResolveUserCommand) {
    ResolveUserCommand.url = api_1.REST_API.USERS.RESOLVE;
    ResolveUserCommand.TSQ_url = ResolveUserCommand.url;
    ResolveUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.RESOLVE, 'post', 'Resolve a user');
    ResolveUserCommand.RequestSchema = zod_1.z
        .object({
        uuid: zod_1.z.string().uuid().optional(),
        id: zod_1.z.number().optional(),
        shortUuid: zod_1.z.string().optional(),
        username: zod_1.z.string().optional(),
    })
        .refine((data) => {
        const provided = [data.uuid, data.id, data.shortUuid, data.username].filter((v) => v !== undefined);
        return provided.length === 1;
    }, {
        message: 'Exactly one of uuid, id, shortUuid, or username must be provided',
    });
    ResolveUserCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            uuid: zod_1.z.string().uuid(),
            username: zod_1.z.string(),
            id: zod_1.z.number(),
            shortUuid: zod_1.z.string(),
        }),
    });
})(ResolveUserCommand || (exports.ResolveUserCommand = ResolveUserCommand = {}));
