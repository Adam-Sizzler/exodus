"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserByTagCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var GetUserByTagCommand;
(function (GetUserByTagCommand) {
    GetUserByTagCommand.url = api_1.REST_API.USERS.GET_BY.TAG;
    GetUserByTagCommand.TSQ_url = GetUserByTagCommand.url(':tag');
    GetUserByTagCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.GET_BY.TAG(':tag'), 'get', 'Get users by tag');
    GetUserByTagCommand.RequestSchema = zod_1.z.object({
        tag: zod_1.z
            .string()
            .regex(/^[A-Z0-9_]+$/, 'Tag can only contain uppercase letters, numbers, underscores')
            .max(16, 'Tag must be less than 16 characters'),
    });
    GetUserByTagCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.ExtendedUsersSchema),
    });
})(GetUserByTagCommand || (exports.GetUserByTagCommand = GetUserByTagCommand = {}));
