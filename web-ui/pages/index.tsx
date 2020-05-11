import GraphiQL from 'graphiql'
import { Box } from '@chakra-ui/core'

const Page = (props) => {
    return (
        <Box height='100%'>
            <GraphiQL
                fetcher={async (graphQLParams) => {
                    const data = await fetch(
                        'https://swapi-graphql.netlify.app/.netlify/functions/index',
                        {
                            method: 'POST',
                            headers: {
                                Accept: 'application/json',
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify(graphQLParams),
                            credentials: 'same-origin',
                        },
                    )
                    return data.json().catch(() => data.text())
                }}
            />
        </Box>
    )
}

export default Page
