import React from 'react'
import './styles.scss'

export default function ({ emptyContent }) {
  return (
    <div className='content-empty'>
      {emptyContent}
    </div>
  )
}
