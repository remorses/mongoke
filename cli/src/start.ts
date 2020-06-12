import { runCommand } from './support'
import { GOKE_START_COMMAND, WEB_UI_ASSETS } from './contsants'
import { CommandModule } from 'yargs'

const command: CommandModule = {
    command: ['start', '*'],
    describe: 'Starts goke server',
    builder: (argv) => {
        argv.option('github-token', {
            type: 'string',
            default: '',
            required: false,
            description:
                'The github token to use for login, instead of using the browser',
        })
        return argv
    },
    handler: async (argv) => {
        console.log('starting the server')
        // console.log(argv.$0)
        // console.log(process.argv)
        await startGokeServer(process.argv.slice(process.env.DEBUG ? 2 : 1)) // TODO this is brittle
    },
} // as CommandModule

export default command

async function startGokeServer(a) {
    const command = GOKE_START_COMMAND
    if (process.env.DEBUG) {
        console.log([...command, ...a])
    }
    const args = [...command.slice(1), ...a]
    return await runCommand({
        command: command[0],
        args,
        env: {
            WEB_UI_ASSETS,
        },
    }).catch((e) => {
        throw new Error('could not start goke server')
    })
}
