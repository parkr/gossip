
package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20150103160145(txn *sql.Tx) {
    txn.Exec(`
    CREATE TABLE messages (
        id int(11) NOT NULL AUTO_INCREMENT,
        room varchar(255) DEFAULT NULL,
        author varchar(255) DEFAULT NULL,
        message text,
        at datetime DEFAULT NULL,
        created_at datetime NOT NULL,
        updated_at datetime NOT NULL,
        PRIMARY KEY (id)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;`)
}

// Down is executed when this migration is rolled back
func Down_20150103160145(txn *sql.Tx) {
	txn.Exec("TRUNCATE TABLE messages; DROP TABLE messages;")
}
