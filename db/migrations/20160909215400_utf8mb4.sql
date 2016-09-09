-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE auth_google ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
ALTER TABLE following_teacher ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
ALTER TABLE lesson ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
ALTER TABLE teacher ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
ALTER TABLE user ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;
ALTER TABLE user_api_token ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC;

ALTER TABLE user MODIFY COLUMN `email` varchar(255) COLLATE utf8mb4_general_ci NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
