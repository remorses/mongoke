import path from 'path'

export const GOKE_START_COMMAND = process.env.DEV_START_COMMAND.split(' ') || [
    'goke-server',
]

const root = path.dirname(path.resolve(path.join('..', 'package.json')))

export const WEB_UI_ASSETS = path.join(root, 'dist', 'web-ui')
