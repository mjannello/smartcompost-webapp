USE smartcompost_db;

CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL,
                                     email VARCHAR(50) NOT NULL
);

INSERT INTO users (username, email) VALUES ('user1', 'user1@example.com'), ('user2', 'user2@example.com');
