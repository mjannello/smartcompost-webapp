USE smartcompost_db;

DROP TABLE IF EXISTS nodes;

CREATE TABLE IF NOT EXISTS nodes (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     description VARCHAR(255) NOT NULL,
                                     type VARCHAR(255),
                                     last_updated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS measurements (
                                            id INT AUTO_INCREMENT PRIMARY KEY,
                                            node_id INT,
                                            value FLOAT,
                                            timestamp TIMESTAMP,
                                            FOREIGN KEY (node_id) REFERENCES nodes(id)
);


INSERT INTO nodes (description, type, last_updated) VALUES ('Node 1', 'Type A', NOW()), ('Node 2', 'Type B', NOW());

INSERT INTO measurements (node_id, value, timestamp) VALUES
    (1, 12.34, NOW()),
    (1, 56.78, NOW()),
    (2, 90.12, NOW());
