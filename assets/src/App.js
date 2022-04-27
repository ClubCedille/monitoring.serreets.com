import './App.css';
import Header from './components/Header'
import Footer from './components/Footer'
import {Container} from 'react-bootstrap'
import {BrowserRouter as Router, Route} from 'react-router-dom'
import LoginPage from './pages/LoginPage';
import HomePage from './pages/HomePage';
import SignupPage from './pages/SignupPage';
import { useEffect , useState} from 'react';

function App() {

  const [email, setEmail] = useState('')
  useEffect(() => {
    (
      async() => {
        const response = await fetch(process.env.REACT_APP_URL + 'api/user', {
          method: 'GET',
          headers: {'Content-Type': 'application-json'},
          credentials: 'include',
        })
        const data = await response.json()
        setEmail(data.email)
      })()
  })
  return (
    <Router>
      <Header email={email} setEmail={setEmail}/>
        <main>
          <Container>
            <Route exact path='/' component={() => <HomePage email={email}/>}></Route>
            <Route path='/login' component={LoginPage}></Route>
            <Route path='/signup' component={SignupPage}></Route>
          </Container>
        </main>
      <Footer/>
    </Router>
  );
}

export default App;
