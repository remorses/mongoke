import {
    Hero,
    Section,
    Steps,
    Feature,
    Logo,
    Button,
    Provider,
    Code,
    FeatureList
} from '../src'
import React from 'react'
import { H1, Image, Text, Box } from 'hybrid-components'
import { Head, SubHead } from '../src/Text'
import { Archive, Airplay, Aperture } from 'styled-icons/feather'

const codeStr = `
cosa:
    x: Str
cosa:
    x: Str
cosa:
    x: Str
cosa:
    x: Str
cosa:
    x: Str
`

const App = () => {
    return (
        <Provider color='rgb(15,52,74)' bg='#eee' gradients={['#ffeae8', '#f1efff',]}>
            <Hero>
                <Logo width={['100%', null, '800px']} src={require('./mongoke.svg')} />
                <Head fontSize='60px'>Mongoke</Head>
                <SubHead>instant Graphql on MongoDb</SubHead>
                <Button>Get Started</Button>
            </Hero>
            <Section>
                <Head>Simple configuration</Head>
                <Code width={['400px', '800px']} language='yaml' code={codeStr} />
            </Section>
            <Section>
                <Head>Cose</Head>
                <SubHead>The generated queries are super optimized. The generated queries are super optimized</SubHead>
                <FeatureList>
                    <FeatureList.Feature
                        icon={<Archive width='90px' />}
                        title='Powerful queries'
                        description='The generated queries are super optimized. The generated queries are super optimized'
                    />
                    <FeatureList.Feature
                        icon={<Airplay width='90px' />}
                        title='Write your db schema'
                        description='prima cosa'
                    />
                    
                </FeatureList>
            </Section>
            <Section>
                <Head>How it Works</Head>
                <Steps>
                    <Steps.Step
                        icon={<Archive width='90px' />}
                        title='Write your db schema'
                        description='prima cosa'
                    />
                    <Steps.Step
                        icon={<Airplay width='90px' />}
                        title='Connect to your MongoDb'
                        description='sec cosa'
                    />
                    <Steps.Step
                        icon={<Aperture width='90px' />}
                        title='Deploy with Docker'
                        description='ultima cosa'
                    />
                </Steps>
            </Section>
            <Section>
                <Head>Features</Head>
                <Feature
                    title='model'
                    description={`
                    Concerto lets you model the data used in your templates in a flexible and expressive way. 
                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.
                    `}
                    image={
                        // <img src='https://bemuse.ninja/project/img/screenshots/mode-selection.jpg' />
                        <Code light language='yaml' code={codeStr} />
                    }
                />
                <Feature
                    right
                    title='model'
                    description={`
                    Concerto lets you model the data used in your templates in a flexible and expressive way. 
                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.
                    `}
                    // image={<img  src='https://developer.cohesity.com/img/python.png'/>}
                    image={<Airplay />}
                />
            </Section>
        </Provider>
    )
}

// render(<App />, document.getElementById('root'))
export default App

// // @ts-ignore
// if (module.hot) {
//     // @ts-ignore
//     module.hot.accept()
// }
