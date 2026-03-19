"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkEnableHostsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
var BulkEnableHostsCommand;
(function (BulkEnableHostsCommand) {
    BulkEnableHostsCommand.url = api_1.REST_API.HOSTS.BULK.ENABLE_HOSTS;
    BulkEnableHostsCommand.TSQ_url = BulkEnableHostsCommand.url;
    BulkEnableHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.ENABLE_HOSTS, 'post', 'Enable hosts by UUIDs');
    BulkEnableHostsCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()),
    });
    BulkEnableHostsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(models_1.HostsSchema),
    });
})(BulkEnableHostsCommand || (exports.BulkEnableHostsCommand = BulkEnableHostsCommand = {}));
