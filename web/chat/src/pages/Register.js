import { Link, useNavigate } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import "bootstrap/dist/css/bootstrap.min.css";
import { register } from "../api/api";
import { toast } from "react-toastify";

const Register = () => {
    const navigate = useNavigate();

    const formik = useFormik({
        initialValues: { username: "", password: "", confirmPassword: "" },
        validationSchema: Yup.object({
            username: Yup.string().required("Обязательное поле"),
            password: Yup.string().required("Обязательное поле"),
            confirmPassword: Yup.string()
                .oneOf([Yup.ref("password"), null], "Пароли должны совпадать")
                .required("Обязательное поле"),
        }),
        onSubmit: async (values) => {
            try {
                await register(values.username, values.password);
                navigate("/");
            } catch (err) {
                toast.error("Произошла ошибка при регистрации");
            }
        },
    });

    return (
        <div className="d-flex justify-content-center align-items-center vh-100 text-white"
             style={{ background: "linear-gradient(135deg, #1a1a2e, #16213e)" }}>
            <div className="card p-4 shadow-lg w-100"
                 style={{ maxWidth: "400px", background: "#222831", borderRadius: "15px" }}>
                <div className="text-center mb-3">
                    <img src={require('../assets/logo.jpg')} alt="Logo"
                         className="mb-3"
                         style={{ width: "200px", borderRadius: "40%", boxShadow: "0px 4px 10px rgba(255, 255, 255, 0.2)" }} />
                    <h3 className="text-white">Создать аккаунт</h3>
                </div>
                <form onSubmit={formik.handleSubmit}>
                    <div className="mb-3">
                        <input
                            type="text"
                            name="username"
                            placeholder="Логин"
                            {...formik.getFieldProps("username")}
                            className="form-control bg-secondary text-white border-0"
                        />
                        {formik.touched.username && formik.errors.username && (
                            <div className="text-danger">{formik.errors.username}</div>
                        )}
                    </div>
                    <div className="mb-3">
                        <input
                            type="password"
                            name="password"
                            placeholder="Пароль"
                            {...formik.getFieldProps("password")}
                            className="form-control bg-secondary text-white border-0"
                        />
                        {formik.touched.password && formik.errors.password && (
                            <div className="text-danger">{formik.errors.password}</div>
                        )}
                    </div>
                    <div className="mb-3">
                        <input
                            type="password"
                            name="confirmPassword"
                            placeholder="Повторите пароль"
                            {...formik.getFieldProps("confirmPassword")}
                            className="form-control bg-secondary text-white border-0"
                        />
                        {formik.touched.confirmPassword && formik.errors.confirmPassword && (
                            <div className="text-danger">{formik.errors.confirmPassword}</div>
                        )}
                    </div>
                    <button type="submit" className="btn btn-info w-100 text-white fw-bold"
                            style={{ transition: "0.3s" }}
                            onMouseEnter={(e) => e.target.style.backgroundColor = "#17a2b8"}
                            onMouseLeave={(e) => e.target.style.backgroundColor = "#138496"}>
                        Зарегистрироваться
                    </button>
                </form>
                <p className="mt-3 text-center text-white">
                    Уже есть аккаунт? <Link to="/" className="text-warning fw-bold"
                                            style={{ transition: "0.3s" }}
                                            onMouseEnter={(e) => e.target.style.color = "#ffc107"}
                                            onMouseLeave={(e) => e.target.style.color = "#ffca2c"}>
                    Войти
                </Link>
                </p>
            </div>
        </div>
    );
};

export default Register;