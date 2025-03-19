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

	r.Post("/compose/post", handlers.PostHandler.CreatePost())                           // Создать пост в своём профиле (по header)
	r.Get("/{username}/posts", handlers.PostHandler.GetUserPosts())                      // Посты пользователя
	r.Get("/{username}/posts/{post_id}", handlers.PostHandler.GetUserPostByID())         // Конкретный пост пользователя
	r.Patch("/{username}/posts/{post_id}", handlers.PostHandler.UpdatePostContentByID()) // Редактирование контента поста (по header)
	r.Delete("/{username}/posts/{post_id}", handlers.PostHandler.DeletePostByID())       // Удаление поста (по header)

	// GET /feed // Лента новостей. Посты тех на кого подписан пользователь (по header)

	// POST /{username}/posts/{post_id}/repost // Репост чужого поста в свой профиль (1 раз) (по header)
	// DELETE /{username}/posts/{post_id}/repost // Удаление репоста из своего профиля (по header)

	// POST /{username}/posts/{post_id}/like // Лайк на пост (1 раз) (по header)
	// DELETE /{username}/posts/{post_id}/like // Удаление лайка с поста

	r.Post("/{username}/follow", handlers.UserHandler.FollowUser())           // Подписаться на пользователя (по header)
	r.Delete("/{username}/follow", handlers.UserHandler.StopFollowingUser())  // Отписаться от пользователя (по header)
	r.Get("/{username}/followers", handlers.UserHandler.GetFollowersByUser()) // Подписчики пользователя
	r.Get("/{username}/following", handlers.UserHandler.GetFollowingByUser()) // На кого пользователь подписан

	// GET /notifications // Уведомления пользователя (по header)

	return r
}
