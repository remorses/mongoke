"use strict";
exports.__esModule = true;
var src_1 = require("../src");
var react_1 = require("react");
var Text_1 = require("../src/Text");
var feather_1 = require("styled-icons/feather");
var codeStr = "\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\n";
var App = function () {
    return (<src_1.Provider color='rgb(15,52,74)' bg='#eee' gradients={['#ffeae8', '#f1efff',]}>
            <src_1.Hero>
                <src_1.Logo width={['100%', null, '800px']} src={require('./mongoke.svg')}/>
                <Text_1.Head fontSize='60px'>Mongoke</Text_1.Head>
                <Text_1.SubHead>instant Graphql on MongoDb</Text_1.SubHead>
                <src_1.Button>Get Started</src_1.Button>
            </src_1.Hero>
            <src_1.Section>
                <Text_1.Head>Simple configuration</Text_1.Head>
                <src_1.Code width={['400px', '800px']} language='yaml' code={codeStr}/>
            </src_1.Section>
            <src_1.Section>
                <Text_1.Head>Cose</Text_1.Head>
                <Text_1.SubHead>The generated queries are super optimized. The generated queries are super optimized</Text_1.SubHead>
                <src_1.FeatureList>
                    <src_1.FeatureList.Feature icon={<feather_1.Archive width='90px'/>} title='Powerful queries' description='The generated queries are super optimized. The generated queries are super optimized'/>
                    <src_1.FeatureList.Feature icon={<feather_1.Airplay width='90px'/>} title='Write your db schema' description='prima cosa'/>
                    
                </src_1.FeatureList>
            </src_1.Section>
            <src_1.Section>
                <Text_1.Head>How it Works</Text_1.Head>
                <src_1.Steps>
                    <src_1.Steps.Step icon={<feather_1.Archive width='90px'/>} title='Write your db schema' description='prima cosa'/>
                    <src_1.Steps.Step icon={<feather_1.Airplay width='90px'/>} title='Connect to your MongoDb' description='sec cosa'/>
                    <src_1.Steps.Step icon={<feather_1.Aperture width='90px'/>} title='Deploy with Docker' description='ultima cosa'/>
                </src_1.Steps>
            </src_1.Section>
            <src_1.Section>
                <Text_1.Head>Features</Text_1.Head>
                <src_1.Feature title='model' description={"\n                    Concerto lets you model the data used in your templates in a flexible and expressive way. \n                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.\n                    "} image={
    // <img src='https://bemuse.ninja/project/img/screenshots/mode-selection.jpg' />
    <src_1.Code light language='yaml' code={codeStr}/>}/>
                <src_1.Feature right title='model' description={"\n                    Concerto lets you model the data used in your templates in a flexible and expressive way. \n                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.\n                    "} 
    // image={<img  src='https://developer.cohesity.com/img/python.png'/>}
    image={<feather_1.Airplay />}/>
            </src_1.Section>
        </src_1.Provider>);
};
// render(<App />, document.getElementById('root'))
exports["default"] = App;
// // @ts-ignore
// if (module.hot) {
//     // @ts-ignore
//     module.hot.accept()
// }
