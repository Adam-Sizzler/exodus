"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExtendUserCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const user_response_1 = require("../user.response");
var ExtendUserCommand;
(function (ExtendUserCommand) {
    ExtendUserCommand.url = api_1.REST_API.USERS.ACTIONS.EXTEND_EXPIRATION_DATE;
    ExtendUserCommand.TSQ_url = ExtendUserCommand.url(':userId');
    ExtendUserCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.ACTIONS.EXTEND_EXPIRATION_DATE(':userId'), 'post', 'Extend user expiration date', { scope: 'extend', kind: 'write' }, 'If user status is EXPIRED, the new expiration date is calculated from the current date and the user becomes ACTIVE. If user status is ACTIVE, the given number of days is added to the existing expiration date. DISABLED and LIMITED users will be extended, but their status will not change.');
    ExtendUserCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema,
    });
    ExtendUserCommand.RequestBodySchema = zod_1.z.object({
        days: zod_1.z.number().min(1).describe('The number of days to extend the expiration date.'),
    });
    ExtendUserCommand.ResponseSchema = user_response_1.UserResponseSchema;
})(ExtendUserCommand || (exports.ExtendUserCommand = ExtendUserCommand = {}));
