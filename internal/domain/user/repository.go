package user

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"message/internal/pkg/pagination"

	"github.com/pkg/errors"

	"message/internal/pkg/apperror"
)

type Repository interface {
	Get(ctx context.Context, id uint) (c *User, err error)
	First(ctx context.Context, cond *User) (c *User, err error)
	Query(ctx context.Context, cond *User, pag *pagination.Pagination) (Users []User, err error)
	Create(ctx context.Context, c *User) (uint, error)
	Update(ctx context.Context, c *User) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context, cond *User) (uint, error)
}

type repository struct {
	readerNodes []*sql.DB
	writerNodes []*sql.DB
}

func NewRepository(readerNodes, writerNodes []*sql.DB) Repository {
	return &repository{
		readerNodes: readerNodes,
		writerNodes: writerNodes,
	}
}

func (m repository) Get(ctx context.Context, id uint) (u *User, err error) {
	u = &User{}

	query := "SELECT id, email, password, name, surname, age, sex, city, interest FROM users WHERE id = ?"

	err = getRandomNode(m.readerNodes).QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Surname, &u.Age, &u.Sex, &u.City, &u.Interest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, errors.Wrap(err, "cannot get by id User")
	}

	return u, nil
}

func (m repository) First(ctx context.Context, cond *User) (u *User, err error) {
	params := []interface{}{}
	query := "SELECT id, email, password, name, surname, age, sex, city, interest FROM users WHERE 1"
	if cond.Email != "" {
		query += " AND email = ?"
		params = append(params, cond.Email)
	}
	if cond.Password != "" {
		query += " AND password = ?"
		params = append(params, cond.Password)
	}
	if cond.Name != "" {
		query += " AND name = ?"
		params = append(params, cond.Name)
	}
	if cond.Surname != "" {
		query += " AND surname = ?"
		params = append(params, cond.Surname)
	}
	if cond.City != "" {
		query += " AND city = ?"
		params = append(params, cond.City)
	}
	if cond.Interest != "" {
		query += " AND interest = ?"
		params = append(params, cond.Interest)
	}

	u = &User{}
	err = getRandomNode(m.readerNodes).QueryRowContext(ctx, query, params...).
		Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Surname, &u.Age, &u.Sex, &u.City, &u.Interest)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m repository) Query(ctx context.Context, cond *User, pag *pagination.Pagination) (Users []User, err error) {
	params := []interface{}{}
	query := "SELECT id, email, password, name, surname, age, sex, city, interest FROM users WHERE 1"
	if cond.Email != "" {
		query += " AND email = ?"
		params = append(params, cond.Email)
	}
	if cond.Password != "" {
		query += " AND password = ?"
		params = append(params, cond.Password)
	}
	if cond.Name != "" {
		query += " AND name = ?"
		params = append(params, cond.Name)
	}
	if cond.Surname != "" {
		query += " AND surname = ?"
		params = append(params, cond.Surname)
	}
	if cond.City != "" {
		query += " AND city = ?"
		params = append(params, cond.City)
	}
	if cond.Interest != "" {
		query += " AND interest = ?"
		params = append(params, cond.Interest)
	}

	query += fmt.Sprintf(" LIMIT %d, %d", pag.GetOffset(), pag.GetLimit())

	rows, err := getRandomNode(m.readerNodes).QueryContext(ctx, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Surname, &u.Age, &u.Sex, &u.City, &u.Interest)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (m repository) Create(ctx context.Context, u *User) (uint, error) {
	query := "INSERT INTO users (email, password, name, surname, age, sex, city, interest) VALUES (?,?,?,?,?,?,?,?);"
	res, err := getRandomNode(m.writerNodes).ExecContext(ctx, query, u.Email, u.Password, u.Name, u.Surname, u.Age, u.Sex, u.City, u.Interest)
	if err != nil {
		return 0, errors.Wrap(err, "cannot create User")
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(userId), nil
}

func (m repository) Update(ctx context.Context, u *User) error {
	query := "UPDATE users SET email=?,password=?, name=?, surname=?, age=?, sex=?, city=?, interest=? WHERE id =? VALUES (?,?,?,?,?,?,?,?,?);"
	_, err := getRandomNode(m.writerNodes).ExecContext(ctx, query, u.Email, u.Password, u.Name, u.Surname, u.Age, u.Sex, u.City, u.Interest, u.ID)
	if err != nil {
		return errors.Wrap(err, "cannot update user")
	}
	return nil
}

func (m repository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := getRandomNode(m.writerNodes).ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "cannot delete user")
	}
	return nil
}

func (m repository) Count(ctx context.Context, cond *User) (uint, error) {
	params := []interface{}{}
	query := "SELECT count(*) FROM users WHERE 1"
	if cond.Email != "" {
		query += " AND email = ?"
		params = append(params, cond.Email)
	}
	if cond.Password != "" {
		query += " AND password = ?"
		params = append(params, cond.Password)
	}
	if cond.Name != "" {
		query += " AND name = ?"
		params = append(params, cond.Name)
	}
	if cond.Surname != "" {
		query += " AND surname = ?"
		params = append(params, cond.Surname)
	}
	if cond.City != "" {
		query += " AND city = ?"
		params = append(params, cond.City)
	}
	if cond.Interest != "" {
		query += " AND interest = ?"
		params = append(params, cond.Interest)
	}

	count := 0
	err := getRandomNode(m.readerNodes).QueryRowContext(ctx, query, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (m repository) getReadShard(userID int) int {
	return userID % len(m.writerNodes)
}

func (m repository) getWriteShard(userID int) int {
	return userID % len(m.readerNodes)
}

func getRandomNode(nodes []*sql.DB) *sql.DB {
	return nodes[rand.Intn(len(nodes))]
}
