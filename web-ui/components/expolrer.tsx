import React, { Component, useEffect, useState, useRef, Fragment } from 'react'
import GraphiQL from 'graphiql'
import GraphiQLExplorer from 'graphiql-explorer'
import { buildClientSchema, getIntrospectionQuery, parse } from 'graphql'


import { makeDefaultArg, getDefaultScalarArgValue } from './customArgs'

import type { GraphQLSchema } from 'graphql'
import { Box } from '@chakra-ui/core'


const DEFAULT_QUERY = `

`

type State = {
    schema?: GraphQLSchema
    query: string
    explorerIsOpen: boolean
}

export const Explorer = (props) => {

    const {...rest} = props

    const _graphiql = useRef<GraphiQL>()
    const [state, setState] = useState({
        schema: null,
        query: DEFAULT_QUERY,
        explorerIsOpen: true,
    })

    const url = 'http://localhsot:8080'

    useEffect(() => {
        fetcher(url, {
            query: getIntrospectionQuery(),
        }).then((result) => {
            const editor = _graphiql.current.getQueryEditor()
            editor.setOption('extraKeys', {
                ...(editor.options.extraKeys || {}),
                'Shift-Alt-LeftClick': _handleInspectOperation,
            })

            setState((x) => ({ ...x, schema: buildClientSchema(result.data) }))
        })
    }, [])

    const _handleInspectOperation = (
        cm: any,
        mousePos: { line: Number; ch: Number },
    ) => {
        const parsedQuery = parse(state.query || '')

        if (!parsedQuery) {
            console.error("Couldn't parse query document")
            return null
        }

        var token = cm.getTokenAt(mousePos)
        var start = { line: mousePos.line, ch: token.start }
        var end = { line: mousePos.line, ch: token.end }
        var relevantMousePos = {
            start: cm.indexFromPos(start),
            end: cm.indexFromPos(end),
        }

        var position = relevantMousePos

        var def = parsedQuery.definitions.find((definition) => {
            if (!definition.loc) {
                console.log('Missing location information for definition')
                return false
            }

            const { start, end } = definition.loc
            return start <= position.start && end >= position.end
        })

        if (!def) {
            console.error(
                'Unable to find definition corresponding to mouse position',
            )
            return null
        }

        var operationKind =
            def.kind === 'OperationDefinition'
                ? def.operation
                : def.kind === 'FragmentDefinition'
                ? 'fragment'
                : 'unknown'

        var operationName =
            def.kind === 'OperationDefinition' && !!def.name
                ? def.name.value
                : def.kind === 'FragmentDefinition' && !!def.name
                ? def.name.value
                : 'unknown'

        var selector = `.graphiql-explorer-root #${operationKind}-${operationName}`

        var el = document.querySelector(selector)
        el && el.scrollIntoView()
    }

    const _handleEditQuery = (query: string): void =>
        setState((x) => ({ ...x, query }))

    const _handleToggleExplorer = () => {
        setState((x) => ({ ...x, explorerIsOpen: !state.explorerIsOpen }))
    }

    
    return (
        <Box {...rest}>
            <GraphiQLExplorer
                schema={state.schema}
                query={state.query}
                onEdit={_handleEditQuery}
                onRunOperation={(operationName) =>
                    _graphiql.current.handleRunQuery(operationName)
                }
                explorerIsOpen={state.explorerIsOpen}
                onToggleExplorer={_handleToggleExplorer}
                getDefaultScalarArgValue={getDefaultScalarArgValue}
                makeDefaultArg={makeDefaultArg}
            />
            <GraphiQL
                ref={_graphiql}
                fetcher={fetcher}
                schema={state.schema}
                query={state.query}
                onEditQuery={_handleEditQuery}
            >
                <GraphiQL.Toolbar>
                    <GraphiQL.Button
                        onClick={() => _graphiql.current.handlePrettifyQuery()}
                        label='Prettify'
                        title='Prettify Query (Shift-Ctrl-P)'
                    />
                    <GraphiQL.Button
                        onClick={() => _graphiql.current.handleToggleHistory()}
                        label='History'
                        title='Show History'
                    />
                    <GraphiQL.Button
                        onClick={_handleToggleExplorer}
                        label='Explorer'
                        title='Toggle Explorer'
                    />
                </GraphiQL.Toolbar>
            </GraphiQL>
            </Box>
    )
}



function fetcher(url, params: Object) {
    return fetch(url, {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(params),
    })
        .then(function (response) {
            return response.text()
        })
        .then(function (responseBody) {
            try {
                return JSON.parse(responseBody)
            } catch (e) {
                return responseBody
            }
        })
}