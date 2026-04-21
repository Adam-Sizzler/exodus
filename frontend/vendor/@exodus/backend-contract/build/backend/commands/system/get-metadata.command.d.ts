import { z } from 'zod';
export declare namespace GetMetadataCommand {
    const url: "/api/system/metadata";
    const TSQ_url: "/api/system/metadata";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            version: z.ZodString;
            build: z.ZodObject<{
                time: z.ZodString;
                number: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                number: string;
                time: string;
            }, {
                number: string;
                time: string;
            }>;
            git: z.ZodObject<{
                repositoryUrl: z.ZodString;
                backend: z.ZodObject<{
                    commitSha: z.ZodString;
                    branch: z.ZodString;
                    commitUrl: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                }, {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                }>;
                frontend: z.ZodObject<{
                    commitSha: z.ZodString;
                    commitUrl: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    commitSha: string;
                    commitUrl: string;
                }, {
                    commitSha: string;
                    commitUrl: string;
                }>;
            }, "strip", z.ZodTypeAny, {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            }, {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            }>;
        }, "strip", z.ZodTypeAny, {
            version: string;
            build: {
                number: string;
                time: string;
            };
            git: {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            };
        }, {
            version: string;
            build: {
                number: string;
                time: string;
            };
            git: {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            version: string;
            build: {
                number: string;
                time: string;
            };
            git: {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            };
        };
    }, {
        response: {
            version: string;
            build: {
                number: string;
                time: string;
            };
            git: {
                repositoryUrl: string;
                backend: {
                    commitSha: string;
                    branch: string;
                    commitUrl: string;
                };
                frontend: {
                    commitSha: string;
                    commitUrl: string;
                };
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-metadata.command.d.ts.map
