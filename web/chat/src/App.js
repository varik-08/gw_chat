import {BrowserRouter as Router, Navigate, Route, Routes} from "react-router-dom";
import {AuthContext, AuthProvider} from "./context/AuthContext";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Chat from "./pages/Chat";
import {useContext} from "react";
import {ToastContainer} from "react-toastify";
import Profile from "./pages/Profile";

const PrivateRoute = ({children}) => {
    const {accessToken, loading} = useContext(AuthContext);

    if (loading) {
        return <div>Loading...</div>;
    }

    return accessToken ? children : <Navigate to="/"/>;
};

function App() {
    return (
        <AuthProvider>
            <Router>
                <Routes>
                    <Route path="/" element={<Login/>}/>
                    <Route path="/register" element={<Register/>}/>
                    <Route path="/chat" element={<PrivateRoute><Chat/></PrivateRoute>}/>
                    <Route path="/profile" element={<PrivateRoute><Profile/></PrivateRoute>}/>
                </Routes>
            </Router>

            <ToastContainer/>
        </AuthProvider>
    );
}

export default App;