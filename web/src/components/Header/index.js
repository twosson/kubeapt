import React from 'react'
import AptLogoHorizontal from '../AptLogoHorizontal'
import './styles.scss'

function Header () {
  return (
    <header className='header'>
      <AptLogoHorizontal className='header--heptio-logo' />
    </header>
  )
}

export default Header
