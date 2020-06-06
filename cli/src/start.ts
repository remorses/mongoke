import { runCommand } from './support'
import { GOKE_START_COMMAND } from './contsants'
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
        // TODO start the server
    },
} // as CommandModule

export default command

async function startGokeServer(a) {
    const command_ = GOKE_START_COMMAND
    console.log([...command_, ...a])
    const args = [...command_.slice(1), ...a]
    return await runCommand({
        command: command_[0],
        args,
    })
}
