import './App.css'
import { AuthProvider } from './context/AuthContext'
import AppRouter from './routes/AppRouter'
// import Test from './test'

function App() {
  return (
    <AuthProvider> 
      <AppRouter />
    </AuthProvider>
  )
}

export default App
