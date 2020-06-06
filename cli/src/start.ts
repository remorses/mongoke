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
        await startGokeServer({})
        // TODO start the server
    },
} // as CommandModule

export default command

interface StartArgs {}

async function startGokeServer(p: StartArgs) {
    const command_ = GOKE_START_COMMAND
    const args = [...command_.slice(1), ...makeStartArgs(p)]
    return await runCommand({
        command: command_[0],
        args,
    })
}

function makeStartArgs(p: StartArgs): string[] {
    return ['']
}
