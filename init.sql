USE smartcompost_db;

DROP TABLE IF EXISTS nodes;

CREATE TABLE IF NOT EXISTS nodes (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     serial_number CHAR(36) NOT NULL,
                                     description VARCHAR(255) NOT NULL,
                                     model VARCHAR(255),
                                     date_created TIMESTAMP,
                                     last_updated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS measurements (
                                            id INT AUTO_INCREMENT PRIMARY KEY,
                                            node_id INT,
                                            value FLOAT,
                                            timestamp TIMESTAMP,
                                            type VARCHAR(255) NOT NULL DEFAULT '',
                                            FOREIGN KEY (node_id) REFERENCES nodes(id)
);


INSERT INTO nodes (serial_number, description, model, date_created, last_updated) VALUES (UUID(),'Node 1', 'Model A', NOW(), NOW());
INSERT INTO nodes (serial_number, description, model, date_created, last_updated) VALUES (UUID(),'Node 2', 'Model B', NOW(), NOW());

INSERT INTO measurements (node_id, value, timestamp, type) VALUES (1, 10.5, '2024-06-14 12:30:45', 'humidity');
INSERT INTO measurements (node_id, value, timestamp, type) VALUES (1, 15.3, '2024-06-14 12:35:21', 'temperature');
INSERT INTO measurements (node_id, value, timestamp, type) VALUES (2, 26.8, '2024-06-15 11:20:55', 'humidity');

