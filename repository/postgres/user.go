package postgres

import (
	"UserGrpcProj/models"
	"UserGrpcProj/repository"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type userRepoPostgres struct {
	db  *pgxpool.Pool
	ctx context.Context
	logger *zap.Logger
}

func NewUserRepo(ctx context.Context, db *pgxpool.Pool, logger *zap.Logger) repository.UserRepository {
	return &userRepoPostgres{
		db: db,
		ctx: ctx,
		logger: logger,
	}
}
func (u *userRepoPostgres) AddUser(user *models.User) (bool, error) {


	if _, err := u.db.Exec(u.ctx, "INSERT INTO users (login, password, name, phone) VALUES ($1,$2,$3,$4)", user.Login, GetMD5Hash(user.Password), user.Name, user.Phone); err != nil{
		return false, err
	}
	return true, nil
}

func (u *userRepoPostgres) RemoveUser(id int) (bool, error) {
	if qwe, err := u.db.Exec(u.ctx, "DELETE FROM users WHERE id=$1", id); err != nil{
		fmt.Println(qwe)
		return false, err
	}

	return true, nil
}

func (u *userRepoPostgres) UserList(filter string) (*models.UserList, error) {
	return nil,nil
}

func (u *userRepoPostgres) IsSetUser(byLogin string, byId int) bool {
	var id sql.NullInt64
	if byLogin != ""{
		if err := u.db.QueryRow(u.ctx, "SELECT id FROM users WHERE login=$1", byLogin).Scan(&id); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				u.logger.Error("error", zap.Error(err))
			}
			return false
		}
		return true
	}
	if byId > 0{
		if err := u.db.QueryRow(u.ctx, "SELECT id FROM users WHERE id=$1", byId).Scan(&id); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				u.logger.Error("error", zap.Error(err))
			}
			return false
		}
	}
	return true
}
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
