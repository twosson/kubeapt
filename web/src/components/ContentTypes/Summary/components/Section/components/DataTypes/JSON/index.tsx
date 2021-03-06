import React from 'react'
import './styles.scss'

interface Props {
  params: ContentType;
}

export default function Item(props: Props) {
  const { params } = props
  const { label, data } = params
  return (
    <div className='summary--data summary--data-json'>
      <div className='summary--data-key'>{label}</div>
      <div className='summary--data-json-value'>
        {JSON.stringify(data, null, 2)}
      </div>
    </div>
  )
}
