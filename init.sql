USE smartcompost_db;

DROP TABLE IF EXISTS nodes;

CREATE TABLE nodes (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       description TEXT
);

-- Dummy data for nodes
INSERT INTO nodes (name, description) VALUES ('Node 1', 'Description for Node 1');
INSERT INTO nodes (name, description) VALUES ('Node 2', 'Description for Node 2');
INSERT INTO nodes (name, description) VALUES ('Node 3', 'Description for Node 3');
INSERT INTO nodes (name, description) VALUES ('Node 4', 'Description for Node 4');
INSERT INTO nodes (name, description) VALUES ('Node 5', 'Description for Node 5');
