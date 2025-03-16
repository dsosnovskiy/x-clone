package router

import (
	"x-clone/internal/handler"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	PostHandler *handler.PostHandler
	UserHandler *handler.UserHandler
}

func New(handlers *Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/createuser", handlers.UserHandler.CreateUser()) // Создать пользователя. Временно, замена ниже...
	// POST /auth/register - зарегистрировать пользователя.
	// POST /auth/login - авторизировать пользователя.
	// POST /auth/logout - выйти из аккаунта.
	// PUT /settings/profile // Обновить данные пользователя (по header)
	r.Get("/{username}", handlers.UserHandler.FindUserByUsername()) // Информация о пользователе

	r.Post("/compose/post", handlers.PostHandler.CreatePost())      // Создать пост в своём профиле (по header)
	r.Get("/{username}/posts", handlers.PostHandler.GetUserPosts()) // Посты пользователя
	// GET /{username}/posts/{post_id} // Конкретный пост пользователя
	// PATCH /{username}/posts/{post_id} // Редактирование контента поста (по header)
	// DELETE /{username}/posts/{post_id} // Удаление поста (по header)

	// GET /feed // Лента новостей. Посты тех на кого подписан пользователь (по header)

	// POST /{username}/posts/{post_id}/repost // Репост чужого поста в свой профиль (1 раз) (по header)
	// DELETE /{username}/posts/{post_id}/repost // Удаление репоста из своего профиля (по header)

	// POST /{username}/posts/{post_id}/like // Лайк на пост (1 раз) (по header)
	// DELETE /{username}/posts/{post_id}/like // Удаление лайка с поста

	// POST /{username}/follow // Подписаться на пользователя (1 раз) (по header)
	// DELETE /{username}/follow // Отписаться (по header)
	// GET /{username}/followers // Подписчики пользователя
	// GET /{username}/following // На кого пользователь подписан

	// GET /notifications // Уведомления пользователя (по header)

	return r
}
