const footer = {
    style: 'dark',
    links: [
        {
            title: 'Docs',
            items: [
                {
                    label: 'Style Guide',
                    to: 'docs/doc1',
                },
                {
                    label: 'Second Doc',
                    to: 'docs/doc2',
                },
            ],
        },
        {
            title: 'Community',
            items: [
                {
                    label: 'Stack Overflow',
                    href:
                        'https://stackoverflow.com/questions/tagged/docusaurus',
                },
                {
                    label: 'Discord',
                    href: 'https://discordapp.com/invite/docusaurus',
                },
            ],
        },
        {
            title: 'Social',
            items: [
                {
                    label: 'Blog',
                    to: 'blog',
                },
                {
                    label: 'GitHub',
                    href: 'https://github.com/facebook/docusaurus',
                },
                {
                    label: 'Twitter',
                    href: 'https://twitter.com/docusaurus',
                },
            ],
        },
    ],
    copyright: `Copyright © ${new Date().getFullYear()} My Project, Inc. Built with Docusaurus.`,
}

const themeConfig = {
    navbar: {
        title: 'Mongoke',
        logo: {
            alt: 'Mongoke',
            src: 'img/logo.svg',
        },
        links: [
            {
                to: 'docs/quickstart',
                activeBasePath: 'docs',
                label: 'Docs',
                position: 'left',
            },
            // { to: 'blog', label: 'Blog', position: 'left' },
            {
                href: 'https://github.com/remorses/mongoke',
                label: 'GitHub',
                position: 'right',
            },
        ],
    },
    footer,
}

module.exports = {
    title: 'mongoke',
    tagline: 'Instant grphalq for Mongodb',
    url: 'https://mongoke.club',
    baseUrl: '/',
    favicon: 'img/favicon.ico',
    organizationName: 'mongoke', // Usually your GitHub org/user name.
    projectName: 'mongoke', // Usually your repo name.
    themeConfig,
    presets: [
        [
            '@docusaurus/preset-classic',
            {
                docs: {
                    sidebarPath: require.resolve('./sidebars.js'),
                    editUrl:
                        'https://github.com/facebook/docusaurus/edit/master/website/',
                },
                theme: {
                    customCss: require.resolve('./src/css/custom.css'),
                },
            },
        ],
    ],
}
