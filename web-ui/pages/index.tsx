import GraphiQL from 'graphiql'
import { Box } from '@chakra-ui/core'



const API_URL = process.env.NODE_ENV == "production" ? "/" : "http://localhost:8080"



const Page = (props) => {
    return (
        <Box height='100%'>
            <GraphiQL
                fetcher={async (graphQLParams) => {
                    const data = await fetch(
                        API_URL,
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
