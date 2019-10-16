package db

const createUserTable = `CREATE TABLE IF NOT EXISTS 'user_info'(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	full_name VARCHAR,
	about TEXT,
	login VARCHAR UNIQUE,
	password VARCHAR
);`

const addUserScript = `INSERT INTO user_info(full_name, about, login, password) 
	VALUES (?, ?, ?, ?);`

const getAllUserScript = `SELECT id, full_name, about, login, password FROM user_info`

const getUserByIdScript = `SELECT id, full_name, about, login, password FROM user_info WHERE id = ?`

const checkUserExistScript = `SELECT id FROM user_info WHERE login = ?`

const updateUserScript = `UPDATE user_info SET full_name = ?, about = ?, login = ?, password = ? WHERE id = ?`

const deleteUserScript = `DELETE FROM user_info WHERE id = ?`
