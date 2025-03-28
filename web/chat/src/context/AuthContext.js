import {createContext, useEffect, useState} from "react";

export const AuthContext = createContext();

export const AuthProvider = ({children}) => {
    const [userID, setUserID] = useState(null);
    const [username, setUsername] = useState(null);
    const [accessToken, setAccessToken] = useState(null);
    const [refreshToken, setRefreshToken] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const accessToken = localStorage.getItem("accessToken");
        const refreshToken = localStorage.getItem("refreshToken");
        const username = localStorage.getItem("username");
        const userID = localStorage.getItem("user_id");

        if (accessToken) {
            setAccessToken(accessToken);
        }
        if (refreshToken) {
            setRefreshToken(accessToken);
        }
        if (username) {
            setUsername(username);
        }
        if (userID) {
            setUserID(Number(userID));
        }

        setLoading(false);
    }, []);

    const login = (token) => {
        localStorage.setItem("accessToken", token.access_token);
        localStorage.setItem("refreshToken", token.refresh_token);
        localStorage.setItem("user_id", token.user_id);
        localStorage.setItem("username", token.username);

        setAccessToken(token.access_token);
        setRefreshToken(token.refresh_token);
        setUserID(token.user_id)
        setUsername(token.username)
    };

    const logout = () => {
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
        localStorage.removeItem("user_id");
        localStorage.removeItem("username");

        setAccessToken(null);
        setRefreshToken(null);
        setUserID(null);
        setUsername(null);
    };

    return (
        <AuthContext.Provider value={{accessToken, refreshToken, login, logout, loading, username, userID}}>
            {children}
        </AuthContext.Provider>
    );
};