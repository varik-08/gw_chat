import {useContext, useEffect, useRef, useState} from "react";
import {AuthContext} from "../context/AuthContext";
import {addMemberToChat, createChat, getChatMessages, getChats, getUsers} from "../api/api";
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
    const [showModal, setShowModal] = useState(false);
    const [chatType, setChatType] = useState("private");
    const [chatName, setChatName] = useState("");
    const [users, setUsers] = useState([]);
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [newChatId, setNewChatId] = useState(null);

    useEffect(() => {
        if (!token) return;

        loadChats();
        loadUsers();

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

    useEffect(() => {
        if (!newChatId) {
            return
        }

        setSelectedChat(chats.find(chat => chat.id === newChatId))

        setNewChatId(null)
    }, [chats])

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({behavior: "smooth"});
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

            if (userID === message.userId) {
                return
            }

            if (message.type === "typing") {
                setTypingUsers((prev) => ({
                    ...prev,
                    [message.userId]: message.isTyping ? message.username : null
                }));

                if (message.isTyping) {
                    setTimeout(() => {
                        setTypingUsers((prev) => ({
                            ...prev,
                            [message.userId]: null
                        }));
                    }, 4000);
                }
            }

            if (message.type === "message") {
                setMessages((prev) => [...prev, message.message]);
                if (message.message.userId !== userID) {
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
                chatId: selectedChat.id,
            }));
            setInput("");
        }
    };

    const handleTyping = () => {
        if (ws && selectedChat) {
            // Отправляем событие "typing"
            ws.send(JSON.stringify({
                type: "typing",
                chatId: selectedChat.id,
                isTyping: true,
            }));

            if (typingTimeout) {
                clearTimeout(typingTimeout);
            }

            const timeout = setTimeout(() => {
                ws.send(JSON.stringify({
                    type: "typing",
                    chatId: selectedChat.id,
                    isTyping: false,
                }));
            }, 2000);

            setTypingTimeout(timeout);
        }
    };

    const loadUsers = async () => {
        try {
            const data = await getUsers(token);
            setUsers(data.users);
        } catch (error) {
            toast.error("Ошибка загрузки пользователей");
        }
    };

    const createNewChat = async () => {
        if (chatType === "private" && selectedUsers.length !== 1) {
            toast.error("Выберите одного пользователя для личного чата");
            return;
        }
        if (chatType === "group" && selectedUsers.length < 2) {
            toast.error("Выберите хотя бы двух пользователей для группового чата");
            return;
        }

        try {
            const newChat = await createChat(
                chatType === "group" ? chatName : `${username} / ${selectedUsers[0].username}`,
                chatType === "group",
                token,
            );

            for (let i = 0; i < selectedUsers.length; i++) {
                await addMemberToChat(newChat.id, selectedUsers[i].id, token);
            }

            await loadChats();
            setShowModal(false);
            setNewChatId(newChat.id)
            setSelectedUsers([]);
            setChatName('')
        } catch (error) {
            toast.error("Ошибка создания чата");
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
                            <button className="btn btn-success mx-2" onClick={() => setShowModal(true)}>+ Чат</button>

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
                    <ul
                        className="list-group"
                        style={{
                            maxHeight: "85vh",
                            overflowY: "auto",
                        }}
                    >
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
                                    border: "none",
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

                                <div className="text-white" style={{fontSize: "0.9rem",}}>
                                    {Object.values(typingUsers).filter(Boolean).length > 0 && (
                                        <span>{Object.values(typingUsers).filter(Boolean).join(", ")} печатает...</span>
                                    )}
                                </div>
                            </div>

                            <div className="flex-grow-1 overflow-auto p-3"
                                 style={{height: "85vh", background: "#1a1a2e"}}>
                                {messages.map((msg, index) => (
                                    <div key={msg.id}>
                                        {msg.userId !== userID &&
                                            <><strong style={{
                                                fontSize: "0.9rem",
                                                fontWeight: "bold",
                                                color: "#f1f1f1",
                                                marginBottom: "5px",
                                                marginLeft: msg.userId === userID ? "10px" : "0",
                                            }}>{msg.username}
                                            </strong><br/></>}
                                        <div key={index}
                                             className={`d-flex mb-2 ${msg.userId === userID ? "justify-content-end" : ""}`}>
                                            <div
                                                className={`p-3 rounded`}
                                                style={{
                                                    maxWidth: "75%",
                                                    background: msg.userId === userID ? "#17a2b8" : "#393e46",
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

                                <div ref={messagesEndRef}/>
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

            {showModal && (
                <div className="modal show d-block" style={{background: "rgba(0,0,0,0.5)"}}>
                    <div className="modal-dialog">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h5 className="modal-title">Создать чат</h5>
                                <button type="button" className="btn-close"
                                        onClick={() => setShowModal(false)}></button>
                            </div>
                            <div className="modal-body">
                                <select className="form-control mb-2" value={chatType}
                                        onChange={(e) => setChatType(e.target.value)}>
                                    <option value="private">Личный</option>
                                    <option value="group">Групповой</option>
                                </select>
                                {chatType === "group" &&
                                    <input type="text" className="form-control mb-2" placeholder="Название чата"
                                           value={chatName} onChange={(e) => setChatName(e.target.value)}/>}
                                <div className="mb-2">
                                    {chatType === "private" ? (
                                        <select className="form-control"
                                                onChange={(e) => setSelectedUsers([users.find(user => user.id === Number(e.target.value))])}>
                                            <option value="">Выберите пользователя</option>
                                            {users.filter(user => user.id !== userID)
                                                .map(user => (
                                                    <option key={user.id} value={user.id}>{user.username}</option>
                                                ))}
                                        </select>
                                    ) : (
                                        users.filter(user => user.id !== userID)
                                            .map(user => (
                                                <div key={user.id} className="form-check">
                                                    <input className="form-check-input" type="checkbox" value={user.id}
                                                           onChange={(e) => {
                                                               if (e.target.checked) {
                                                                   setSelectedUsers(prev => [...prev, user]);
                                                               } else {
                                                                   setSelectedUsers(prev => prev.filter(sUser => sUser.id !== user.id));
                                                               }
                                                           }}/>
                                                    <label
                                                        className="form-check-label text-black">{user.username}</label>
                                                </div>
                                            ))
                                    )}
                                </div>
                            </div>
                            <div className="modal-footer">
                                <button className="btn btn-primary" onClick={createNewChat}>Создать</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Chat;