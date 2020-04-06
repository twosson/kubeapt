import React, { Component } from 'react'
import { Switch, Route } from 'react-router-dom'
import Promise from 'promise'
import { getNavigation, getNamespaces } from 'api'
import Home from 'pages/Home'
import Header from '../Header'
import Navigation from '../Navigation'
import './styles.scss'

class App extends Component {
  constructor (props) {
    super(props)
    this.state = {
      navigation: [],
      namespaceOptions: [],
      namespaceValue: null
    }
  }

  async componentDidMount () {
    const [navigation, namespacesPayload] = await Promise.all([
      getNavigation(),
      getNamespaces()
    ])
    let namespaceOptions = []
    if (
      namespacesPayload &&
      namespacesPayload.namespaces &&
      namespacesPayload.namespaces.length
    ) {
      namespaceOptions = namespacesPayload.namespaces.map(ns => ({
        label: ns,
        value: ns
      }))
    }

    this.setState({ navigation, namespaceOptions })
  }

  render () {
    const { navigation, namespaceOptions, namespaceValue } = this.state
    return (
      <div className='app'>
        <Header />
        <div className='app-page'>
          <div className='app-nav'>
            <Navigation
              navigation={navigation}
              namespaceOptions={namespaceOptions}
              namespaceValue={namespaceValue}
              onNamespaceChange={value => this.setState({ namespaceValue: value })
              }
            />
          </div>
          <div className='app-main'>
            <Switch>
              <Route path='/' component={Home} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}

export default App
