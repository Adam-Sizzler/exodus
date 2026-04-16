"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSnippetCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateSnippetCommand;
(function (CreateSnippetCommand) {
    CreateSnippetCommand.url = api_1.REST_API.SNIPPETS.CREATE;
    CreateSnippetCommand.TSQ_url = CreateSnippetCommand.url;
    CreateSnippetCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SNIPPETS_ROUTES.CREATE, 'post', 'Create snippet');
    CreateSnippetCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(255, 'Name must be less than 255 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
        snippet: zod_1.z.array(zod_1.z.object({}).passthrough()),
    });
    CreateSnippetCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            snippets: zod_1.z.array(models_1.SnippetsSchema),
        }),
    });
})(CreateSnippetCommand || (exports.CreateSnippetCommand = CreateSnippetCommand = {}));
