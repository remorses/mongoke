// @jsx jsx
import { MiniGraphiQL } from 'mini-graphiql'

import { API_URL } from '../constants'
import { jsx, css } from '@emotion/core'
import { Stack, Box } from 'layout-kit-react'

export const Graphiql = ({ query }) => {
    return (
        <Stack w='100%' my='20px'>
            <Box w='100%' borderRadius='lg' overflow='hidden' borderWidth='1px'>
                <MiniGraphiQL url={API_URL} query={query} />
            </Box>
        </Stack>
    )
}

jsx
