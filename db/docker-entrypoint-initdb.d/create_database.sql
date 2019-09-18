CREATE USER IF NOT EXISTS 'lekcije'@'localhost' identified by 'lekcije';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER, LOCK TABLES ON lekcije.* TO 'lekcije'@'localhost';
CREATE USER IF NOT EXISTS 'lekcije'@'%' identified by 'lekcije';

CREATE DATABASE IF NOT EXISTS lekcije DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_bin;
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER, LOCK TABLES ON lekcije.* TO 'lekcije'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER, LOCK TABLES ON lekcije.* TO 'lekcije'@'%';

CREATE DATABASE IF NOT EXISTS lekcije_test DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_bin;
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER, LOCK TABLES ON lekcije_test.* TO 'lekcije'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER, LOCK TABLES ON lekcije_test.* TO 'lekcije'@'%';

FLUSH PRIVILEGES;
