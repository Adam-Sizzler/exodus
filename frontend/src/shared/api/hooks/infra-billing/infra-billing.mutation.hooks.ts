import { notifications } from '@mantine/notifications'
import {
    CreateInfraBillingRecordCommand,
    CreateInfraBillingNodeCommand,
    CreateInfraProviderCommand,
    DeleteInfraBillingRecordCommand,
    DeleteInfraBillingNodeCommand,
    DeleteInfraProviderCommand,
    UpdateInfraBillingNodeCommand,
    UpdateInfraProviderCommand
} from '@exodus/backend-contract'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateInfraProvider = createMutationHook({
    endpoint: UpdateInfraProviderCommand.TSQ_url,
    bodySchema: UpdateInfraProviderCommand.RequestBodySchema,
    responseSchema: UpdateInfraProviderCommand.ResponseSchema,
    requestMethod: UpdateInfraProviderCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Provider updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Infra Provider`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteInfraProvider = createMutationHook({
    endpoint: DeleteInfraProviderCommand.TSQ_url,
    routeParamsSchema: DeleteInfraProviderCommand.RequestParamSchema,
    requestMethod: DeleteInfraProviderCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Provider deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Infra Provider`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCreateInfraProvider = createMutationHook({
    endpoint: CreateInfraProviderCommand.TSQ_url,
    responseSchema: CreateInfraProviderCommand.ResponseSchema,
    bodySchema: CreateInfraProviderCommand.RequestBodySchema,
    requestMethod: CreateInfraProviderCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Provider created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Infra Provider`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteInfraBillingHistoryRecord = createMutationHook({
    endpoint: DeleteInfraBillingRecordCommand.TSQ_url,
    routeParamsSchema: DeleteInfraBillingRecordCommand.RequestParamSchema,
    requestMethod: DeleteInfraBillingRecordCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Billing History Record deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Infra Billing History Record`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCreateInfraBillingHistoryRecord = createMutationHook({
    endpoint: CreateInfraBillingRecordCommand.TSQ_url,
    responseSchema: CreateInfraBillingRecordCommand.ResponseSchema,
    bodySchema: CreateInfraBillingRecordCommand.RequestBodySchema,
    requestMethod: CreateInfraBillingRecordCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Billing History Record created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Infra Billing History Record`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCreateInfraBillingNode = createMutationHook({
    endpoint: CreateInfraBillingNodeCommand.TSQ_url,
    responseSchema: CreateInfraBillingNodeCommand.ResponseSchema,
    bodySchema: CreateInfraBillingNodeCommand.RequestBodySchema,
    requestMethod: CreateInfraBillingNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Billing Node created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Infra Billing Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteInfraBillingNode = createMutationHook({
    endpoint: DeleteInfraBillingNodeCommand.TSQ_url,
    routeParamsSchema: DeleteInfraBillingNodeCommand.RequestParamSchema,
    requestMethod: DeleteInfraBillingNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Billing Node deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Infra Billing Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useUpdateInfraBillingNode = createMutationHook({
    endpoint: UpdateInfraBillingNodeCommand.TSQ_url,
    bodySchema: UpdateInfraBillingNodeCommand.RequestBodySchema,
    responseSchema: UpdateInfraBillingNodeCommand.ResponseSchema,
    requestMethod: UpdateInfraBillingNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Infra Billing Node updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Infra Billing Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
