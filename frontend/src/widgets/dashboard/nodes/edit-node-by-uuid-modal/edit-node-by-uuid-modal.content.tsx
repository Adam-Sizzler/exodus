import { UpdateNodeCommand } from '@exodus/backend-contract'
import { zodResolver } from 'mantine-form-zod-resolver'
import { useForm } from '@mantine/form'
import { motion } from 'motion/react'
import { useEffect } from 'react'

import {
    configProfilesQueryKeys,
    nodesQueryKeys,
    updateNodeFormSchema,
    UpdateNodeRequest,
    useGetNode,
    useGetNodePlugins,
    useGetPubKey,
    useUpdateNode
} from '@shared/api/hooks'
import { BaseNodeForm } from '@shared/ui/forms/nodes/base-node-form/base-node-form'
import { bytesToGbUtil, gbToBytesUtil } from '@shared/utils/bytes'
import { LoaderModalShared } from '@shared/ui/loader-modal'
import { queryClient } from '@shared/api'

import { NodeDetailsCardWidget } from '../node-details-card/node-details-card.widget'
import { NodeSystemCardWidget } from '../node-system-card/node-system-card.widget'

interface IProps {
    generatedCredentials?: {
        grpcToken?: string
        pubKey: string
    }
    nodeUuid: string
    onClose: () => void
}

export const EditNodeByUuidModalContent = (props: IProps) => {
    const { generatedCredentials, nodeUuid, onClose } = props

    const form = useForm<UpdateNodeRequest>({
        name: 'edit-node-form',
        mode: 'uncontrolled',
        validate: zodResolver(updateNodeFormSchema)
    })

    const { data: pubKey } = useGetPubKey({
        rQueryParams: {
            enabled: !generatedCredentials
        }
    })
    const { data: nodePlugins } = useGetNodePlugins()

    const { data: fetchedNode } = useGetNode({
        route: {
            uuid: nodeUuid
        },
        rQueryParams: {
            enabled: !form.isTouched()
        }
    })

    const { mutate: updateNode, isPending: isUpdateNodePending } = useUpdateNode({
        mutationFns: {
            onSuccess: async () => {
                queryClient.refetchQueries({
                    queryKey: nodesQueryKeys.getAllNodes.queryKey
                })
                queryClient.refetchQueries({
                    queryKey: configProfilesQueryKeys.getConfigProfiles.queryKey
                })
                queryClient.refetchQueries({
                    queryKey: nodesQueryKeys.getAllNodes.queryKey
                })
            }
        }
    })

    useEffect(() => {
        if (fetchedNode) {
            const normalizedFetchedSchema = (fetchedNode.apiSchema ?? '').toLowerCase()
            const apiSchema: 'mtls' | 'tls' = normalizedFetchedSchema === 'tls' ? 'tls' : 'mtls'
            form.initialize({
                uuid: fetchedNode.uuid,
                countryCode: fetchedNode.countryCode,
                name: fetchedNode.name,
                address: fetchedNode.address,
                port: fetchedNode.port ?? undefined,
                apiSchema,
                apiPath: fetchedNode.apiPath ?? '/',
                isTrafficTrackingActive: fetchedNode.isTrafficTrackingActive ?? undefined,
                trafficLimitBytes: bytesToGbUtil(fetchedNode.trafficLimitBytes ?? undefined),
                trafficResetDay: fetchedNode.trafficResetDay ?? undefined,
                notifyPercent: fetchedNode.notifyPercent ?? undefined,
                consumptionMultiplier: fetchedNode.consumptionMultiplier ?? undefined,
                tags: fetchedNode.tags ?? undefined,

                configProfile: {
                    activeConfigProfileUuid:
                        fetchedNode.configProfile.activeConfigProfileUuid ?? '',
                    activeInbounds:
                        fetchedNode.configProfile.activeInbounds.map((inbound) => inbound.uuid) ??
                        []
                },

                providerUuid: fetchedNode.providerUuid ?? undefined,
                activePluginUuid: fetchedNode.activePluginUuid ?? undefined
            })
        }
    }, [fetchedNode])

    const handleSubmit = form.onSubmit(async (values) => {
        if (!fetchedNode) {
            return
        }

        updateNode({
            variables: {
                ...values,
                name: values.name?.trim(),
                address: values.address?.trim(),
                apiSchema: values.apiSchema === 'tls' ? 'tls' : 'mtls',
                apiPath: values.apiPath?.trim() || '/',
                trafficLimitBytes: gbToBytesUtil(values.trafficLimitBytes),
                configProfile: {
                    activeConfigProfileUuid: values.configProfile?.activeConfigProfileUuid ?? '',
                    activeInbounds: values.configProfile?.activeInbounds ?? []
                }
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
            nodePlugins={nodePlugins?.nodePlugins ?? []}
            nodeSystemCard={<NodeSystemCardWidget node={fetchedNode} />}
            pubKey={
                (generatedCredentials ?? pubKey)
                    ? {
                          ...(generatedCredentials ?? pubKey)!,
                          grpcToken:
                              fetchedNode.grpcAuthToken ||
                              (generatedCredentials ?? pubKey)?.grpcToken ||
                              ''
                      }
                    : undefined
            }
        />
    )
}
