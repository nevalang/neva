import * as React from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"

import { ComponentTypes, Program } from "../types/program"
import { Api } from "../api"
import { useEffect, useState } from "react"
import { ProgramEditor } from "./program"
import { Menu } from "./menu"
import { Redirect } from "react-router"
// import {drag/} from 'reaflow'
import {
  Breadcrumb,
  Breadcrumbs,
  Icon,
  Intent,
  Spinner,
  SpinnerSize,
} from "@blueprintjs/core"
import { BREADCRUMBS } from "@blueprintjs/core/lib/esm/common/classes"
import { Palette } from "./palette"

const defaultProgram: Program = {
  root: "",
  scope: {
    "": {
      type: ComponentTypes.MODULE,
      constants: {},
      io: {
        in: {},
        out: {},
      },
      deps: {},
      net: [],
      workers: {},
    },
  },
}

interface AppProps {
  api: Api
}

function App(props: AppProps) {
  const [path, setPath] = useState("examples/program/pkg")
  const [program, setProgram] = useState(defaultProgram)
  const [err, setErr] = useState(null)

  useEffect(() => {
    async function aux() {
      let err = null
      try {
        setProgram(await props.api.getProgram(path))
      } catch (err) {
        err = err
      } finally {
        setErr(err)
      }
    }

    aux()
  }, [path])

  const handleRemoveFromScope = async (name: string) => {
    const filtered = Object.entries(program.scope).filter(([k]) => k !== name)
    const newProgram = {
      root: program.root,
      scope: Object.fromEntries(filtered),
    }
    try {
      await props.api.editProgram(path, newProgram)
      setProgram(newProgram)
      setErr(null)
    } catch (err) {
      setErr(err)
    }
  }

  if (err !== null) {
    return <span>{err}</span>
  }

  return (
    <Router>
      <Breadcrumbs
        currentBreadcrumbRenderer={props => (
          <Breadcrumb {...props}>
            text <Icon icon="star" />
          </Breadcrumb>
        )}
        items={[]}
      />
      <Switch>
        <Redirect exact from="/" to="/menu" />
        <Route path="/menu" component={Menu} exact />
        <Route
          path="/program"
          component={props => (
            <ProgramEditor
              {...props}
              program={program}
              onAddToScope={console.log}
              onRemoveFromScope={handleRemoveFromScope}
            />
          )}
        />
        {/* <Palette onClick={console.log} scope={state.program.scope} /> */}
      </Switch>
    </Router>
  )
}

export { App }
