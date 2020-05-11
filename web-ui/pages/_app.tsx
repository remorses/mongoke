import 'graphiql/graphiql.min.css'
import '../index.css'
import { ThemeProvider, CSSReset, Box } from '@chakra-ui/core'

function MyApp({ Component, pageProps }) {
    return (
        <Box h='100vh'>
            <ThemeProvider>
                <CSSReset />
                <Component {...pageProps} />
            </ThemeProvider>
        </Box>
    )
}

export default MyApp
