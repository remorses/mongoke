"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
exports.__esModule = true;
var react_landing_page_components_1 = require("react-landing-page-components");
var react_1 = __importDefault(require("react"));
var feather_1 = require("styled-icons/feather");
var codeStr = "\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\ncosa:\n    x: Str\n";
var App = function (_a) {
    var _b = _a.config, baseUrl = _b.baseUrl, docsUrl = _b.docsUrl;
    return (react_1["default"].createElement(react_landing_page_components_1.Provider, { color: 'rgb(15,52,74)', bg: '#eee', gradients: ['#ffeae8', '#f1efff'] },
        react_1["default"].createElement(react_landing_page_components_1.Hero, null,
            react_1["default"].createElement(react_landing_page_components_1.Logo, { width: ['100%', null, '800px'], src: baseUrl + 'img/mongoke.svg' }),
            react_1["default"].createElement(react_landing_page_components_1.Head, { fontSize: '60px' }, "Mongoke"),
            react_1["default"].createElement(react_landing_page_components_1.SubHead, null, "instant Graphql on MongoDb"),
            react_1["default"].createElement(react_landing_page_components_1.Button, null, "Get Started")),
        react_1["default"].createElement(react_landing_page_components_1.Section, null,
            react_1["default"].createElement(react_landing_page_components_1.Head, null, "Simple configuration"),
            react_1["default"].createElement(react_landing_page_components_1.Code, { width: ['400px', '800px'], language: 'yaml', code: codeStr })),
        react_1["default"].createElement(react_landing_page_components_1.Section, null,
            react_1["default"].createElement(react_landing_page_components_1.Head, null, "Cose"),
            react_1["default"].createElement(react_landing_page_components_1.SubHead, null, "The generated queries are super optimized. The generated queries are super optimized"),
            react_1["default"].createElement(react_landing_page_components_1.FeatureList, null,
                react_1["default"].createElement(react_landing_page_components_1.FeatureList.Feature, { icon: react_1["default"].createElement(feather_1.Archive, { width: '90px' }), title: 'Powerful queries', description: 'The generated queries are super optimized. The generated queries are super optimized' }),
                react_1["default"].createElement(react_landing_page_components_1.FeatureList.Feature, { icon: react_1["default"].createElement(feather_1.Airplay, { width: '90px' }), title: 'Write your db schema', description: 'prima cosa' }))),
        react_1["default"].createElement(react_landing_page_components_1.Section, null,
            react_1["default"].createElement(react_landing_page_components_1.Head, null, "How it Works"),
            react_1["default"].createElement(react_landing_page_components_1.Steps, null,
                react_1["default"].createElement(react_landing_page_components_1.Steps.Step, { icon: react_1["default"].createElement(feather_1.Archive, { width: '90px' }), title: 'Write your db schema', description: 'prima cosa' }),
                react_1["default"].createElement(react_landing_page_components_1.Steps.Step, { icon: react_1["default"].createElement(feather_1.Airplay, { width: '90px' }), title: 'Connect to your MongoDb', description: 'sec cosa' }),
                react_1["default"].createElement(react_landing_page_components_1.Steps.Step, { icon: react_1["default"].createElement(feather_1.Aperture, { width: '90px' }), title: 'Deploy with Docker', description: 'ultima cosa' }))),
        react_1["default"].createElement(react_landing_page_components_1.Section, null,
            react_1["default"].createElement(react_landing_page_components_1.Head, null, "Features"),
            react_1["default"].createElement(react_landing_page_components_1.Feature, { title: 'model', description: "\n                    Concerto lets you model the data used in your templates in a flexible and expressive way. \n                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.\n                    ", image: 
                // <img src='https://bemuse.ninja/project/img/screenshots/mode-selection.jpg' />
                react_1["default"].createElement(react_landing_page_components_1.Code, { light: true, language: 'yaml', code: codeStr }) }),
            react_1["default"].createElement(react_landing_page_components_1.Feature, { right: true, title: 'model', description: "\n                    Concerto lets you model the data used in your templates in a flexible and expressive way. \n                    Models can be written in a modular and portable way so they can be reused in a variety of contracts.\n                    ", 
                // image={<img  src='https://developer.cohesity.com/img/python.png'/>}
                image: react_1["default"].createElement(feather_1.Airplay, null) }))));
};
// render(<App />, document.getElementById('root'))
// export default App
module.exports = App;
// // @ts-ignore
// if (module.hot) {
//     // @ts-ignore
//     module.hot.accept()
// }
