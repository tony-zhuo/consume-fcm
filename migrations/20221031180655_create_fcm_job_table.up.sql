CREATE TABLE fcm_jobs (
    id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    identifier VARCHAR(100) NOT NULL,
    deliver_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);