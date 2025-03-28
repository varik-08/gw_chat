const API_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

export const login = async (username, password) => {
    const response = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({username, password}),
    });
    if (!response.ok) throw new Error("Invalid credentials");

    return await response.json();
};

export const register = async (username, password) => {
    const response = await fetch(`${API_URL}/auth/registration`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({username, password}),
    });

    if (!response.ok) throw new Error("Registration failed");

    return response.json();
};

export const getChats = async (token) => {
    const options = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
    };

    const response = await fetch(`${API_URL}/chats`, options);

    if (!response.ok) {
        throw new Error(`Request failed: ${response.statusText}`);
    }

    return response.json();
};

export const getChatMessages = async (chatID, token) => {
    const options = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
    };

    const response = await fetch(`${API_URL}/chats/${chatID}/messages`, options);

    if (!response.ok) {
        throw new Error(`Request failed: ${response.statusText}`);
    }

    return response.json();
};

export const logout = () => {
    localStorage.removeItem("token");
};


export const updatePassword = async (newPassword, oldPassword, token) => {
    const response = await fetch(`${API_URL}/users/change-password`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({newPassword, oldPassword}),
    });

    if (!response.ok) throw new Error("Password update failed");
}