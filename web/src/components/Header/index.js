import React from 'react'
import AptLogoHorizontal from '../AptLogoHorizontal'
import './styles.scss'

function Header () {
  return (
    <header className='header'>
      <AptLogoHorizontal className='header--apt-logo' />
    </header>
  )
}

export default Header
