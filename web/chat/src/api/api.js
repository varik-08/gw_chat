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

export const getUsers = async (token) => {
    const options = {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
    };

    const response = await fetch(`${API_URL}/users`, options);

    if (!response.ok) {
        throw new Error(`Request failed: ${response.statusText}`);
    }

    return response.json();
};

export const createChat = async (name, isPublic, token) => {
    const response = await fetch(`${API_URL}/chats`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
            name: name,
            is_public: isPublic,
        }),
    });

    if (!response.ok) throw new Error("Chat creation failed");

    return response.json();
}

export const addMemberToChat = async (chatID, userID, token) => {
    const response = await fetch(`${API_URL}/chats/add-member`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
            chat_id: chatID,
            user_id: userID,
        }),
    });

    if (!response.ok) throw new Error("Member addition failed");
}