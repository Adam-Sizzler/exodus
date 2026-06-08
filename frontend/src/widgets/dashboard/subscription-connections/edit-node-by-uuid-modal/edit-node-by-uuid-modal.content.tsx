import { zodResolver } from 'mantine-form-zod-resolver'
import { useEffect } from 'react'
import { useForm } from '@mantine/form'
import { motion } from 'motion/react'

import {
    subscriptionConnectionsQueryKeys,
    updateSubscriptionConnectionFormSchema,
    useGetSubscriptionConnection,
    useGetSubscriptionConnectionsPubKey,
    useUpdateSubscriptionConnection
} from '@shared/api/hooks'
import { BaseNodeForm } from '@shared/ui/forms/subscription-connections/base-subscription-connection-form/base-subscription-connection-form'
import { LoaderModalShared } from '@shared/ui/loader-modal'
import { queryClient } from '@shared/api'

import { NodeDetailsCardWidget } from '../node-details-card/node-details-card.widget'

interface IProps {
    generatedCredentials?: {
        grpcToken?: string
        pubKey: string
    }
    nodeUuid: string
    onClose: () => void
}

type EditSubscriptionConnectionForm = {
    uuid: string
    name?: string
    address?: string
    publicDomain?: string | null
    port?: number
    apiSchema?: 'mtls' | 'tls'
    apiPath?: string
    subpageConfigUuid?: string | null
    providerUuid?: string | null
    tags?: string[]
}

export const EditNodeByUuidModalContent = (props: IProps) => {
    const { generatedCredentials, nodeUuid, onClose } = props

    const form = useForm<EditSubscriptionConnectionForm>({
        name: 'edit-subscription-connection-form',
        mode: 'uncontrolled',
        validate: zodResolver(updateSubscriptionConnectionFormSchema)
    })

    const { data: pubKey } = useGetSubscriptionConnectionsPubKey({
        rQueryParams: {
            enabled: !generatedCredentials
        }
    })

    const { data: fetchedNode } = useGetSubscriptionConnection({
        route: {
            uuid: nodeUuid
        },
        rQueryParams: {
            enabled: !form.isTouched()
        }
    })

    const { mutate: updateNode, isPending: isUpdateNodePending } = useUpdateSubscriptionConnection({
        mutationFns: {
            onSuccess: async () => {
                form.resetDirty()
                queryClient.refetchQueries({
                    queryKey: subscriptionConnectionsQueryKeys.getAllNodes.queryKey
                })
                queryClient.refetchQueries({
                    queryKey: subscriptionConnectionsQueryKeys.getNode({ uuid: nodeUuid }).queryKey
                })
            }
        }
    })

    useEffect(() => {
        if (fetchedNode) {
            const normalizedFetchedSchema = (fetchedNode.apiSchema ?? '').toLowerCase()
            const apiSchema: 'mtls' | 'tls' =
                normalizedFetchedSchema === 'tls' ? 'tls' : 'mtls'
            form.initialize({
                uuid: fetchedNode.uuid,
                name: fetchedNode.name,
                address: fetchedNode.address,
                publicDomain: fetchedNode.publicDomain ?? null,
                port: fetchedNode.port ?? undefined,
                apiSchema,
                apiPath: fetchedNode.apiPath ?? '/',
                subpageConfigUuid: fetchedNode.subpageConfigUuid ?? null,
                tags: fetchedNode.tags ?? undefined,
                providerUuid: fetchedNode.providerUuid ?? null
            })
        }
    }, [fetchedNode])

    const handleSubmit = form.onSubmit(async (values) => {
        if (!fetchedNode) {
            return
        }

        const schema = values.apiSchema === 'tls' ? 'tls' : 'mtls'
        const resolvedSubpageConfigUuid =
            values.subpageConfigUuid === undefined
                ? (fetchedNode.subpageConfigUuid ?? null)
                : ((values.subpageConfigUuid ?? '').trim() || null)
        const resolvedPublicDomain =
            values.publicDomain === undefined
                ? (fetchedNode.publicDomain ?? null)
                : ((values.publicDomain ?? '').trim() || null)

        updateNode({
            variables: {
                ...values,
                uuid: fetchedNode.uuid,
                name: values.name?.trim(),
                address: values.address?.trim(),
                publicDomain: resolvedPublicDomain,
                apiSchema: schema,
                apiPath: values.apiPath?.trim() || '/',
                subpageConfigUuid: resolvedSubpageConfigUuid
            }
        })
    })

    if (!fetchedNode) {
        return (
            <motion.div
                animate={{ opacity: 1 }}
                initial={{ opacity: 0 }}
                transition={{ duration: 0.3 }}
            >
                <LoaderModalShared h="78vh" />
            </motion.div>
        )
    }

    return (
        <BaseNodeForm
            form={form}
            handleClose={onClose}
            handleSubmit={handleSubmit}
            isDataSubmitting={isUpdateNodePending}
            node={fetchedNode}
            nodeDetailsCard={<NodeDetailsCardWidget node={fetchedNode} />}
            pubKey={
                (generatedCredentials ?? pubKey)
                    ? {
                          ...(generatedCredentials ?? pubKey)!,
                          grpcToken:
                              fetchedNode.grpcAuthToken ||
                              (generatedCredentials ?? pubKey)?.grpcToken
                      }
                    : undefined
            }
        />
    )
}
