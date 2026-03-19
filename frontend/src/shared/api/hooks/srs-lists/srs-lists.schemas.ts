import { z } from 'zod'

export const SRSListSchema = z.object({
    uuid: z.string().uuid(),
    tag: z.string(),
    format: z.string(),
    url: z.string().url(),
    updateInterval: z.string(),
    path: z.string().nullable().optional(),
    fileName: z.string(),
    shortName: z.string(),
    viewPosition: z.number(),
    isEnabled: z.boolean(),
    isAvailable: z.boolean(),
    lastCheckedAt: z.string().datetime().nullable().optional(),
    lastError: z.string().nullable().optional(),
    createdAt: z.string().datetime(),
    updatedAt: z.string().datetime()
})

export const GetSRSListsResponseSchema = z.object({
    response: z.object({
        srsLists: z.array(SRSListSchema)
    })
})

export const CreateSRSListsRequestSchema = z.object({
    url: z.string().optional(),
    urls: z.array(z.string()).optional(),
    format: z.string().optional(),
    updateInterval: z.string().optional(),
    tag: z.string().optional(),
    path: z.string().optional(),
    isEnabled: z.boolean().optional()
})

export const UpdateSRSListRequestSchema = z.object({
    uuid: z.string().uuid(),
    url: z.string().optional(),
    tag: z.string().optional(),
    format: z.string().optional(),
    updateInterval: z.string().optional(),
    path: z.string().optional(),
    isEnabled: z.boolean().optional()
})

export const ReorderSRSListsRequestSchema = z.object({
    items: z.array(
        z.object({
            uuid: z.string().uuid(),
            viewPosition: z.number()
        })
    )
})

export const BulkDeleteSRSListsRequestSchema = z.object({
    uuids: z.array(z.string().uuid()).min(1)
})

export const CheckSRSListsRequestSchema = z.object({
    uuids: z.array(z.string().uuid()).optional()
})

export const BulkEnableSRSListsRequestSchema = z.object({
    uuids: z.array(z.string().uuid()).min(1)
})

export const BulkSetIntervalSRSListsRequestSchema = z.object({
    uuids: z.array(z.string().uuid()).min(1),
    updateInterval: z.string().min(1)
})

export const GenericSRSListsMutationResponseSchema = z.object({
    response: z.any()
})
