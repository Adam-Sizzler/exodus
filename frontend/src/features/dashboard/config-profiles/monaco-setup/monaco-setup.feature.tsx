import {
    ResponseRulesConfigSchema,
    TSubscriptionTemplateType
} from '@exodus/backend-contract'
import zodToJsonSchema, { jsonDescription } from 'zod-to-json-schema'
import { Monaco } from '@monaco-editor/react'
import consola from 'consola'
import axios from 'axios'

import { monacoTheme } from '@shared/constants/monaco-theme'
import { app } from 'src/config'

export const MonacoSetupFeature = {
    setup: async (monaco: Monaco, currentLanguage: string) => {
        try {
            let { jsonSchemaUrl } = app.configEditor
            switch (currentLanguage) {
                case 'zh':
                    jsonSchemaUrl = app.configEditor.jsonSchemaCnUrl
                    break
                default:
                    jsonSchemaUrl = app.configEditor.jsonSchemaUrl
            }

            const response = await axios.get(jsonSchemaUrl)
            const schema = await response.data

            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                allowComments: false,
                enableSchemaRequest: true,
                schemaRequest: 'warning',
                schemas: [
                    {
                        fileMatch: ['*'],
                        schema,
                        uri: 'https://singbox-config-schema.json'
                    }
                ],
                validate: true
            })
        } catch (error) {
            consola.error('Failed to load JSON schema:', error)
        }
    }
}
export const MonacoSetupResponseRulesFeature = {
    setup: async (
        monaco: Monaco,
        groupedTemplates: Record<TSubscriptionTemplateType, string[]>
    ) => {
        try {
            const schema = zodToJsonSchema(ResponseRulesConfigSchema, {
                name: 'Response Rules Config Schema',
                applyRegexFlags: true,
                errorMessages: true,
                postProcess: jsonDescription
            })

            const templateOptions = {
                BROWSER: [],
                BLOCK: [],
                STATUS_CODE_404: [],
                STATUS_CODE_451: [],
                SOCKET_DROP: [],
                ...groupedTemplates
            }

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const schemaDefinitions = (schema as any).definitions?.['Response Rules Config Schema']
            const rulesItems = schemaDefinitions?.properties?.rules?.items

            if (rulesItems) {
                if (!rulesItems.allOf) {
                    rulesItems.allOf = []
                }

                Object.entries(templateOptions).forEach(([responseType, templates]) => {
                    if (templates.length > 0) {
                        rulesItems.allOf.push({
                            if: {
                                properties: {
                                    responseType: { const: responseType }
                                },
                                required: ['responseType']
                            },
                            then: {
                                properties: {
                                    responseModifications: {
                                        properties: {
                                            subscriptionTemplate: {
                                                enum: templates,
                                                markdownDescription: `Available templates for **${responseType}** response type.`,
                                                markdownEnumDescriptions: templates.map(
                                                    (t) => `Use ${t} template`
                                                )
                                            }
                                        }
                                    }
                                }
                            }
                        })
                    } else {
                        rulesItems.allOf.push({
                            if: {
                                properties: {
                                    responseType: { const: responseType }
                                },
                                required: ['responseType']
                            },
                            then: {
                                properties: {
                                    responseModifications: {
                                        properties: {
                                            subscriptionTemplate: {
                                                type: 'null',
                                                not: { type: 'string' },
                                                markdownDescription: `⚠️ No templates available for **${responseType}** response type. This field should not be used.`
                                            }
                                        }
                                    }
                                }
                            }
                        })
                    }
                })
                //     } else {
                //         rulesItems.allOf.push({
                //             if: {
                //                 properties: {
                //                     responseType: { const: responseType }
                //                 },
                //                 required: ['responseType']
                //             },
                //             then: {
                //                 properties: {
                //                     responseModifications: {
                //                         properties: {
                //                             overrideSubscriptionTemplateWith: false
                //                         }
                //                     }
                //                 }
                //             }
                //         })
                //     }
                // })
            }

            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                schemaValidation: 'error',
                comments: 'error',
                trailingCommas: 'error',

                schemas: [
                    {
                        fileMatch: ['response-rules://*'],
                        schema,
                        uri: 'https://response-rules-schema.json'
                    }
                ],
                validate: true
            })

            monaco.languages.json.jsonDefaults.setModeConfiguration({
                documentFormattingEdits: true,
                documentRangeFormattingEdits: true,
                completionItems: true,
                hovers: true,
                documentSymbols: true,
                tokens: true,
                colors: true,
                foldingRanges: true,
                diagnostics: true,
                selectionRanges: true
            })

            monaco.editor.defineTheme('GithubDark', {
                ...monacoTheme,
                base: 'vs-dark'
            })
        } catch (error) {
            consola.error('Failed to load JSON schema:', error)
        }
    }
}
