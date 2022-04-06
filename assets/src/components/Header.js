import {React} from 'react'
import {Navbar, Nav, Container} from 'react-bootstrap'

const Header = ({email, setEmail}) => {

  const logoutHandler = async (e) => {
    e.preventDefault()
    await fetch('http://localhost:3001/api/logout', {
      method: 'GET',
      headers: {'Content-Type': 'application-json'},
      credentials: 'include',
    })

    setEmail('')
  }

  return (
    <Navbar bg="dark" variant='dark' collapseOnSelect expand="lg">
      <Container>
        <Navbar.Brand href="/">Auth</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          {email ? (
            <Nav className="ms-auto">
              <Nav.Link onClick={logoutHandler}>Logout</Nav.Link>
            </Nav>
          ) : (
            <Nav className="ms-auto">
              <Nav.Link href="/signup">Sign Up</Nav.Link>
              <Nav.Link href="/login">LogIn</Nav.Link>
            </Nav>
          )}
          
        </Navbar.Collapse>
      </Container>
    </Navbar>
  )
}

export default Header
