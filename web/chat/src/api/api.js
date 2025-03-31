const API_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

const originalFetch = window.fetch;

window.fetch = async (url, options = {}) => {
    let response = await originalFetch(url, options);

    if (response.status === 401) {
        console.warn("Токен истек. Обновляем...");

        const newToken = await refreshAccessToken();

        if (!newToken) {
            console.error("Не удалось обновить токен. Требуется авторизация.");
            return response;
        }

        options.headers.Authorization = `Bearer ${newToken}`;
        response = await originalFetch(url, options);
    }

    return response;
};

export const refreshAccessToken = async () => {
    const refreshToken = localStorage.getItem("refreshToken");

    if (!refreshToken) {
        console.error("Нет refresh_token, требуется авторизация.");
        return null;
    }

    try {
        const response = await fetch(`${API_URL}/auth/refresh`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ refresh_token: refreshToken }),
        });

        if (!response.ok) {
            throw new Error("Не удалось обновить токен");
        }

        const data = await response.json();
        localStorage.setItem("accessToken", data.access_token);
        localStorage.setItem("refreshToken", data.refresh_token);

        return data.access_token;
    } catch (error) {
        console.error("Ошибка при обновлении токена:", error);
        return null;
    }
};

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