import { notifications } from '@mantine/notifications'

import { createMutationHook } from '../../tsq-helpers'
import {
    BulkDeleteSRSListsRequestSchema,
    BulkEnableSRSListsRequestSchema,
    BulkSetIntervalSRSListsRequestSchema,
    CheckSRSListsRequestSchema,
    CreateSRSListsRequestSchema,
    GenericSRSListsMutationResponseSchema,
    ReorderSRSListsRequestSchema,
    UpdateSRSListRequestSchema
} from './srs-lists.schemas'

export const useCreateSRSLists = createMutationHook({
    endpoint: '/api/srs-lists',
    bodySchema: CreateSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'SRS lists created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: 'Create SRS lists',
                message: error instanceof Error ? error.message : 'Request failed',
                color: 'red'
            })
        }
    }
})

export const useUpdateSRSList = createMutationHook({
    endpoint: '/api/srs-lists',
    bodySchema: UpdateSRSListRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'patch'
})

export const useDeleteSRSList = createMutationHook({
    endpoint: '/api/srs-lists/:uuid',
    routeParamsSchema: UpdateSRSListRequestSchema.pick({ uuid: true }),
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'delete'
})

export const useBulkDeleteSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/bulk/delete',
    bodySchema: BulkDeleteSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})

export const useReorderSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/actions/reorder',
    bodySchema: ReorderSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})

export const useCheckSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/actions/check',
    bodySchema: CheckSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})

export const useSyncSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/actions/sync',
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'SRS sync queued for connected nodes',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: 'Sync SRS lists',
                message: error instanceof Error ? error.message : 'Request failed',
                color: 'red'
            })
        }
    }
})

export const useBulkEnableSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/bulk/enable',
    bodySchema: BulkEnableSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})

export const useBulkDisableSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/bulk/disable',
    bodySchema: BulkEnableSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})

export const useBulkSetIntervalSRSLists = createMutationHook({
    endpoint: '/api/srs-lists/bulk/set-interval',
    bodySchema: BulkSetIntervalSRSListsRequestSchema,
    responseSchema: GenericSRSListsMutationResponseSchema,
    requestMethod: 'post'
})
