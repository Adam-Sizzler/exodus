import { Center, Loader, Stack, Text, Title, Transition } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { TOAuth2ProvidersKeys } from '@exodus/backend-contract'
import { IconCheck } from '@tabler/icons-react'
import { CSSProperties, useEffect } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router'

import { useOauth2Callback } from '@shared/api/hooks'
import { ROUTES } from '@shared/constants'
import { useAuth } from '@shared/hooks/use-auth'
import { Page } from '@shared/ui/page'
import { consumeReturnTo } from '@shared/utils/return-to.util'

export const Oauth2CallbackPage = () => {
    const { provider } = useParams()
    const { setIsAuthenticated } = useAuth()

    const [searchParams] = useSearchParams()
    const navigate = useNavigate()

    const code = searchParams.get('code')
    const state = searchParams.get('state')

    const { mutate: oauth2Callback, isPending } = useOauth2Callback({
        mutationFns: {
            onSuccess: () => {
                setIsAuthenticated(true)

                navigate(consumeReturnTo() ?? ROUTES.DASHBOARD.HOME)
            },
            onError: (error) => {
                notifications.show({
                    title: 'OAuth2 Callback',
                    message: error.message,
                    color: 'red'
                })
                setIsAuthenticated(false)

                navigate(ROUTES.AUTH.LOGIN)
            }
        }
    })

    useEffect(() => {
        if (code && state && provider) {
            oauth2Callback({
                variables: {
                    provider: provider as TOAuth2ProvidersKeys,
                    code,
                    state
                }
            })
        }
    }, [code, state, provider])

    return (
        <Page title="OAuth2 Authentication">
            <Center style={{ minHeight: '60vh' }}>
                <Stack align="center" gap="xl">
                    <Transition
                        duration={300}
                        mounted={true}
                        timingFunction="ease"
                        transition="fade"
                    >
                        {(styles: CSSProperties) => (
                            <div style={styles}>
                                {isPending ? (
                                    <Loader size="xl" variant="dots" />
                                ) : (
                                    <IconCheck color="teal" size={48} />
                                )}
                            </div>
                        )}
                    </Transition>
                    <div>
                        <Title mb="md" order={2} ta="center">
                            Authenticating...
                        </Title>
                        <Text c="dimmed" ta="center">
                            Verifying credentials...
                        </Text>
                    </div>
                </Stack>
            </Center>
        </Page>
    )
}

export default Oauth2CallbackPage
