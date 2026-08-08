import { ConnectionsByNodeResultCommand } from '@exodus/backend-contract'
import { useCallback, useEffect, useState } from 'react'

import { useConnectionsByNode, useConnectionsByNodeResult } from '@shared/api/hooks'

export type ActiveSessionUser = NonNullable<
    ConnectionsByNodeResultCommand.Response['response']['result']
>['users'][number]

export const useNodeActiveSessions = (nodeUuid: string) => {
    const [jobId, setJobId] = useState<null | string>(null)
    const [isCompleted, setIsCompleted] = useState(false)
    const [isFailed, setIsFailed] = useState(false)

    const { mutate: connectionsByNode } = useConnectionsByNode({
        route: {
            nodeUuid
        },
        mutationFns: {
            onSuccess: (data) => {
                setJobId(data.jobId)
            }
        }
    })

    const shouldPoll = !!jobId && !isCompleted && !isFailed

    const { data: connectionsByNodeResult } = useConnectionsByNodeResult({
        route: {
            jobId: jobId ?? ''
        },
        rQueryParams: {
            enabled: shouldPoll,
            refetchInterval: shouldPoll ? 1000 : false
        }
    })

    useEffect(() => {
        if (!connectionsByNodeResult) return undefined
        if (
            connectionsByNodeResult.isFailed ||
            (connectionsByNodeResult.isCompleted && !connectionsByNodeResult.result?.success)
        ) {
            // oxlint-disable-next-line
            setIsFailed(true)
            return undefined
        }
        if (connectionsByNodeResult.isCompleted) {
            const timer = setTimeout(() => setIsCompleted(true), 500)
            return () => clearTimeout(timer)
        }
        return undefined
    }, [connectionsByNodeResult])

    const refresh = useCallback(() => {
        setJobId(null)
        setIsCompleted(false)
        setIsFailed(false)
        connectionsByNode({})
    }, [connectionsByNode])

    useEffect(() => {
        connectionsByNode({})
    }, [])

    return {
        isCompleted,
        isFailed,
        refresh,
        users: connectionsByNodeResult?.result?.users
    }
}
