import { useState, useEffect, useContext } from "react";
import { AuthContext } from "../context/AuthContext";
// import { fetchWithAuth, changePassword } from "../api/authApi";
import "bootstrap/dist/css/bootstrap.min.css";
import { toast } from "react-toastify";
import {Link, useNavigate} from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import {updatePassword} from "../api/api";

const Profile = () => {
    const { accessToken: token, username, logout } = useContext(AuthContext);
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    // Формируем схему валидации с Yup
    const formik = useFormik({
        initialValues: { currentPassword: "", newPassword: "", confirmPassword: "" },
        validationSchema: Yup.object({
            currentPassword: Yup.string().required("Текущий пароль обязателен"),
            newPassword: Yup.string()
                .required("Новый пароль обязателен"),
            confirmPassword: Yup.string()
                .oneOf([Yup.ref("newPassword"), null], "Пароли должны совпадать")
                .required("Подтверждение пароля обязательно"),
        }),
        onSubmit: async (values) => {
            if (values.newPassword !== values.confirmPassword) {
                toast.error("Пароли не совпадают");
                return;
            }
            setIsLoading(true);
            try {
                await updatePassword(values.newPassword, values.currentPassword, token);
                toast.success("Пароль успешно изменён!");
                formik.resetForm();
            } catch (err) {
                console.log(err)
                toast.error("Произошла ошибка при смене пароля");
            } finally {
                setIsLoading(false);
            }
        },
    });

    // Выход из аккаунта
    const handleLogout = () => {
        logout();
        navigate("/");
    };

    return (
        <div className="container-fluid vh-100 d-flex flex-column text-white"
             style={{ background: "linear-gradient(135deg, #1a1a2e, #16213e)" }}>
            <div className="row flex-grow-1">
                {/* Левая панель с меню */}
                <div className="col-md-4 border-end p-3"
                     style={{ background: "#222831", borderRight: "2px solid #393e46" }}>
                    {/* Логотип и меню */}
                    <div className="d-flex justify-content-between align-items-center mb-3">
                        <img src={require('../assets/logo.jpg')} alt="Logo"
                             style={{ width: "50px", borderRadius: "40%", boxShadow: "0px 2px 5px rgba(255, 255, 255, 0.2)" }} />
                        {/* Сэндвич-меню */}
                        <div className="dropdown">
                            <button className="btn btn-secondary dropdown-toggle" type="button" id="menuButton"
                                    data-bs-toggle="dropdown" aria-expanded="false">
                                ☰
                            </button>
                            <ul className="dropdown-menu dropdown-menu-dark" aria-labelledby="menuButton">
                                <li><Link className="dropdown-item" to="/chat">Чаты</Link></li>
                                <li><Link className="dropdown-item active" to="/profile">Профиль</Link></li>
                            </ul>
                        </div>
                    </div>
                    {/* Меню завершено */}
                </div>

                {/* Правая панель - профиль пользователя */}
                <div className="col-md-8 d-flex flex-column">
                    <div className="p-3 border-bottom shadow-sm" style={{ background: "#222831", borderBottom: "2px solid #393e46" }}>
                        <h5 className="m-0 text-white">Профиль</h5>
                    </div>

                    {/* Имя пользователя */}
                    <div className="p-3" style={{ background: "#1a1a2e" }}>
                        <h6 className="text-white">Имя пользователя: <span className="text-info">{username}</span></h6>
                    </div>

                    {/* Смена пароля */}
                    <div className="p-3" style={{ background: "#1a1a2e" }}>
                        <h5 className="mb-3 text-white">Сменить пароль</h5>
                        <form onSubmit={formik.handleSubmit}>
                            <div className="mb-3">
                                <input
                                    type="password"
                                    className="form-control bg-secondary text-white border-0"
                                    placeholder="Текущий пароль"
                                    {...formik.getFieldProps("currentPassword")}
                                />
                                {formik.touched.currentPassword && formik.errors.currentPassword && (
                                    <div className="text-danger">{formik.errors.currentPassword}</div>
                                )}
                            </div>
                            <div className="mb-3">
                                <input
                                    type="password"
                                    className="form-control bg-secondary text-white border-0"
                                    placeholder="Новый пароль"
                                    {...formik.getFieldProps("newPassword")}
                                />
                                {formik.touched.newPassword && formik.errors.newPassword && (
                                    <div className="text-danger">{formik.errors.newPassword}</div>
                                )}
                            </div>
                            <div className="mb-3">
                                <input
                                    type="password"
                                    className="form-control bg-secondary text-white border-0"
                                    placeholder="Повторите новый пароль"
                                    {...formik.getFieldProps("confirmPassword")}
                                />
                                {formik.touched.confirmPassword && formik.errors.confirmPassword && (
                                    <div className="text-danger">{formik.errors.confirmPassword}</div>
                                )}
                            </div>
                            <button type="submit" className="btn btn-info w-100 text-white fw-bold"
                                    style={{ transition: "0.3s" }}
                                    onMouseEnter={(e) => e.target.style.backgroundColor = "#17a2b8"}
                                    onMouseLeave={(e) => e.target.style.backgroundColor = "#138496"}
                                    disabled={isLoading}>
                                {isLoading ? "Загрузка..." : "Сменить пароль"}
                            </button>
                        </form>
                    </div>

                    {/* Кнопка Выйти */}
                    <div className="p-3" style={{ background: "#1a1a2e" }}>
                        <button className="btn btn-danger w-100 text-white fw-bold"
                                onClick={handleLogout}
                                style={{ transition: "0.3s" }}
                                onMouseEnter={(e) => e.target.style.backgroundColor = "#dc3545"}
                                onMouseLeave={(e) => e.target.style.backgroundColor = "#c82333"}>
                            Выйти
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Profile;