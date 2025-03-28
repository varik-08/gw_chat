import {useContext, useEffect, useRef, useState} from "react";
import {AuthContext} from "../context/AuthContext";
import {getChatMessages, getChats} from "../api/api";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min.js";
import {Link} from "react-router-dom";
import {toast} from "react-toastify";

const Chat = () => {
    const {accessToken: token, username, userID} = useContext(AuthContext);
    const messagesEndRef = useRef(null);
    const [chats, setChats] = useState([]);
    const [selectedChat, setSelectedChat] = useState(null);
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState("");
    const [ws, setWs] = useState(null);
    const [typingTimeout, setTypingTimeout] = useState(null);
    const [typingUsers, setTypingUsers] = useState({});
    const [shouldScroll, setShouldScroll] = useState(false);

    useEffect(() => {
        if (!token) return;

        loadChats();

        const socket = initWS();

        setWs(socket);

        return () => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.close();
            }
        };
    }, [token]);

    useEffect(() => {
        if (!selectedChat) {
            return
        }

        loadMessages(selectedChat.id);
        setShouldScroll(true)
    }, [selectedChat])

    useEffect(() => {
        if (shouldScroll) {
            scrollToBottom();
            setShouldScroll(false);
        }
    }, [messages]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    const loadMessages = async (chatID) => {
        try {
            const data = await getChatMessages(chatID, token);

            setMessages(data.messages);
            scrollToBottom()
        } catch (error) {
            toast.error("Произошла ошибка при загрузке чатов");
        }
    }

    const playNotificationSound = () => {
        const audio = new Audio("/sounds/sound_notif_imessage.mp3");
        audio.play().catch((err) => console.error("Ошибка воспроизведения:", err));
    };

    const initWS = () => {
        const socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

        socket.onopen = () => console.log("Connected to WebSocket");
        socket.onmessage = (event) => {
            const message = JSON.parse(event.data);

            if (userID === message.user_id) {
                return
            }

            if (message.type === "typing") {
                setTypingUsers((prev) => ({
                    ...prev,
                    [message.user_id]: message.is_typing ? message.username : null
                }));

                if (message.is_typing) {
                    setTimeout(() => {
                        setTypingUsers((prev) => ({
                            ...prev,
                            [message.user_id]: null
                        }));
                    }, 4000);
                }
            }

            // добавить обработку события "user_status"

            if (message.type === "message") {
                setMessages((prev) => [...prev, message.message]);
                if (message.message.user_id !== userID) {
                    playNotificationSound();
                }
            }
        };

        return socket;
    }

    const loadChats = async () => {
        try {
            const data = await getChats(token);
            setChats(data.chats);
        } catch (error) {
            toast.error("Произошла ошибка при загрузке чатов");
        }
    };

    const sendMessage = async () => {
        if (ws && input.trim()) {
            ws.send(JSON.stringify({
                type: "message",
                content: input,
                chat_id: selectedChat.id,
            }));
            setInput("");
        }
    };

    const handleTyping = () => {
        if (ws && selectedChat) {
            // Отправляем событие "typing"
            ws.send(JSON.stringify({
                type: "typing",
                chat_id: selectedChat.id,
                is_typing: true,
            }));

            if (typingTimeout) {
                clearTimeout(typingTimeout);
            }

            const timeout = setTimeout(() => {
                ws.send(JSON.stringify({
                    type: "typing",
                    chat_id: selectedChat.id,
                    is_typing: false,
                }));
            }, 2000);

            setTypingTimeout(timeout);
        }
    };

    return (
        <div className="container-fluid vh-100 d-flex flex-column text-white"
             style={{background: "linear-gradient(135deg, #1a1a2e, #16213e)"}}>
            <div className="row flex-grow-1">
                <div className="col-md-4 border-end p-3"
                     style={{background: "#222831", borderRight: "2px solid #393e46"}}>
                    {/* Список чатов */}
                    <div className="d-flex justify-content-between align-items-center mb-3">
                        <img src={require('../assets/logo.jpg')} alt="Logo"
                             style={{
                                 width: "50px",
                                 borderRadius: "40%",
                                 boxShadow: "0px 2px 5px rgba(255, 255, 255, 0.2)"
                             }}/>
                        <div className="dropdown">
                            <button className="btn btn-secondary dropdown-toggle" type="button" id="menuButton"
                                    data-bs-toggle="dropdown" aria-expanded="false">
                                ☰
                            </button>
                            <ul className="dropdown-menu dropdown-menu-dark" aria-labelledby="menuButton">
                                <li><Link className="dropdown-item active" to="/chat">Чаты</Link></li>
                                <li><Link className="dropdown-item" to="/profile">Профиль</Link></li>
                            </ul>
                        </div>
                    </div>

                    <h5 className="text-center mb-3">Чаты</h5>
                    <ul className="list-group">
                        {chats.map((chat) => (
                            <li
                                key={chat.id}
                                className={`list-group-item list-group-item-action d-flex justify-content-between align-items-center 
                                    ${selectedChat?.id === chat.id ? "active" : ""}`}
                                onClick={() => setSelectedChat(chat)}
                                style={{
                                    cursor: "pointer",
                                    background: selectedChat?.id === chat.id ? "#17a2b8" : "#393e46",
                                    color: "white",
                                    border: "none"
                                }}
                            >
                                {chat.name}
                            </li>
                        ))}
                    </ul>
                </div>

                <div className="col-md-8 d-flex flex-column">
                    {selectedChat ? (
                        <div>
                            <div className="p-3 border-bottom shadow-sm"
                                 style={{background: "#222831", borderBottom: "2px solid #393e46", height: "70px"}}>
                                <h5 className="m-0 text-white">{selectedChat.name}</h5>

                                <div className="text-white" style={{ fontSize: "0.9rem", }}>
                                    {Object.values(typingUsers).filter(Boolean).length > 0 && (
                                        <span>{Object.values(typingUsers).filter(Boolean).join(", ")} печатает...</span>
                                    )}
                                </div>
                            </div>

                            <div className="flex-grow-1 overflow-auto p-3"
                                 style={{height: "85vh", background: "#1a1a2e"}}>
                                {messages.map((msg, index) => (
                                    <div key={msg.id}>
                                        {msg.user_id !== userID &&
                                            <><strong style={{
                                                fontSize: "0.9rem",
                                                fontWeight: "bold",
                                                color: "#f1f1f1",
                                                marginBottom: "5px",
                                                marginLeft: msg.user_id === userID ? "10px" : "0",
                                            }}>{msg.username}
                                            </strong><br/></>}
                                        <div key={index}
                                             className={`d-flex mb-2 ${msg.user_id === userID ? "justify-content-end" : ""}`}>
                                            <div
                                                className={`p-3 rounded`}
                                                style={{
                                                    maxWidth: "75%",
                                                    background: msg.user_id === userID ? "#17a2b8" : "#393e46",
                                                    color: "white",
                                                    wordWrap: "break-word",
                                                    borderRadius: "10px",
                                                }}
                                            >

                                                {msg.content}
                                            </div>
                                        </div>
                                    </div>
                                ))}

                                <div ref={messagesEndRef} />
                            </div>

                            {/* Поле ввода */}
                            <div className="p-3 border-top d-flex"
                                 style={{background: "#222831", borderTop: "2px solid #393e46"}}>
                                <input
                                    className="form-control me-2 text-white"
                                    style={{background: "#393e46", border: "none"}}
                                    value={input}
                                    onChange={(e) => {
                                        setInput(e.target.value);
                                        handleTyping();
                                    }}
                                    placeholder="Введите сообщение..."
                                />
                                <button className="btn btn-info text-white fw-bold"
                                        style={{transition: "0.3s"}}
                                        onMouseEnter={(e) => e.target.style.backgroundColor = "#17a2b8"}
                                        onMouseLeave={(e) => e.target.style.backgroundColor = "#138496"}
                                        onClick={sendMessage}>
                                    Отправить
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div className="d-flex justify-content-center align-items-center flex-grow-1 text-white">
                            Выберите чат, чтобы начать общение
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Chat;