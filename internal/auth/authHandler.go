package auth

import (
	"LinkShorty/configs"
	"LinkShorty/pkg/jwt"
	"LinkShorty/pkg/request"
	"LinkShorty/pkg/res"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(mux *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	mux.HandleFunc("POST /auth/register", handler.Register())
	mux.HandleFunc("POST /auth/login", handler.Login())
}

// Register обрабатывает регистрацию нового пользователя.
// @Tags Auth
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и возвращает JWT-токен
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} auth.RegisterResponse "Успешная регистрация"
// @Failure 400 {object} res.ErrorResponse "Невалидные данные (например, некорректный email или пароль меньше 8 символов)"
// @Failure 409 {object} res.ErrorResponse "Пользователь с таким email уже существует"
// @Failure 500 {object} res.ErrorResponse "Внутренняя ошибка сервера (например, ошибка генерации токена)"
// @Router /auth/register [post]
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenJWT, err := jwt.NewJWT(handler.Config.Auth.Secret).CreateToken(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: tokenJWT,
		}
		res.JsonResponse(w, data, http.StatusCreated)
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandelBody[LoginRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenJWT, err := jwt.NewJWT(handler.Config.Auth.Secret).CreateToken(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: tokenJWT,
		}
		res.JsonResponse(w, data, http.StatusOK)
	}
}
