import { DokzProvider, GithubLink, ColorModeSwitch } from 'dokz/dist'
import React from 'react'
import 'mini-graphiql/dist/style.css'

export default function App(props) {
    const { Component, pageProps } = props
    return (
        <DokzProvider
            headerItems={[
                <GithubLink key='0' url='https://github.com/remorses/dokz' />,
                <ColorModeSwitch key='1' />,
            ]}
            docsRootPath='pages/docs'
            sidebarOrdering={{
                'index.mdx': true,
                Documents_Group: {
                    'another.mdx': true,
                },
            }}
        >
            <Component {...pageProps} />
        </DokzProvider>
    )
}
