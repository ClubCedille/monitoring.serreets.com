import {React, useState} from 'react'
import {Form, Button} from 'react-bootstrap'
import FormContainer from '../components/FormContainer'

const SignupPage = ({history}) => {

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const submitHandler = async (e) => {
    e.preventDefault()

    await fetch(process.env.REACT_APP_URL + 'api/signup', {
      method: 'POST',
      headers: {'Content-Type': 'application-json'},
      body: JSON.stringify({
        email,
        password
      })
    })

    history.push('/login')
  }

  return (
    <FormContainer>
      <h1>Signup page</h1>
      <Form onSubmit={submitHandler}>
        <Form.Group className="mb-3" controlId="email">
          <Form.Label>Email</Form.Label>
          <Form.Control type="email" placeholder="Enter your email" 
          value={email}
          onChange={e => setEmail(e.target.value)}/>
        </Form.Group>

        <Form.Group className="mb-3" controlId="password">
          <Form.Label>Password</Form.Label>
          <Form.Control type="password" placeholder="Password" 
          value={password}
          onChange={e => setPassword(e.target.value)}/>
        </Form.Group>

        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
    </FormContainer>   
  )
}

export default SignupPage
