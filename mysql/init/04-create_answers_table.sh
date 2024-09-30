#!/bin/sh

echo "### initialize answers start ####"
CMD_MYSQL="mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}"
$CMD_MYSQL -e "CREATE TABLE IF NOT EXISTS answers (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    question_id CHAR(36) NOT NULL,
    folder_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);"
echo "### initialize finish ####"
