import './App.css'
import { LoginForm } from './components/LoginForm'
import { AuthProvider } from './context/AuthContext'

function App() {
  return (
    <AuthProvider> 
      <div><LoginForm /></div>
    </AuthProvider>
  )
}

export default App
