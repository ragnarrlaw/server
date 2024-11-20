package userrepository

import (
	"context"
	"fmt"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types"
)

type UserRepository interface {
	// user
	Add(context.Context, *types.UserSignUpPayload) (*types.User, error)
	GetById(context.Context, string) (*types.User, error)
	GetBy(context.Context, string, any) (*types.User, error)
	GetAll(context.Context) (*[]types.User, error)
	Remove(context.Context, []string) error
	Update(context.Context, string, *types.UserUpdatePayload) (*types.User, error)

	// user preferences
	GetUserPreferences(context.Context, string) (*types.UserPreferences, error)
	UpdateUserPreferences(context.Context, string, *types.UserPreferences) (*types.UserPreferences, error)

	// user input lists
	GetUserInputLists(context.Context, string) (*[]types.UserInputProductList, error)
	GetUserInputList(context.Context, string, string) (*types.UserInputProductList, error)
	AddUserInputList(context.Context, string, *types.UserInputListPayload) (*types.UserInputProductList, error)
	UpdateUserInputLists(context.Context, string, string, *types.UserInputListPayload) (*types.UserInputProductList, error)
	RemoveUserInputLists(context.Context, string, []string) error
}

type UserRepo struct {
	storage *database.Storage
}

func NewUserRepo(storage *database.Storage) UserRepository {
	return &UserRepo{storage: storage}
}

func (repo *UserRepo) Add(ctx context.Context, u *types.UserSignUpPayload) (*types.User, error) {
	query := `
		INSERT INTO users (username, first_name, last_name, email, contact_number, password_digest) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, username, first_name, last_name, email, contact_number, password_digest
	`
	var user types.User
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.ContactNumber,
		u.Password,
	)
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.Password,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *UserRepo) GetById(ctx context.Context, id string) (*types.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, contact_number, password_digest 
		FROM users WHERE id = $1
	`
	var user types.User

	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.Password,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *UserRepo) GetBy(ctx context.Context, field string, value any) (*types.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, first_name, last_name, email,contact_number, password_digest
		FROM users WHERE %s = $1
	`,
		field,
	)
	var user types.User

	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		value,
	).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.Password,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *UserRepo) GetAll(ctx context.Context) (*[]types.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, contact_number, password_digest
		FROM users
	`
	if rows, err := repo.storage.Pool.Query(
		ctx,
		query,
	); err != nil {
		return nil, err
	} else {
		var users []types.User
		for rows.Next() {
			var user types.User
			if err := rows.Scan(
				&user.Id,
				&user.Username,
				&user.FirstName,
				&user.LastName,
				&user.Email,
				&user.ContactNumber,
				&user.Password,
			); err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		return &users, nil
	}
}

func (repo *UserRepo) Remove(ctx context.Context, ids []string) error {
	query := `DELETE FROM user WHERE id = ANY($1)`

	if _, err := repo.storage.Pool.Exec(
		ctx,
		query,
		ids,
	); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *UserRepo) Update(ctx context.Context, id string, u *types.UserUpdatePayload) (*types.User, error) {
	query := `
		UPDATE users SET username = $1, first_name = $2, last_name = $3 WHERE id = $4
		RETURNING id, username, first_name, last_name, email, contact_number
	`
	var user types.User

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		u.Username,
		u.FirstName,
		u.LastName,
		id,
	)

	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *UserRepo) GetUserPreferences(ctx context.Context, userId string) (*types.UserPreferences, error) {
	query := `
		SELECT user_id, preferences, created_at, updated_at FROM user_preferences WHERE user_id = $1 
	`
	var userPreferences types.UserPreferences
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		userId,
	).Scan(
		&userPreferences.UserId,
		&userPreferences.Preferences,
		&userPreferences.CreatedAt,
		&userPreferences.UpdateAt,
	); err != nil {
		return nil, err
	} else {
		return &userPreferences, nil
	}
}

func (repo *UserRepo) UpdateUserPreferences(ctx context.Context, userId string, preferences *types.UserPreferences) (*types.UserPreferences, error) {
	query := `
		UPDATE user_preferences SET preferences = $1 WHERE id = $2
		RETURNING userId, preferences, created_at, updated_at
	`
	var userPreferences types.UserPreferences
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		preferences,
		userId,
	).Scan(
		&userPreferences.UserId,
		&userPreferences.Preferences,
		&userPreferences.CreatedAt,
		&userPreferences.UpdateAt,
	); err != nil {
		return nil, err
	} else {
		return &userPreferences, nil
	}
}

func (repo *UserRepo) GetUserInputLists(ctx context.Context, userId string) (*[]types.UserInputProductList, error) {
	query := `
		SELECT id, user_id, lists, created_at, updated_at FROM user_shopping_lists WHERE user_id = $1
	`
	var userInputProductList []types.UserInputProductList
	if rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		userId,
	); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var list types.UserInputProductList
			if err := rows.Scan(
				&list.Id,
				&list.UserId,
				&list.Items,
				&list.CreatedAt,
				&list.UpdatedAt,
			); err != nil {
				return nil, err
			}
			userInputProductList = append(userInputProductList, list)
		}
		return &userInputProductList, nil
	}
}

func (repo *UserRepo) GetUserInputList(ctx context.Context, userId string, listId string) (*types.UserInputProductList, error) {
	query := `
		SELECT id, user_id, lists, created_at, updated_at FROM user_shopping_lists WHERE user_id = $1 AND id = $2
	`
	var userInputProductList types.UserInputProductList
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		userId,
		listId,
	).Scan(
		&userInputProductList.Id,
		&userInputProductList.UserId,
		&userInputProductList.Items,
		&userInputProductList.CreatedAt,
		&userInputProductList.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &userInputProductList, nil
	}
}

func (repo *UserRepo) AddUserInputList(ctx context.Context, userId string, list *types.UserInputListPayload) (*types.UserInputProductList, error) {
	query := `
		INSERT INTO user_shopping_lists (user_id, lists) VALUES ($1, $2) RETURNING id, 
		user_id, lists, created_at, updated_at`
	var userInputProductList types.UserInputProductList
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		userId,
		list.List,
	).Scan(
		&userInputProductList.Id,
		&userInputProductList.UserId,
		&userInputProductList.Items,
		&userInputProductList.CreatedAt,
		&userInputProductList.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &userInputProductList, nil
	}
}

func (repo *UserRepo) UpdateUserInputLists(ctx context.Context, userId string, listId string, list *types.UserInputListPayload) (*types.UserInputProductList, error) {
	query := `
		UPDATE user_shopping_lists SET lists = $1 WHERE user_id = $2 AND id = $3
		RETURNING id, user_id, lists, created_at, updated_at
	`
	var userInputProductList types.UserInputProductList
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		list.List,
		userId,
		listId,
	).Scan(
		&userInputProductList.Id,
		&userInputProductList.UserId,
		&userInputProductList.Items,
		&userInputProductList.CreatedAt,
		&userInputProductList.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &userInputProductList, nil
	}
}

func (repo *UserRepo) RemoveUserInputLists(ctx context.Context, userId string, listId []string) error {
	query := `
		DELETE FROM user_shopping_lists WHERE user_id = $1 AND id = ANY($2)
	`
	if _, err := repo.storage.Pool.Exec(
		ctx,
		query,
		userId,
		listId,
	); err != nil {
		return err
	} else {
		return nil
	}
}
