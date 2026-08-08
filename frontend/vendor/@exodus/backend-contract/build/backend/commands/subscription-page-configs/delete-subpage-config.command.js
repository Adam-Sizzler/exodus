"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteSubpageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteSubpageConfigCommand;
(function (DeleteSubpageConfigCommand) {
    DeleteSubpageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.DELETE;
    DeleteSubpageConfigCommand.TSQ_url = DeleteSubpageConfigCommand.url(':uuid');
    DeleteSubpageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.DELETE(':uuid'), 'delete', 'Delete subscription page config', { scope: 'delete', kind: 'write' });
    DeleteSubpageConfigCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteSubpageConfigCommand || (exports.DeleteSubpageConfigCommand = DeleteSubpageConfigCommand = {}));
