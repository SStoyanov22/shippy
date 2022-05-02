// shippy-service-user/repository.go

package main

import (
	"context"

	pb "github.com/SStoyanov22/shippy/shippy-service-user/proto/user"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID       string `sql:"id"`
	NAME     string `sql:"email"`
	COMPANY  string `sql:"company"`
	EMAIL    string `sql:"email"`
	PASSWORD string `sql:"password"`
}

type Repository interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func MarshalCollection(users []*pb.User) []*User {
	u := make([]*User, len(users))
	for _, val := range users {
		u = append(u, MarshalUser(val))
	}

	return u
}

func MarshalUser(user *pb.User) *User {
	return &User{
		ID:       user.Id,
		NAME:     user.Name,
		COMPANY:  user.Company,
		EMAIL:    user.Email,
		PASSWORD: user.Password,
	}
}

func UnmarshalCollection(users []*User) []*pb.User {
	u := make([]*pb.User, len(users))
	for _, val := range users {
		u = append(u, UnmarshalUser(val))
	}

	return u
}

func UnmarshalUser(user *User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Name:     user.NAME,
		Company:  user.COMPANY,
		Email:    user.EMAIL,
		Password: user.PASSWORD,
	}
}
func (r *PostgresRepository) GetAll(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)

	if err := r.db.GetContext(ctx, users, "select * from users"); err != nil {
		return users, err
	}

	return users, nil
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	query := "insert into users(id, name, email, company, password) values ($1, $2, $3, $4, $5)"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.NAME, user.EMAIL, user.COMPANY, user.PASSWORD)

	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*User, error) {
	var user *User

	if err := r.db.GetContext(ctx, user, "select * from user where id = $1", id); err != nil {
		return user, err
	}

	return user, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user *User

	if err := r.db.GetContext(ctx, "select * from user where email = $1", email); err != nil {
		return user, err
	}

	return user, nil
}
