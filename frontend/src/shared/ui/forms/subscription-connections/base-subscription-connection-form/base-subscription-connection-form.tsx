import { Button, CopyButton, Menu, px, SimpleGrid, Stack } from '@mantine/core'
import { PiFloppyDiskDuotone } from 'react-icons/pi'
import { UseFormReturnType } from '@mantine/form'
import { TbCopy, TbDots } from 'react-icons/tb'
import { motion } from 'framer-motion'
import { ReactNode } from 'react'
import { t } from 'i18next'

import { ToggleNodeStatusButtonFeature } from '@features/ui/dashboard/subscription-connections/toggle-node-status-button'
import { RestartNodeButtonFeature } from '@features/ui/dashboard/subscription-connections/restart-node-button'
import { DeleteNodeFeature } from '@features/ui/dashboard/subscription-connections/delete-node'
import { ModalAccordionWidget } from '@widgets/dashboard/subscription-connections/modal-accordeon-widget'
import { SubscriptionConnectionResponse } from '@shared/api/hooks'
import { ModalFooter } from '@shared/ui/modal-footer'

import { NodeVitalsCard } from './node-vitals.card'
import { SubscriptionConfigCard } from './subscription-config.card'

const MotionWrapper = motion.div
const MotionStack = motion.create(Stack)

const containerVariants = {
    hidden: {},
    visible: {
        transition: {
            staggerChildren: 0.1
        }
    }
}

const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
        opacity: 1,
        y: 0,
        transition: { duration: 0.3 }
    }
}

interface BaseSubscriptionConnectionFormValues {
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

interface SubscriptionConnectionKeygenResponse {
    pubKey: string
    grpcToken?: string
}

interface IProps<T extends BaseSubscriptionConnectionFormValues> {
    form: UseFormReturnType<T>
    handleClose: () => void
    handleSubmit: () => void
    isDataSubmitting: boolean
    node: SubscriptionConnectionResponse
    nodeDetailsCard?: ReactNode
    pubKey: SubscriptionConnectionKeygenResponse | undefined
}

export const BaseNodeForm = <T extends BaseSubscriptionConnectionFormValues>(props: IProps<T>) => {
    const { form, node, pubKey, nodeDetailsCard, handleClose, handleSubmit, isDataSubmitting } = props

    return (
        <>
            <MotionStack animate="visible" gap="md" initial="hidden" variants={containerVariants}>
                <ModalAccordionWidget node={node} />
                <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md">
                    <Stack gap="md">
                        {nodeDetailsCard && (
                            <MotionWrapper variants={cardVariants}>{nodeDetailsCard}</MotionWrapper>
                        )}
                        <SubscriptionConfigCard
                            cardVariants={cardVariants}
                            form={form}
                            motionWrapper={MotionWrapper}
                        />
                    </Stack>

                    <Stack gap="md">
                        <NodeVitalsCard
                            cardVariants={cardVariants}
                            form={form}
                            motionWrapper={MotionWrapper}
                            pubKey={pubKey}
                        />
                    </Stack>
                </SimpleGrid>
            </MotionStack>

            <ModalFooter>
                {node && (
                    <Menu keepMounted={true} position="top-end" shadow="md">
                        <Menu.Target>
                            <Button
                                color="gray"
                                leftSection={<TbDots size={px('1.2rem')} />}
                                size="md"
                            >
                                {t('base-node-form.more-actions')}
                            </Button>
                        </Menu.Target>

                        <Menu.Dropdown>
                            <DeleteNodeFeature handleClose={handleClose} node={node} />
                            <Menu.Divider />

                            <Menu.Label>{t('base-node-form.management')}</Menu.Label>
                            <CopyButton value={node.uuid}>
                                {({ copy }) => (
                                    <Menu.Item leftSection={<TbCopy size="16px" />} onClick={copy}>
                                        {t('common.copy-uuid')}
                                    </Menu.Item>
                                )}
                            </CopyButton>

                            <RestartNodeButtonFeature handleClose={handleClose} node={node} />
                            <ToggleNodeStatusButtonFeature handleClose={handleClose} node={node} />
                        </Menu.Dropdown>
                    </Menu>
                )}
                <Button
                    color="teal"
                    disabled={!form.isDirty()}
                    leftSection={<PiFloppyDiskDuotone size="16px" />}
                    loading={isDataSubmitting}
                    onClick={handleSubmit}
                    size="md"
                    variant="light"
                >
                    {t('common.save')}
                </Button>
            </ModalFooter>
        </>
    )
}
