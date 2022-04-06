import React from 'react'

const HomePage = ({email}) => {
  return (
    email ? <h1>Welcome {email}</h1> : 
    <div>
      Welcome to the home page!
    </div>
  )
}

export default HomePage
