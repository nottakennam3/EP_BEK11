package store

import (
	"database/sql"
	"fmt"

	"gosocial/types"
	"github.com/go-sql-driver/mysql"
)

type UserStorage interface {
	GetUserByID(int) (*types.User, error)
	GetUserByUsername(string) (*types.User, error)
	CreateUser(*types.User) error
	UpdateUser(*types.User) error
}

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &MySQLStorage{
		db: db,
	}, nil
}

func (store *MySQLStorage) Ping() error {
	return store.db.Ping()
}

func (store *MySQLStorage) Init() error {
	createUsersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		username VARCHAR(50) NOT NULL,
		password VARCHAR(255) NOT NULL,
		userProfile VARCHAR(255),
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		UNIQUE KEY(username)
	);`
	_, err := store.db.Exec(createUsersTableQuery)
	if err != nil {
		return err
	}

	createPostsTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		userID INT UNSIGNED NOT NULL,
		content TEXT,
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (userID) REFERENCES users(id)
	);`
	_, err = store.db.Exec(createPostsTableQuery)
	if err != nil {
		return err
	}

	createLikesTableQuery := `
	CREATE TABLE IF NOT EXISTS likes (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		postID INT UNSIGNED NOT NULL,
		userID INT UNSIGNED NOT NULL,
		timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (postID) REFERENCES posts(id),
		FOREIGN KEY (userID) REFERENCES users(id)
	);`
	_, err = store.db.Exec(createLikesTableQuery)
	if err != nil {
		return err
	}

	createCommentsTableQuery := `
	CREATE TABLE IF NOT EXISTS comments (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		postID INT UNSIGNED NOT NULL,
		userID INT UNSIGNED NOT NULL,
		content TEXT,
		timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (postID) REFERENCES posts(id),
		FOREIGN KEY (userID) REFERENCES users(id)
	);`
	_, err = store.db.Exec(createCommentsTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (store *MySQLStorage) GetUserByID(id int) (*types.User, error) {
	q := "SELECT * FROM users WHERE id = ?"
	rows, err := store.db.Query(q, id)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		if err := scanRowToUser(rows, u); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (store *MySQLStorage) GetUserByUsername(username string) (*types.User, error) {
	q := "SELECT * FROM users WHERE username = ?"
	rows, err := store.db.Query(q, username)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		if err := scanRowToUser(rows, u); err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (store *MySQLStorage) CreateUser(u *types.User) error {
	q := "INSERT INTO users (username, password, userProfile) VALUES (?, ?, ?)"
	_, err := store.db.Exec(q, u.Username, u.Password, u.UserProfile)
	if err != nil {
		return err
	}
	
	return nil
}

func (store *MySQLStorage) UpdateUser(u *types.User) error {
	var err error
	if u.Password == "" {
		q := "UPDATE users SET userProfile = ? WHERE id = ?;"
		_, err = store.db.Exec(q, u.UserProfile, u.ID)
		if err != nil {
			return err
		}
	}
	if u.UserProfile == "" {
		q := "UPDATE users SET password = ? WHERE id = ?;"
		_, err = store.db.Exec(q, u.Password, u.ID)
		if err != nil {
			return err
		}
	}
	q := "UPDATE users SET password = ?, userProfile = ? WHERE id = ?;"
		_, err = store.db.Exec(q, u.Password, u.UserProfile, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLStorage) CreatePost(p *types.Post) error {
	q := "INSERT INTO posts (userID, content) VALUES (?, ?)"
	_, err := store.db.Exec(q, p.UserID, p.Content)
	if err != nil {
		return err
	}
	
	return nil
}

func (store *MySQLStorage) GetPostByID(id int) (*types.Post, error) {
	q := "SELECT * FROM posts WHERE id = ?"
	rows, err := store.db.Query(q, id)
	if err != nil {
		return nil, err
	}
	p := new(types.Post)
	for rows.Next() {
		if err := scanRowToPost(rows, p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (store *MySQLStorage) UpdatePost(p *types.Post) error {
	q := "UPDATE posts SET content = ? WHERE id = ?"
	_, err := store.db.Exec(q, p.Content, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLStorage) GetPostLikeByUserID(postID, userID int) (*types.PostLike, error) {
	q := "SELECT * FROM likes WHERE postID = ? AND userID = ?"
	rows, err := store.db.Query(q, postID, userID)
	if err != nil {
		return nil, err
	}
	like := new(types.PostLike)
	for rows.Next() {
		if err := scanRowToPostLike(rows, like); err != nil {
			return nil, err
		}
	}
	return like, nil
}

func (store *MySQLStorage) LikePost(postID, userID int) error {
	q := "INSERT INTO likes (postID, userID) VALUES (?, ?)"
	_, err := store.db.Exec(q, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLStorage) UnlikePost(postID, userID int) error {
	q := "DELETE FROM likes WHERE postID = ? AND userID = ?"
	_, err := store.db.Exec(q, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLStorage) CommentPost(pc *types.PostComment) error {
	q := "INSERT INTO comments (postID, userID, content) VALUES (?, ?, ?)"
	_, err := store.db.Exec(q, pc.PostID, pc.UserID, pc.Content)
	if err != nil {
		return err
	}
	return nil
}

func scanRowToUser(rows *sql.Rows, u *types.User) error {
	return rows.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.UserProfile,
		&u.CreatedAt,
	)
} 

func scanRowToPost(rows *sql.Rows, p *types.Post) error {
	return rows.Scan(
		&p.ID,
		&p.UserID,
		&p.Content,
		&p.CreatedAt,
	)
}

func scanRowToPostLike(rows *sql.Rows, pl *types.PostLike) error {
	return rows.Scan(
		&pl.ID,
		&pl.PostID,
		&pl.UserID,
		&pl.Timestamp,
	)
}