#!/usr/bin/env node
import yargs from 'yargs'
import startCommand from './start'
import { withErrorHandling } from './support'

yargs
    .option('file', {
        type: 'string',
        alias: 'f',
        default: 'docker-compose.yml',
    })
    .option('verbose', {
        alias: 'v',
        type: 'boolean',
        default: false,
    })
    .option('env-file', {
        type: 'string',
        default: '.env',
    })
    // .middleware([
    //     (argv) => {
    //         if (argv.verbose) {
    //             winston.configure({
    //                 ...winstonConf,
    //                 level: 'debug',
    //             })
    //             return
    //         }
    //         winston.configure({ ...winstonConf, silent: true, level: 'error' })
    //     },
    // ])
    .command(withErrorHandling(startCommand))
    // .demandCommand()
    .help('h').argv
