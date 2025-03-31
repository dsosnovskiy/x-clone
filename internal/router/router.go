package router

import (
	"net/http"
	"x-clone/internal/handler"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	AuthHandler *handler.AuthHandler
	PostHandler *handler.PostHandler
	UserHandler *handler.UserHandler
}

func New(handlers *Handlers, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// В будущем добавить Refresh токен...
	// r.Post("/auth/refresh", handlers.AuthHandler.Refresh())   // Рефреш токенов
	// r.Post("/auth/logout", handlers.AuthHandler.Logout())     // Выйти из аккаунта

	// Пока только Access токен
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", handlers.AuthHandler.Register()) // Регистрация
		r.Post("/auth/login", handlers.AuthHandler.Login())       // Логин
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware) // Применяем middleware ко всем роутам в группе

		r.Get("/{username}", handlers.UserHandler.FindUserByUsername()) // Информация о пользователе

		r.Post("/compose/post", handlers.PostHandler.CreatePost())                           // Создать пост в своём профиле
		r.Get("/{username}/posts", handlers.PostHandler.GetUserPosts())                      // Посты пользователя
		r.Get("/{username}/posts/{post_id}", handlers.PostHandler.GetUserPostByID())         // Конкретный пост пользователя
		r.Patch("/{username}/posts/{post_id}", handlers.PostHandler.UpdatePostContentByID()) // Редактирование контента поста
		r.Delete("/{username}/posts/{post_id}", handlers.PostHandler.DeletePostByID())       // Удаление поста

		// GET /feed // Лента новостей. Посты тех на кого подписан пользователь
		// GET /notifications // Уведомления пользователя

		// PUT /settings/profile // Обновить данные пользователя

		r.Get("/{username}/reposts", handlers.PostHandler.GetUserReposts())                   // Все репосты пользователя
		r.Post("/{username}/posts/{post_id}/repost", handlers.PostHandler.RepostPost())       // Репост чужого поста в свой профиль
		r.Delete("/{username}/posts/{post_id}/repost", handlers.PostHandler.UndoRepostPost()) // Удаление репоста из своего профиля

		r.Post("/{username}/posts/{post_id}/quote", handlers.PostHandler.QuotePost()) // Цитата на пост чужого поста

		r.Post("/{username}/posts/{post_id}/like", handlers.PostHandler.LikePost())     // Лайк на пост
		r.Delete("/{username}/posts/{post_id}/like", handlers.PostHandler.UnlikePost()) // Удаление лайка с поста

		r.Post("/{username}/follow", handlers.UserHandler.FollowUser())           // Подписаться на пользователя
		r.Delete("/{username}/follow", handlers.UserHandler.StopFollowingUser())  // Отписаться от пользователя
		r.Get("/{username}/followers", handlers.UserHandler.GetFollowersByUser()) // Подписчики пользователя
		r.Get("/{username}/following", handlers.UserHandler.GetFollowingByUser()) // На кого пользователь подписан
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found"))
	})

	return r
}
