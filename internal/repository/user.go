package repository

import (
	"context"
	"fmt"
	"rest1/internal/domain"

	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	conn *pgx.Conn
}

func NewUserRepo(conn *pgx.Conn) *UserRepo {
	return &UserRepo{conn: conn}
}
var  count = 0
func (u *UserRepo) CreateTable() error {
	if count == 0 {
		_1 , err1 := u.conn.Exec(context.Background() , "DROP TABLE users;")
		if(err1 != nil){
			fmt.Println(err1)
			fmt.Println(_1)
		}
		count++
	}
	_ , err := u.conn.Exec(context.Background() , "CREATE TABLE users (id INT GENERATED ALWAYS AS IDENTITY, accountno INT, name VARCHAR ( 50 )  NOT NULL,password VARCHAR ( 50 ) NOT NULL, PRIMARY KEY(id), CONSTRAINT fk_account FOREIGN KEY(accountno) REFERENCES accounts(accountno));")
	if err != nil{
		fmt.Println(err)
	}
	return nil
}

func (u *UserRepo) GetAll() ([]domain.User, error) {
	rows, err := u.conn.Query(context.Background(), "SELECT id, name, accountNo, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.AccountNo, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepo) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := u.conn.QueryRow(context.Background(), "SELECT id, name, accountNo, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.AccountNo, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) CreateUser(user *domain.User) error {
	var id int
	err := u.conn.QueryRow(context.Background(), "INSERT INTO users(id , name, accountNo, password) VALUES($1, $2, $3 , $4) RETURNING id", user.ID , user.Name , user.AccountNo , user.Password).Scan(&id)
	if(err != nil) {
		fmt.Println(err)
	}
	return nil
}

func (u *UserRepo) Withdraw(user *domain.User, amount int) error {
	qry := "UPDATE accounts SET accounts.balance = (accounts.balance - $1) WHERE accounts.accountno = $2"
	_, err := u.conn.Exec(context.Background(), qry, amount, user.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *UserRepo) Deposit(user *domain.User, amount int) error {
	/* 
		UPDATE product
		SET net_price = price - price * discount
		FROM product_segment
		WHERE product.segment_id = product_segment.id;
	*/
	qry := "UPDATE accounts SET accounts.balance = (accounts.balance + $1) WHERE accounts.accountno = $2"
	_, err := u.conn.Exec(context.Background(), qry, amount, user.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}