CREATE TABLE IF NOT EXISTS account(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(40) NOT NULL,
    password VARCHAR(40) NOT NULL,
    account_type INT,
    UNIQUE INDEX name_uniq (name)
);

INSERT INTO account(name, password, account_type) VALUES("admin", "qwerty", "1");

CREATE TABLE IF NOT EXISTS todo(
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner_id INT NOT NULL,
    name VARCHAR(40) NOT NULL,
    `desc` TEXT NOT NULL,
    status INT,
    FOREIGN KEY(owner_id)
        REFERENCES account(id)
);
